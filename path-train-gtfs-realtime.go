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
	"encoding/json"
	"fmt"
	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	exitCodeMalformedEnvVar         = 101
	exitCodeNoRoutesResponse        = 102
	exitCodeMalformedRoutesResponse = 103
	exitCodeNoStopsResponse         = 104
	exitCodeMalformedStopsResponse  = 105
	exitCodeCannotWriteToDisk       = 106
	envVarPeriodicity               = "PATH_GTFS_RT_PERIODICITY_MILLISECS"
	envVarOutputPath                = "PATH_GTFS_RT_OUTPUT_PATH"
	apiBaseUrl                      = "https://path.api.razza.dev/v1/"
	apiRoutesEndpoint               = "routes/"
	apiStationsEndpoint             = "stations/"
	apiRealtimeEndpoint             = "stations/%s/realtime/"
)

// (1)

// Data structure representing train data in the API.
type apiTrain struct {
	ProjectedArrival string
	LastUpdated      string
	Route            string
	Direction        string
}

// Data structure for communicating the full set of relevant data in the Stations endpoint.
// This structure just exists for using Go channels.
type apiTrainsAtStation struct {
	ApiTrains    []apiTrain
	ApiStationId string
	Err          error
}

// A container for the most updated data retrieved from the API so far.
// During the program execution, there is one instance of this structure.
type apiData struct {
	apiStationIdToStopId    map[string]string
	apiRouteIdToRouteId     map[string]string
	apiStationIdToApiTrains map[string][]apiTrain
}

// Initialize the apiData data structure by populating its stop and route fields.
// If the initialization fails, the program exits.
func (data *apiData) initialize() {
	routesContent, err := getApiContent(apiRoutesEndpoint)
	if err != nil {
		os.Exit(exitCodeNoRoutesResponse)
	}
	data.apiRouteIdToRouteId, err = buildApiRouteIdToRouteId(routesContent)
	if err != nil {
		os.Exit(exitCodeMalformedRoutesResponse)
	}
	stationsContent, err := getApiContent(apiStationsEndpoint)
	if err != nil {
		os.Exit(exitCodeNoStopsResponse)
	}
	data.apiStationIdToStopId, err = buildApiStationIdToStopId(stationsContent)
	if err != nil {
		os.Exit(exitCodeMalformedStopsResponse)
	}
	data.apiStationIdToApiTrains = map[string][]apiTrain{}
	for apiStationId := range data.apiStationIdToStopId {
		data.apiStationIdToApiTrains[apiStationId] = []apiTrain{}
	}
}

// Update the apiData structure using the most recent stations data from the API.
// If data for one of the stations cannot be retrieved, the old data will be preserved.
// This method is highly I/O bound: the time it takes to execute is largely spent waiting for HTTP responses
// from the 14 station endpoints.
// To speed it up, the 14 endpoints are hit in parallel.
func (data apiData) update() (err error) {
	updateResults := make(chan apiTrainsAtStation, len(data.apiStationIdToStopId))
	for apiStationId := range data.apiStationIdToStopId {
		apiStationId := apiStationId
		go func() { updateResults <- getTrainsAtStation(apiStationId) }()
	}
	for range data.apiStationIdToStopId {
		updateResult := <-updateResults
		if updateResult.Err == nil {
			data.apiStationIdToApiTrains[updateResult.ApiStationId] = updateResult.ApiTrains
		} else {
			err = updateResult.Err
		}
	}
	return
}

// Convert directions as described in the API ("TO_NY"/"TO_NJ") to boolean direction IDs as
// they appear in the GTFS Static feed.
func (data apiData) convertApiDirectionToDirectionId(apiDirection string) (directionId uint32) {
	if apiDirection == "TO_NY" {
		directionId = 1
	} else {
		directionId = 0
	}
	return
}

// (2)

// Get the raw bytes from an endpoint in the API.
func getApiContent(endpoint string) (bytes []byte, err error) {
	// TODO: timeout?
	resp, err := http.Get(apiBaseUrl + endpoint)
	if err != nil {
		return
	}
	// TODO: handle error properly
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// Get realtime data about a station from the API.
func getTrainsAtStation(apiStationId string) (result apiTrainsAtStation) {
	result = apiTrainsAtStation{ApiStationId: apiStationId}
	type apiRealtimeResponse struct {
		Trains []apiTrain `json:"upcomingTrains"`
	}
	realtimeApiContent, err := getApiContent(fmt.Sprintf(apiRealtimeEndpoint, apiStationId))
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

// Given the raw response from the list stations endpoint of the API, construct
// a map of the station ID in the API to the stop ID as it appears in the GTFS Static feed.
func buildApiStationIdToStopId(stationApiContent []byte) (apiStationIdToStopId map[string]string, err error) {
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
	apiStationIdToStopId = map[string]string{}
	for _, station := range response.Stations {
		apiStationIdToStopId[strings.ToLower(station.ApiId)] = station.Id
	}
	return
}

// Given the raw response from the list routes endpoint of the API, construct
// a map of the route ID in the API to the route ID as it appears in the GTFS Static feed.
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

// (3)

// Build a GTFS Realtime message from a snapshot of the current data.
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
	for apiStationId, trains := range data.apiStationIdToApiTrains {
		for _, train := range trains {
			tripUuid, err := uuid.NewRandom()
			if err != nil {
				continue
			}
			tripId := tripUuid.String()
			tripUpdate, err := convertApiTrainToTripUpdate(data, train, tripId, apiStationId)
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
	stopId := data.apiStationIdToStopId[apiStationId]
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
		os.Exit(exitCodeMalformedEnvVar)
	}
	return periodicity
}

func ensureCanWrite(outputPath string) {
	err := ioutil.WriteFile(outputPath, []byte{}, 0644)
	if err != nil {
		fmt.Println(fmt.Sprintf("Unable to write to '%s', exiting.", outputPath))
		os.Exit(exitCodeCannotWriteToDisk)
	}
}

// Run one feed update iteration.
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

func main() {
	fmt.Println("Starting up.")
	data := apiData{}
	data.initialize()
	outputPath, envVarSet := os.LookupEnv(envVarOutputPath)
	if !envVarSet {
		outputPath = "path.gtfsrt"
	}
	ensureCanWrite(outputPath)
	fmt.Println(fmt.Sprintf("Feed will be written to '%s'.", outputPath))
	periodicity := getPeriodicity()
	fmt.Println(fmt.Sprintf("Feed will be updated every %d milliseconds.", periodicity))
	fmt.Println("Ready.")
	for {
		run(data, outputPath)
		time.Sleep(time.Duration(periodicity) * time.Millisecond)
	}
}
