package pathgtfsrt

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	sourceapi "github.com/jamespfennell/path-train-gtfs-realtime/proto/sourceapi"
)

const (
	apiBaseUrl          = "https://path.api.razza.dev/v1/"
	apiRoutesEndpoint   = "routes/"
	apiStationsEndpoint = "stations/"
	apiRealtimeEndpoint = "stations/%s/realtime/"
)

// HttpSourceClient is a source client that gets data using the Razza HTTP API.
type HttpSourceClient struct {
	httpClient    HttpClient
}

func NewHttpSourceClient(httpClient HttpClient) *HttpSourceClient {
	return &HttpSourceClient{httpClient: httpClient}
}

func (client *HttpSourceClient) GetTrainsAtStation(_ context.Context, station sourceapi.Station) ([]Train, error) {
	type jsonUpcomingTrain struct {
		ProjectedArrival  string
		LastUpdated       string
		RouteAsString     string `json:"route"`
		DirectionAsString string `json:"direction"`
		LineName          string `json:"lineName"`
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
			LineName:         rawUpcomingTrain.LineName,
			Direction:        client.convertDirectionAsStringToDirection(rawUpcomingTrain.DirectionAsString),
			ProjectedArrival: client.convertApiTimeStringToTimestamp(rawUpcomingTrain.ProjectedArrival),
			LastUpdated:      client.convertApiTimeStringToTimestamp(rawUpcomingTrain.LastUpdated),
		}
		applyRouteQaToTrain(&upcomingTrain)
		trains = append(trains, &upcomingTrain)
	}
	return trains, nil
}

func (client *HttpSourceClient) GetStationToStopId(_ context.Context) (map[sourceapi.Station]string, error) {
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

func (client *HttpSourceClient) GetRouteToRouteId(_ context.Context) (map[sourceapi.Route]string, error) {
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

func (client *HttpSourceClient) convertDirectionAsStringToDirection(directionAsString string) sourceapi.Direction {
	return sourceapi.Direction(sourceapi.Direction_value[directionAsString])
}

func (client *HttpSourceClient) convertStationAsStringToStation(stationAsString string) sourceapi.Station {
	return sourceapi.Station(sourceapi.Station_value[stationAsString])
}

func (client *HttpSourceClient) convertRouteAsStringToRoute(routeAsString string) sourceapi.Route {
	return sourceapi.Route(sourceapi.Route_value[routeAsString])
}
func (client *HttpSourceClient) convertApiTimeStringToTimestamp(timeString string) *timestamp.Timestamp {
	timeObj, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return nil
	}
	value := timestamp.Timestamp{Seconds: timeObj.Unix()}
	return &value
}

// Get the raw bytes from an endpoint in the API.
func (client HttpSourceClient) getContent(endpoint string) (bytes []byte, err error) {
	resp, err := client.httpClient.Get(apiBaseUrl + endpoint)
	if err != nil {
		return
	}
	defer func() {
		closingErr := resp.Body.Close()
		if err == nil {
			err = closingErr
		}
	}()
	return io.ReadAll(resp.Body)
}
