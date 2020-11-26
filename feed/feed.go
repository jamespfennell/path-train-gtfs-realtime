package feed

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	gtfs "github.com/jamespfennell/path-train-gtfs-realtime/feed/gtfsrt"
	sourceapi "github.com/jamespfennell/path-train-gtfs-realtime/feed/sourceapi"
	"github.com/jamespfennell/path-train-gtfs-realtime/monitoring"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	grpcApiUrl          = "path.grpc.razza.dev:443"
	apiBaseUrl          = "https://path.api.razza.dev/v1/"
	apiRoutesEndpoint   = "routes/"
	apiStationsEndpoint = "stations/"
	apiRealtimeEndpoint = "stations/%sourceapi/realtime/"
)

// (1)

type Train *sourceapi.GetUpcomingTrainsResponse_UpcomingTrain

type source struct {
	data   sourceData
	client sourceClient
}

// A container for the most updated data retrieved from the API so far.
// During the program execution, there is one instance of this structure.
type sourceData struct {
	stationToStopId           map[sourceapi.Station]string
	routeToRouteId            map[sourceapi.Route]string
	stationIdToUpcomingTrains map[sourceapi.Station][]Train
}

// Initialize the sourceData data structure by populating its stop and route fields.
// If the initialization fails, the program exits.
func (s *source) initializeData() error {
	var err error
	s.data.routeToRouteId, err = s.client.GetRouteToRouteId()
	if err != nil {
		return err
	}
	s.data.stationToStopId, err = s.client.GetStationToStopId()
	if err != nil {
		return err
	}
	s.data.stationIdToUpcomingTrains = map[sourceapi.Station][]Train{}
	for apiStationId := range s.data.stationToStopId {
		s.data.stationIdToUpcomingTrains[apiStationId] = []Train{}
	}
	return nil
}

// Update the sourceData structure using the most recent stations data from the API.
// If data for one of the stations cannot be retrieved, the old data will be preserved.
// This method is highly I/O bound: the time it takes to execute is largely spent waiting for HTTP responses
// from the 14 station endpoints.
// To speed it up, the 14 endpoints are hit in parallel.
// TODO: return the errors and number of trains returned
func (s *source) updateData() (err error) {
	type trainsAtStation struct {
		Station sourceapi.Station
		Trains  []Train
		Err     error
	}
	// TODO: make this less fragile, perhaps using a wait group
	allTrainsAtStations := make(chan trainsAtStation, len(s.data.stationToStopId))
	for station := range s.data.stationToStopId {
		station := station
		go func() {
			r := trainsAtStation{Station: station}
			r.Trains, r.Err = s.client.GetTrainsAtStation(station)
			allTrainsAtStations <- r
		}()
	}
	for range s.data.stationToStopId {
		trainsAtStation := <-allTrainsAtStations
		if trainsAtStation.Err == nil {
			s.data.stationIdToUpcomingTrains[trainsAtStation.Station] = trainsAtStation.Trains
		} else {
			err = trainsAtStation.Err
		}
	}
	return
}

func convertDirectionToBoolean(direction sourceapi.Direction) *uint32 {
	var result uint32
	if direction == sourceapi.Direction_TO_NY {
		result = 1
	} else if direction == sourceapi.Direction_TO_NJ {
		result = 0
	}
	return &result
}

// (2)

// sourceClient is a way to get data from the Razza API
type sourceClient interface {
	GetStationToStopId() (map[sourceapi.Station]string, error)
	GetRouteToRouteId() (map[sourceapi.Route]string, error)
	GetTrainsAtStation(sourceapi.Station) ([]Train, error)
	Close() error
}

type httpClient struct {
	timeoutPeriod time.Duration
}

func (client *httpClient) GetTrainsAtStation(station sourceapi.Station) ([]Train, error) {
	type jsonUpcomingTrain struct {
		ProjectedArrival  string
		LastUpdated       string
		RouteAsString     string `json:"route"`
		DirectionAsString string `json:"direction"`
	}
	type jsonGetUpcomingTrainsResponse struct {
		Trains []jsonUpcomingTrain `json:"upcomingTrains"`
	}
	stationAsString := strings.ToLower(sourceapi.Station_name[int32(station)])
	realtimeApiContent, err := client.getContent(fmt.Sprintf(apiRealtimeEndpoint, stationAsString))
	if err != nil {
		return nil, err
	}
	response := jsonGetUpcomingTrainsResponse{}
	err = json.Unmarshal(realtimeApiContent, &response)
	if err != nil {
		return nil, err
	}
	var trains []Train
	for _, rawUpcomingTrain := range response.Trains {
		upcomingTrain := sourceapi.GetUpcomingTrainsResponse_UpcomingTrain{
			Route:            client.convertRouteAsStringToRoute(rawUpcomingTrain.RouteAsString),
			Direction:        client.convertDirectionAsStringToDirection(rawUpcomingTrain.DirectionAsString),
			ProjectedArrival: client.convertApiTimeStringToTimestamp(rawUpcomingTrain.ProjectedArrival),
			LastUpdated:      client.convertApiTimeStringToTimestamp(rawUpcomingTrain.LastUpdated),
		}
		trains = append(trains, &upcomingTrain)
	}
	return trains, nil
}

