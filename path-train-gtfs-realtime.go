package main

import (
	"encoding/json"
	"fmt"
	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type apiTrain struct {
	ProjectedArrival string
	LastUpdated      string
	Route            string
	Direction        string
}

var apiStopIdToStopId = map[string]string{}
var apiRouteIdToRouteId = map[string]string{}
var apiStopIdToApiTrains = map[string][]apiTrain{}

const apiUrlRoutes = "https://path.api.razza.dev/v1/routes/"
const apiUrlStations = "https://path.api.razza.dev/v1/stations/"
const apiUrlRealtime = "https://path.api.razza.dev/v1/stations/%s/realtime/"

// TODO: see if this needs to be flipped by looking at the GTFS static
var apiDirectionToDirectionId = map[string]uint32{
	"TO_NY": uint32(0),
	"TO_NJ": uint32(1),
}

func updateTrainsAtStation(station string) (err error) {
	type apiRealtimeResponse struct {
		Trains []apiTrain `json:"upcomingTrains"`
	}
	realtimeApiContent, err := getApiContent(fmt.Sprintf(apiUrlRealtime, station))
	if err != nil {
		return
	}
	response := apiRealtimeResponse{}
	err = json.Unmarshal(realtimeApiContent, &response)
	if err != nil {
		return
	}
	apiStopIdToApiTrains[station] = response.Trains
	return
}

func updateTrainsAtAllStations() (err error) {
	for apiStopId := range apiStopIdToStopId {
		err = updateTrainsAtStation(apiStopId)
	}
	return
}

func getApiContent(url string) (bytes []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	// TODO: handle error properly
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
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
	apiStopIdToApiTrains["harrison"] = []apiTrain{}
	_ = updateTrainsAtStation("harrison")
	fmt.Println(apiStopIdToApiTrains)
	/*
		b := "str"
		c := gtfs.FeedHeader_FULL_DATASET
		t := uint64(405)
		a := gtfs.FeedMessage{
			Header: &gtfs.FeedHeader{
				GtfsRealtimeVersion: &b,
				Incrementality:      &c,
				Timestamp:           &t,
			},
		}

		// out, err := proto.Marshal(&a)
		// fmt.Println(out)
		initializeIdMaps()
		fmt.Println(apiStopIdToStopId)
		fmt.Println(apiRouteIdToRouteId)
	*/
}

func initializeApiIdMaps() {
	routesContent, err := getApiContent(apiUrlRoutes)
	if err != nil {
		os.Exit(1)
	}
	apiRouteIdToRouteId, err = buildApiRouteIdToRouteId(routesContent)
	if err != nil {
		os.Exit(2)
	}
	stationsContent, err := getApiContent(apiUrlStations)
	if err != nil {
		os.Exit(3)
	}
	apiStopIdToStopId, err = buildApiStopIdToStopId(stationsContent)
	if err != nil {
		os.Exit(4)
	}
	for apiStopId := range apiStopIdToStopId {
		apiStopIdToApiTrains[apiStopId] = []apiTrain{}
	}
}

func convertApiTrainToTripUpdate(train apiTrain, stopId string) gtfs.TripUpdate {
	lastUpdated := uint64(convertApiTimeStringToTimestamp(train.LastUpdated))
	arrivalTime := convertApiTimeStringToTimestamp(train.ProjectedArrival)
	routeId := apiRouteIdToRouteId[train.Route]
	directionId := apiDirectionToDirectionId[train.Direction]
	stopTimeUpdate := gtfs.TripUpdate_StopTimeUpdate{
		StopSequence: nil,
		StopId:       &stopId,
		Arrival: &gtfs.TripUpdate_StopTimeEvent{
			Time: &arrivalTime,
		},
	}
	return gtfs.TripUpdate{
		Trip: &gtfs.TripDescriptor{
			TripId:      nil,
			RouteId:     &routeId,
			DirectionId: &directionId,
		},
		StopTimeUpdate: []*gtfs.TripUpdate_StopTimeUpdate{
			&stopTimeUpdate,
		},
		Timestamp: &lastUpdated,
	}
}

func convertApiTimeStringToTimestamp(timeString string) int64 {
	return int64(4)
}
