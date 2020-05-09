package main

import (
	"encoding/json"
	"fmt"
	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	apiUrlRoutes   = "https://path.api.razza.dev/v1/routes/"
	apiUrlStations = "https://path.api.razza.dev/v1/stations/"
	apiUrlRealtime = "https://path.api.razza.dev/v1/stations/%s/realtime/"
)

type apiTrain struct {
	ProjectedArrival string
	LastUpdated      string
	Route            string
	Direction        string
}

type apiTrainsAtStation struct {
	ApiTrains    []apiTrain
	ApiStationId string
	Err          error
}

type apiData struct {
	apiStopIdToStopId    map[string]string
	apiRouteIdToRouteId  map[string]string
	apiStopIdToApiTrains map[string][]apiTrain
}

func (data apiData) initialize() {
	routesContent, err := getApiContent(apiUrlRoutes)
	if err != nil {
		os.Exit(101)
	}
	data.apiRouteIdToRouteId, err = buildApiRouteIdToRouteId(routesContent)
	if err != nil {
		os.Exit(102)
	}
	stationsContent, err := getApiContent(apiUrlStations)
	if err != nil {
		os.Exit(103)
	}
	data.apiStopIdToStopId, err = buildApiStopIdToStopId(stationsContent)
	if err != nil {
		os.Exit(104)
	}
	data.apiStopIdToApiTrains = map[string][]apiTrain{}
	for apiStopId := range data.apiStopIdToStopId {
		data.apiStopIdToApiTrains[apiStopId] = []apiTrain{}
	}
}

func (data apiData) update() (err error) {
	updateResults := make(chan apiTrainsAtStation, len(data.apiStopIdToApiTrains))
	for apiStopId := range data.apiStopIdToStopId {
		apiStopId := apiStopId
		go func() { updateResults <- getTrainsAtStation(apiStopId) }()
	}
	for range data.apiStopIdToApiTrains {
		updateResult := <-updateResults
		if updateResult.Err == nil {
			data.apiStopIdToApiTrains[updateResult.ApiStationId] = updateResult.ApiTrains
		} else {
			err = updateResult.Err
		}
	}
	return
}

func (data apiData) convertApiDirectionToDirectionId(apiDirection string) (directionId uint32) {
	if apiDirection == "TO_NY" {
		directionId = 1
	} else {
		directionId = 0
	}
	return
}

func getApiContent(url string) (bytes []byte, err error) {
	// TODO: timeout?
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	// TODO: handle error properly
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func getTrainsAtStation(apiStationId string) (result apiTrainsAtStation) {
	result = apiTrainsAtStation{ApiStationId: apiStationId}
	type apiRealtimeResponse struct {
		Trains []apiTrain `json:"upcomingTrains"`
	}
	realtimeApiContent, err := getApiContent(fmt.Sprintf(apiUrlRealtime, apiStationId))
	if err != nil {
		result.Err = err
		return
	}
	response := apiRealtimeResponse{}
	err = json.Unmarshal(realtimeApiContent, &response)
	if err != nil {
		result.Err = err
		return
	}
	result.ApiTrains = response.Trains
	return
}

func buildApiRouteIdToRouteId(routeApiContent []byte) (apiRouteIdToRouteId map[string]string, err error) {
	type apiRoute struct {
		ApiId string `json:"route"`
		Id    string
	}
	type apiRoutesResponse struct {
		Routes []apiRoute `json:"routes"`
	}
	response := apiRoutesResponse{}
	err = json.Unmarshal(routeApiContent, &response)
	if err != nil {
		return
	}
	apiRouteIdToRouteId = map[string]string{}
	for _, apiRoute := range response.Routes {
		apiRouteIdToRouteId[apiRoute.ApiId] = apiRoute.Id
	}
	return
}

func buildGtfsRealtimeFeedMessage(data apiData) gtfs.FeedMessage {
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
	for apiStopId, trains := range data.apiStopIdToApiTrains {
		for _, train := range trains {
			tripUuid, err := uuid.NewRandom()
			if err != nil {
				continue
			}
			tripId := tripUuid.String()
			tripUpdate, err := convertApiTrainToTripUpdate(data, train, tripId, apiStopId)
			if err != nil {
				continue
			}
			feedEntity := gtfs.FeedEntity{
				Id:         &tripId,
				TripUpdate: &tripUpdate,
			}
			feedMessage.Entity = append(feedMessage.Entity, &feedEntity)
		}
	}
	return feedMessage
}

func convertApiTrainToTripUpdate(data apiData, train apiTrain, tripId string, apiStationId string) (update gtfs.TripUpdate, err error) {
	lastUpdated, err := convertApiTimeStringToTimestamp(train.LastUpdated)
	if err != nil {
		return
	}
	lastUpdatedUnsigned := uint64(lastUpdated)
	arrivalTime, err := convertApiTimeStringToTimestamp(train.ProjectedArrival)
	if err != nil {
		return
	}
	stopId := data.apiStopIdToStopId[apiStationId]
	routeId := data.apiRouteIdToRouteId[train.Route]
	directionId := data.convertApiDirectionToDirectionId(train.Direction)
	stopTimeUpdate := gtfs.TripUpdate_StopTimeUpdate{
		StopSequence: nil,
		StopId:       &stopId,
		Arrival: &gtfs.TripUpdate_StopTimeEvent{
			Time: &arrivalTime,
		},
	}
	return gtfs.TripUpdate{
		Trip: &gtfs.TripDescriptor{
			TripId:      &tripId,
			RouteId:     &routeId,
			DirectionId: &directionId,
		},
		StopTimeUpdate: []*gtfs.TripUpdate_StopTimeUpdate{
			&stopTimeUpdate,
		},
		Timestamp: &lastUpdatedUnsigned,
	}, nil
}

func convertApiTimeStringToTimestamp(timeString string) (t int64, err error) {
	timeObj, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return
	}
	t = timeObj.Unix()
	return
}

func buildApiStopIdToStopId(stationApiContent []byte) (apiStopIdToStopId map[string]string, err error) {
	type apiStation struct {
		ApiId string `json:"station"`
		Id    string
	}
	type apiStationsResponse struct {
		Stations []apiStation `json:"stations"`
	}
	response := apiStationsResponse{}
	err = json.Unmarshal(stationApiContent, &response)
	if err != nil {
		return
	}
	apiStopIdToStopId = map[string]string{}
	for _, station := range response.Stations {
		apiStopIdToStopId[strings.ToLower(station.ApiId)] = station.Id
	}
	return
}

func main() {
	fmt.Println("Starting up.")
	data := apiData{}
	data.initialize()
	outputPath, envVarSet := os.LookupEnv("PATH_GTFS_REALTIME_OUTPUT_PATH")
	if !envVarSet {
		outputPath = "path.gtfsrt"
	}
	fmt.Println(fmt.Sprintf("Feed will be written to '%s'.", outputPath))
	fmt.Println("Ready.")
	for {
		run(data, outputPath)
		time.Sleep(5 * time.Second)
	}
}

func run(data apiData, outputPath string) {
	fmt.Println("Updating GTFS Realtime feed.")
	err := data.update()
	if err != nil {
		fmt.Println("There was an error while retrieving the data; update will continue with some data stale.")
	}
	feedMessage := buildGtfsRealtimeFeedMessage(data)
	out, err := proto.Marshal(&feedMessage)
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