func (client *httpClient) GetStationToStopId() (map[sourceapi.Station]string, error) {
	stationsContent, err := client.getContent(apiStationsEndpoint)
	if err != nil {
		return nil, err
	}
	type jsonStationData struct {
		StationAsString string `json:"station"`
		Id              string
	}
	type jsonListStationsResponse struct {
		Stations []jsonStationData `json:"stations"`
	}
	response := jsonListStationsResponse{}
	err = json.Unmarshal(stationsContent, &response)
	if err != nil {
		return nil, err
	}
	stationToStopId := map[sourceapi.Station]string{}
	for _, stationData := range response.Stations {
		stationToStopId[client.convertStationAsStringToStation(stationData.StationAsString)] = stationData.Id
	}
	return stationToStopId, nil
}

func (client *httpClient) GetRouteToRouteId() (map[sourceapi.Route]string, error) {
	routesContent, err := client.getContent(apiRoutesEndpoint)
	if err != nil {
		return nil, err
	}
	type jsonRouteData struct {
		RouteAsString string `json:"route"`
		Id            string
	}
	type jsonListRoutesResponse struct {
		Routes []jsonRouteData `json:"routes"`
	}
	response := jsonListRoutesResponse{}
	err = json.Unmarshal(routesContent, &response)
	if err != nil {
		return nil, err
	}
	routeToRouteId := map[sourceapi.Route]string{}
	for _, routeData := range response.Routes {
		routeToRouteId[client.convertRouteAsStringToRoute(routeData.RouteAsString)] = routeData.Id
	}
	return routeToRouteId, nil
}

func (client *httpClient) Close() error { return nil }

func (client *httpClient) convertDirectionAsStringToDirection(directionAsString string) sourceapi.Direction {
	return sourceapi.Direction(sourceapi.Direction_value[directionAsString])
}

func (client *httpClient) convertStationAsStringToStation(stationAsString string) sourceapi.Station {
	return sourceapi.Station(sourceapi.Station_value[stationAsString])
}

func (client *httpClient) convertRouteAsStringToRoute(routeAsString string) sourceapi.Route {
	return sourceapi.Route(sourceapi.Route_value[routeAsString])
}
func (client *httpClient) convertApiTimeStringToTimestamp(timeString string) *timestamp.Timestamp {
	timeObj, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return nil
	}
	value := timestamp.Timestamp{Seconds: timeObj.Unix()}
	return &value
}

// Get the raw bytes from an endpoint in the API.
func (client httpClient) getContent(endpoint string) (bytes []byte, err error) {
	httpClient := &http.Client{Timeout: client.timeoutPeriod}
	resp, err := httpClient.Get(apiBaseUrl + endpoint)
	if err != nil {
		return
	}
	defer func() {
		closingErr := resp.Body.Close()
		if err == nil {
			err = closingErr
		}
	}()
	return ioutil.ReadAll(resp.Body)
}

type grpcClient struct {
	conn          *grpc.ClientConn
	stations      *sourceapi.StationsClient
	routes        *sourceapi.RoutesClient
	timeoutPeriod time.Duration
}

