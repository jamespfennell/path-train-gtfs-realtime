// PATH Train GTFS Realtime feed generator.
//
// Copyright (c) James Fennell 2020.
//
// Released under the MIT License.
//
// This file is divided into 4 parts:
// (1) Structures for holding data from the API and methods for populating these structure using API data.
// (2) Functions for getting data from the API.
// (3) Functions for converting the data structures in (1) to GTFS Realtime protobuf data structures.
// (4) Code that launches the application and uses (1), (2) and (3) to create the feed.

package main

import (
	gtfs "./gtfsrt"
	s "./sourceapi"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	exitCodeMalformedEnvVarPeriodicity = 101
	exitCodeMalformedEnvVarApi         = 102
	exitCodeCannotWriteToDisk          = 103
	exitCodeCannotGetRoutesData        = 104
	exitCodeCannotGetStationsData      = 105
	requestTimeoutMilliseconds         = 3000
	envVarSourceApi                    = "PATH_GTFS_RT_SOURCE_API"
	envVarPeriodicity                  = "PATH_GTFS_RT_PERIODICITY_MILLISECONDS"
	envVarOutputPath                   = "PATH_GTFS_RT_OUTPUT_PATH"
	grpcApiUrl                         = "path.grpc.razza.dev:443"
	apiBaseUrl                         = "https://path.api.razza.dev/v1/"
	apiRoutesEndpoint                  = "routes/"
	apiStationsEndpoint                = "stations/"
	apiRealtimeEndpoint                = "stations/%s/realtime/"
)

// (1)

// A container for the most updated data retrieved from the API so far.
// During the program execution, there is one instance of this structure.
type apiData struct {
	client                    *apiClient
	stationToStopId           map[s.Station]string
	routeToRouteId            map[s.Route]string
	stationIdToUpcomingTrains map[s.Station][]*s.GetUpcomingTrainsResponse_UpcomingTrain
}

// Initialize the apiData data structure by populating its stop and route fields.
// If the initialization fails, the program exits.
func (data *apiData) initialize() {
	var err error
	data.routeToRouteId, err = (*data.client).GetRouteToRouteId()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(exitCodeCannotGetRoutesData)
	}
	data.stationToStopId, err = (*data.client).GetStationToStopId()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(exitCodeCannotGetStationsData)
	}
	data.stationIdToUpcomingTrains = map[s.Station][]*s.GetUpcomingTrainsResponse_UpcomingTrain{}
	for apiStationId := range data.stationToStopId {
		data.stationIdToUpcomingTrains[apiStationId] = []*s.GetUpcomingTrainsResponse_UpcomingTrain{}
	}
}

// Update the apiData structure using the most recent stations data from the API.
// If data for one of the stations cannot be retrieved, the old data will be preserved.
// This method is highly I/O bound: the time it takes to execute is largely spent waiting for HTTP responses
// from the 14 station endpoints.
// To speed it up, the 14 endpoints are hit in parallel.
func (data *apiData) update() (err error) {
	allTrainsAtStations := make(chan trainsAtStation, len(data.stationToStopId))
	for station := range data.stationToStopId {
		station := station
		go func() { allTrainsAtStations <- (*data.client).GetTrainsAtStation(station) }()
	}
	for range data.stationToStopId {
		trainsAtStation := <-allTrainsAtStations
		if trainsAtStation.Err == nil {
			data.stationIdToUpcomingTrains[trainsAtStation.Station] = trainsAtStation.Trains
		} else {
			err = trainsAtStation.Err
		}
	}
	return
}

func (data *apiData) convertDirectionToBoolean(direction s.Direction) *uint32 {
	var result uint32
	if direction == s.Direction_TO_NY {
		result = 1
	} else if direction == s.Direction_TO_NJ {
		result = 0
	}
	return &result
}

// (2)

type trainsAtStation struct {
	Station s.Station
	Trains  []*s.GetUpcomingTrainsResponse_UpcomingTrain
	Err     error
}

type apiClient interface {
	GetStationToStopId() (map[s.Station]string, error)
	GetRouteToRouteId() (map[s.Route]string, error)
	GetTrainsAtStation(s.Station) trainsAtStation
	Close() error
}

type httpApiClient struct{}

