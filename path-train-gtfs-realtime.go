package main

import (
	"encoding/json"
	"fmt"
	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"net/http"
	"strings"
)

var apiStopToStopId = map[string]string{

}

const apiUrlRoutes = "https://path.api.razza.dev/v1/routes/"
const apiUrlStations = "https://path.api.razza.dev/v1/stations/"

// TODO: see if this needs to be flipped
var apiDirectionToDirectionId = map[string]uint32{
	"TO_NY": uint32(0),
	"TO_NJ": uint32(1),
}
var myJsonString = `
{
"upcomingTrains": [
{
"lineName": "World Trade Center",
"lineColors": [
"#D93A30"
],
"projectedArrival": "2020-05-09T01:41:50Z",
"lastUpdated": "2020-05-09T01:38:47Z",
"status": "ON_TIME",
"headsign": "World Trade Center",
"route": "NWK_WTC",
"routeDisplayName": "Newark - World Trade Center",
"direction": "TO_NY"
},
{
"lineName": "33rd Street",
"lineColors": [
"#FF9900"
],
"projectedArrival": "2020-05-09T01:43:15Z",
"lastUpdated": "2020-05-09T01:38:47Z",
"status": "ON_TIME",
"headsign": "33rd Street",
"route": "JSQ_33",
"routeDisplayName": "Journal Square - 33rd Street",
"direction": "TO_NY"
},
{
"lineName": "Newark",
"lineColors": [
"#D93A30"
],
"projectedArrival": "2020-05-09T01:38:47Z",
"lastUpdated": "2020-05-09T01:38:47Z",
"status": "ARRIVING_NOW",
"headsign": "Newark",
"route": "NWK_WTC",
"routeDisplayName": "World Trade Center - Newark",
"direction": "TO_NJ"
},
{
"lineName": "Journal Square",
"lineColors": [
"#FF9900"
],
"projectedArrival": "2020-05-09T01:43:32Z",
"lastUpdated": "2020-05-09T01:38:47Z",
"status": "ON_TIME",
"headsign": "Journal Square",
"route": "JSQ_33",
"routeDisplayName": "33rd Street - Journal Square",
"direction": "TO_NJ"
}
]
}
`


var routeJson = `
{
  "routes": [
    {
      "route": "JSQ_33_HOB",
      "id": "1024",
      "name": "Journal Square - 33rd Street (via Hoboken)",
      "color": "ff9900",
      "lines": [
        {
          "displayName": "33rd Street (via Hoboken) - Journal Square",
          "headsign": "Journal Square via Hoboken",
          "direction": "TO_NJ"
        },
        {
          "displayName": "Journal Square - 33rd Street (via Hoboken)",
          "headsign": "33rd via Hoboken",
          "direction": "TO_NY"
        }
      ]
    }
  ]
}
`



type apiTrain struct {
	ProjectedArrival string
	LastUpdated      string
	Route            string
	Direction        string
	Stop             string
}

type apiRealtimeResponse struct {
	Trains []apiTrain `json:"upcomingTrains"`
}
// https://path.api.razza.dev/v1/routes/

func getApiContent(url string) []byte {
	resp, _ := http.Get(url) // TODO: handle error
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

func buildApiRouteIdToRouteId(routeApiContent []byte) (apiRouteIdToRouteId map[string]string) {
	type apiRoute struct {
		ApiId string `json:"route"`
		Id string
	}
	type apiRoutesResponse struct {
		Routes []apiRoute `json:"routes"`
	}
	r := apiRoutesResponse{}
	json.Unmarshal(routeApiContent, &r)  // TODO: handle error
	apiRouteIdToRouteId = map[string]string{}
	for _, apiRoute:= range r.Routes {
		apiRouteIdToRouteId[apiRoute.ApiId] = apiRoute.Id
	}
	return
}


func buildApiStopIdToStopId(stationApiContent []byte) (apiStopIdToStopId map[string]string) {
	type apiStation struct {
		ApiId string `json:"station"`
		Id string
	}
	type apiStationsResponse struct {
		Stations []apiStation `json:"stations"`
	}
	r := apiStationsResponse{}
	json.Unmarshal(stationApiContent, &r)  // TODO: handle error
	apiStopIdToStopId = map[string]string{}
	for _, station:= range r.Stations {
		apiStopIdToStopId[strings.ToLower(station.ApiId)] = station.Id
	}
	fmt.Println(apiStopIdToStopId)
	return
}




func main() {

	r := apiRealtimeResponse{}
	err := json.Unmarshal([]byte(myJsonString), &r)
	r.Trains[0].Stop = "a"
	if err != nil {
		return
	}
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

	fmt.Println("hello world")
	fmt.Println(r)
	fmt.Println(a)
	out, err := proto.Marshal(&a)
	fmt.Println(out)
	// body := getApiContent("https://path.api.razza.dev/v1/routes/") // TODO: handle error
	// fmt.Println(buildApiRouteIdToRouteId(body))
	// body := getApiContent(apiUrlStations)
	// buildApiStopIdToStopId(body)
}

func convertApiTrainToTripUpdate(train apiTrain) gtfs.TripUpdate {
	lastUpdated := uint64(convertApiTimeStringToTimestamp(train.LastUpdated))
	arrivalTime := convertApiTimeStringToTimestamp(train.ProjectedArrival)
	stopId := apiStopToStopId[train.Stop]
	routeId := "a" //apiRouteToRouteId[train.Route]
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