func newGrpcClient(timeoutPeriod time.Duration) (*grpcClient, error) {
	conn, err := grpc.Dial(grpcApiUrl, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	stationsClient := sourceapi.NewStationsClient(conn)
	routesClient := sourceapi.NewRoutesClient(conn)
	return &grpcClient{conn: conn, stations: &stationsClient, routes: &routesClient, timeoutPeriod: timeoutPeriod}, nil
}

func (client *grpcClient) GetStationToStopId() (stationToStopId map[sourceapi.Station]string, err error) {
	response, err := (*client.stations).ListStations(client.createContext(), &sourceapi.ListStationsRequest{})
	if err != nil {
		return
	}
	stationToStopId = map[sourceapi.Station]string{}
	for _, stationData := range response.Stations {
		stationToStopId[stationData.Station] = stationData.Id
	}
	return
}

func (client *grpcClient) GetRouteToRouteId() (routeToRouteId map[sourceapi.Route]string, err error) {
	response, err := (*client.routes).ListRoutes(client.createContext(), &sourceapi.ListRoutesRequest{})
	if err != nil {
		return
	}
	routeToRouteId = map[sourceapi.Route]string{}
	for _, routeData := range response.Routes {
		routeToRouteId[routeData.Route] = routeData.Id
	}
	return
}

func (client *grpcClient) Close() error {
	return client.conn.Close()
}

func (client *grpcClient) GetTrainsAtStation(station sourceapi.Station) ([]Train, error) {
	request := sourceapi.GetUpcomingTrainsRequest{Station: station}
	response, err := (*client.stations).GetUpcomingTrains(client.createContext(), &request)
	if err != nil {
		return nil, err
	}
	var trains []Train
	for _, train := range response.UpcomingTrains {
		trains = append(trains, train)
	}
	return trains, nil
}

func (client *grpcClient) createContext() context.Context {
	deadline := time.Now().Add(client.timeoutPeriod)
	ctx, _ := context.WithDeadline(context.Background(), deadline)
	return ctx
}

// (3)

// Build a GTFS Realtime message from a snapshot of the current data.
func buildGtfsRealtimeFeedMessage(data *sourceData) *gtfs.FeedMessage {
	gtfsVersion := "0.2"
	incrementality := gtfs.FeedHeader_FULL_DATASET
	currentTimestamp := uint64(time.Now().Unix())
	feedMessage := gtfs.FeedMessage{
		Header: &gtfs.FeedHeader{
			GtfsRealtimeVersion: &gtfsVersion,
			Incrementality:      &incrementality,
			Timestamp:           &currentTimestamp,
		},
		Entity: []*gtfs.FeedEntity{},
	}
	for apiStationId, trains := range data.stationIdToUpcomingTrains {
		for _, train := range trains {
			tripId := newPseudoTripId()
			tripUpdate, err := convertApiTrainToTripUpdate(data, train, tripId, apiStationId)
			if err != nil {
				continue
			}
			feedEntity := gtfs.FeedEntity{
				Id:         &tripId,
				TripUpdate: tripUpdate,
			}
			feedMessage.Entity = append(feedMessage.Entity, &feedEntity)
		}
	}
	return &feedMessage
}

func convertApiTrainToTripUpdate(
	data *sourceData,
	train *sourceapi.GetUpcomingTrainsResponse_UpcomingTrain,
	tripId string,
	station sourceapi.Station) (update *gtfs.TripUpdate, err error) {
	lastUpdatedUnsigned := uint64(train.LastUpdated.Seconds)
	arrivalTime := train.ProjectedArrival.Seconds
	stopId := data.stationToStopId[station]
	route := train.Route
	routeId := data.routeToRouteId[route]
	stopTimeUpdate := gtfs.TripUpdate_StopTimeUpdate{
		StopSequence: nil,
		StopId:       &stopId,
		Arrival: &gtfs.TripUpdate_StopTimeEvent{
			Time: &arrivalTime,
		},
	}
	return &gtfs.TripUpdate{
		Trip: &gtfs.TripDescriptor{
			TripId:      &tripId,
			RouteId:     &routeId,
			DirectionId: convertDirectionToBoolean(train.Direction),
		},
		StopTimeUpdate: []*gtfs.TripUpdate_StopTimeUpdate{
			&stopTimeUpdate,
		},
		Timestamp: &lastUpdatedUnsigned,
	}, nil
}

func newPseudoTripId() string {
	randomUuid, err := uuid.NewRandom()
	if err != nil {
		return ""
	} else {
		return randomUuid.String()
	}
}

// (4)

// Run one feed update iteration.
func (f *Feed) update() {
	fmt.Println("Updating GTFS Realtime feed.")
	err := f.source.updateData()
	if err != nil {
		fmt.Println("There was an error while retrieving the data; update will continue with some data stale.")
	}
	feedMessage := buildGtfsRealtimeFeedMessage(&f.source.data)
	out, err := proto.Marshal(feedMessage)
	if err != nil {
		fmt.Println("Update failed: there was an error while generating the realtime protobuf file. ")
		return
	}
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.gtfs = out
	fmt.Println("Done")
}

type Feed struct {
	gtfs    []byte
	mutex   sync.RWMutex
	source  source
	monitor *monitoring.Monitor
}

func (f *Feed) runInBackground(initChan chan<- error, updatePeriod, timeoutPeriod time.Duration, useHTTPSourceAPI bool) {
	// TODO: use log instead
	fmt.Println("Starting up")
	if !useHTTPSourceAPI {
		fmt.Println("Source API: gRPC")
		var err error
		f.source.client, err = newGrpcClient(timeoutPeriod)
		if err != nil {
			initChan <- err
			return
		}
	} else {
		fmt.Println("Source API: HTTP")
		f.source.client = &httpClient{timeoutPeriod: timeoutPeriod}
	}
	defer func() {
		err := f.source.client.Close()
		if err != nil {
			fmt.Println("Error while closing client connection:", err.Error())
		}
	}()
	err := f.source.initializeData()
	if err != nil {
		initChan <- err
		return
	}
	// Signal that initialization completed successfully
	close(initChan)
	f.update()
	ticker := time.NewTicker(updatePeriod)
	for range ticker.C {
		f.update()
	}
}

func NewFeed(monitor *monitoring.Monitor, updatePeriod, timeoutPeriod time.Duration, useHTTPSourceAPI bool) (*Feed, error) {
	initChan := make(chan error)
	f := Feed{monitor: monitor}
	go f.runInBackground(initChan, updatePeriod, timeoutPeriod, useHTTPSourceAPI)
	return &f, <-initChan
}

func (f *Feed) Get() []byte {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	r := make([]byte, len(f.gtfs))
	copy(r, f.gtfs)
	return r
}