func (client *httpApiClient) GetTrainsAtStation(station s.Station) (result trainsAtStation) {
	result = trainsAtStation{Station: station, Trains: []*s.GetUpcomingTrainsResponse_UpcomingTrain{}}
	type jsonUpcomingTrain struct {
		ProjectedArrival  string
		LastUpdated       string
		RouteAsString     string `json:"route"`
		DirectionAsString string `json:"direction"`
	}
	type jsonGetUpcomingTrainsResponse struct {
		Trains []jsonUpcomingTrain `json:"upcomingTrains"`
	}
	stationAsString := strings.ToLower(s.Station_name[int32(station)])
	realtimeApiContent, err := client.getContent(fmt.Sprintf(apiRealtimeEndpoint, stationAsString))
	if err != nil {
		result.Err = err
		return
	}
	response := jsonGetUpcomingTrainsResponse{}
	err = json.Unmarshal(realtimeApiContent, &response)
	if err != nil {
		result.Err = err
		return
	}
	for _, rawUpcomingTrain := range response.Trains {
		upcomingTrain := s.GetUpcomingTrainsResponse_UpcomingTrain{
			Route:            client.convertRouteAsStringToRoute(rawUpcomingTrain.RouteAsString),
			Direction:        client.convertDirectionAsStringToDirection(rawUpcomingTrain.DirectionAsString),
			ProjectedArrival: client.convertApiTimeStringToTimestamp(rawUpcomingTrain.ProjectedArrival),
			LastUpdated:      client.convertApiTimeStringToTimestamp(rawUpcomingTrain.LastUpdated),
		}
		result.Trains = append(result.Trains, &upcomingTrain)
	}
	return
}

func (client *httpApiClient) GetStationToStopId() (stationToStopId map[s.Station]string, err error) {
	stationsContent, err := client.getContent(apiStationsEndpoint)
	if err != nil {
		return
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
		return
	}
	stationToStopId = map[s.Station]string{}
	for _, stationData := range response.Stations {
		stationToStopId[client.convertStationAsStringToStation(stationData.StationAsString)] = stationData.Id
	}
	return
}

func (client *httpApiClient) GetRouteToRouteId() (routeToRouteId map[s.Route]string, err error) {
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
		return
	}
	routeToRouteId = map[s.Route]string{}
	for _, routeData := range response.Routes {
		routeToRouteId[client.convertRouteAsStringToRoute(routeData.RouteAsString)] = routeData.Id
	}
	return
}

func (client *httpApiClient) Close() error { return nil }

func (client *httpApiClient) convertDirectionAsStringToDirection(directionAsString string) s.Direction {
	return s.Direction(s.Direction_value[directionAsString])
}

func (client *httpApiClient) convertStationAsStringToStation(stationAsString string) s.Station {
	return s.Station(s.Station_value[stationAsString])
}

func (client *httpApiClient) convertRouteAsStringToRoute(routeAsString string) s.Route {
	return s.Route(s.Route_value[routeAsString])
}
func (client *httpApiClient) convertApiTimeStringToTimestamp(timeString string) *timestamp.Timestamp {
	timeObj, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return nil
	}
	value := timestamp.Timestamp{Seconds: timeObj.Unix()}
	return &value
}

// Get the raw bytes from an endpoint in the API.
func (_ httpApiClient) getContent(endpoint string) (bytes []byte, err error) {
	httpClient := &http.Client{Timeout: requestTimeoutMilliseconds * time.Millisecond}
	resp, err := httpClient.Get(apiBaseUrl + endpoint)
	if err != nil {
		return
	}
	defer func() {
		cerr := resp.Body.Close()
		if err == nil {
			err = cerr
		}
	}()
	return ioutil.ReadAll(resp.Body)
}

type grpcApiClient struct {
	conn     *grpc.ClientConn
	stations *s.StationsClient
	routes   *s.RoutesClient
}

func newGrpcApiClient() *grpcApiClient {
	conn, err := grpc.Dial(grpcApiUrl, grpc.WithInsecure())
	if err != nil {
		os.Exit(57)
	}
	stationsClient := s.NewStationsClient(conn)
	routesClient := s.NewRoutesClient(conn)
	return &grpcApiClient{conn: conn, stations: &stationsClient, routes: &routesClient}
}

func (client *grpcApiClient) GetStationToStopId() (stationToStopId map[s.Station]string, err error) {
	response, err := (*client.stations).ListStations(client.createContext(), &s.ListStationsRequest{})
	if err != nil {
		return
	}
	stationToStopId = map[s.Station]string{}
	for _, stationData := range response.Stations {
		stationToStopId[stationData.Station] = stationData.Id
	}
	return
}

func (client *grpcApiClient) GetRouteToRouteId() (routeToRouteId map[s.Route]string, err error) {
	response, err := (*client.routes).ListRoutes(client.createContext(), &s.ListRoutesRequest{})
	if err != nil {
		return
	}
	routeToRouteId = map[s.Route]string{}
	for _, routeData := range response.Routes {
		routeToRouteId[routeData.Route] = routeData.Id
	}
	return
}

func (client *grpcApiClient) Close() error {
	return client.conn.Close()
}

func (client *grpcApiClient) GetTrainsAtStation(station s.Station) (result trainsAtStation) {
	request := s.GetUpcomingTrainsRequest{Station: station}
	result = trainsAtStation{Station: station, Trains: []*s.GetUpcomingTrainsResponse_UpcomingTrain{}}
	response, err := (*client.stations).GetUpcomingTrains(client.createContext(), &request)
	if err != nil {
		return
	}
	result.Trains = response.UpcomingTrains
	return
}

func (client *grpcApiClient) createContext() context.Context {
	deadline := time.Now().Add(time.Duration(requestTimeoutMilliseconds) * time.Millisecond)
	ctx, _ := context.WithDeadline(context.Background(), deadline)
	return ctx
}

// (3)

// Build a GTFS Realtime message from a snapshot of the current data.
func buildGtfsRealtimeFeedMessage(data *apiData) *gtfs.FeedMessage {
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
	data *apiData,
	train *s.GetUpcomingTrainsResponse_UpcomingTrain,
	tripId string,
	station s.Station) (update *gtfs.TripUpdate, err error) {
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
			DirectionId: data.convertDirectionToBoolean(train.Direction),
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

// Determine the desired periodicity of the update from an env var.
func getPeriodicity() int {
	periodicityString, envVarSet := os.LookupEnv(envVarPeriodicity)
	if !envVarSet {
		return 5000
	}
	periodicity, err := strconv.Atoi(periodicityString)
	if err != nil || periodicity < 1000 {
		fmt.Println(fmt.Sprintf("Expected periodicity in milliseconds to be a number; recieved '%s'; exiting.", periodicityString))
		os.Exit(exitCodeMalformedEnvVarPeriodicity)
	}
	return periodicity
}

func getSourceApi() string {
	sourceString, envVarSet := os.LookupEnv(envVarSourceApi)
	if !envVarSet {
		return "grpc"
	}
	if strings.ToLower(sourceString) == "grpc" {
		return "grpc"
	}
	if strings.ToLower(sourceString) == "http" {
		return "grpc"
	}
	os.Exit(exitCodeMalformedEnvVarApi)
	return ""
}

func ensureCanWrite(outputPath string) {
	err := ioutil.WriteFile(outputPath, []byte{}, 0644)
	if err != nil {
		fmt.Println(fmt.Sprintf("Unable to write to '%s', exiting.", outputPath))
		os.Exit(exitCodeCannotWriteToDisk)
	}
}

// Run one feed update iteration.
func run(data *apiData, outputPath string) {
	fmt.Println("Updating GTFS Realtime feed.")
	err := data.update()
	if err != nil {
		fmt.Println("There was an error while retrieving the data; update will continue with some data stale.")
	}
	feedMessage := buildGtfsRealtimeFeedMessage(data)
	out, err := proto.Marshal(feedMessage)
	if err != nil {
		fmt.Println("Update failed: there was an error while generating the realtime protobuf file. ")
		return
	}
	err = ioutil.WriteFile(outputPath, out, 0644)
	if err != nil {
		fmt.Println("Update failed: there was an error writing the GTFS Realtime file to disk.")
		return
	}
	fmt.Println("Update successful.")
}

func main() {
	fmt.Println("Starting up")

	var client apiClient
	if getSourceApi() == "grpc" {
		fmt.Println("Source API: gRPC")
		client = newGrpcApiClient()
	} else {
		fmt.Println("Source API: HTTP")
		client = &httpApiClient{}
	}
	defer func() {
		err := client.Close()
		if err != nil {
			fmt.Println("Error while closing client connection:", err.Error())
		}
	}()
	data := apiData{client: &client}

	outputPath, envVarSet := os.LookupEnv(envVarOutputPath)
	if !envVarSet {
		outputPath = "path.gtfsrt"
	}
	ensureCanWrite(outputPath)
	fmt.Println(fmt.Sprintf("Output path: %s", outputPath))

	periodicity := getPeriodicity()
	fmt.Println(fmt.Sprintf("Update periodcity: %d", periodicity))
	fmt.Println("Ready")

	data.initialize()
	for {
		run(&data, outputPath)
		time.Sleep(time.Duration(periodicity) * time.Millisecond)
	}
}
