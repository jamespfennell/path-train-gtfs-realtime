package pathgtfsrt

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/golang/protobuf/ptypes/timestamp"
	sourceapi "github.com/jamespfennell/path-train-gtfs-realtime/proto/sourceapi"
)

const (
	paNyNjApiUrl      = "https://www.panynj.gov/bin/portauthority/ridepath.json"
	cacheValidityTime = 10 * time.Second
)

type cachedContent struct {
	timestamp time.Time
	data      []byte
	error     error
}

// PaNyNjClient is a source client that gets data from the Port Authority of New York and New Jersey.
// It is what is used to power the official realtime schedules on the PATH website: https://www.panynj.gov/path/en/index.html
type PaNyNjClient struct {
	httpClient    HttpClient
	clock         clock.Clock
	cachedContent *cachedContent
	mu            sync.RWMutex
}

var panynjStationToSourceStation = map[string]sourceapi.Station{
	"NWK": sourceapi.Station_NEWARK,
	"HAR": sourceapi.Station_HARRISON,
	"JSQ": sourceapi.Station_JOURNAL_SQUARE,
	"GRV": sourceapi.Station_GROVE_STREET,
	"NEW": sourceapi.Station_NEWPORT,
	"EXP": sourceapi.Station_EXCHANGE_PLACE,
	"HOB": sourceapi.Station_HOBOKEN,
	"WTC": sourceapi.Station_WORLD_TRADE_CENTER,
	"CHR": sourceapi.Station_CHRISTOPHER_STREET,
	"09S": sourceapi.Station_NINTH_STREET,
	"14S": sourceapi.Station_FOURTEENTH_STREET,
	"23S": sourceapi.Station_TWENTY_THIRD_STREET,
	"33S": sourceapi.Station_THIRTY_THIRD_STREET,
}

var panynjLineColorToRoute = map[string]sourceapi.Route{
	"4D92FB":        sourceapi.Route_HOB_33,
	"4D92FB,FF9900": sourceapi.Route_JSQ_33_HOB,
	"65C100":        sourceapi.Route_HOB_WTC,
	"FF9900":        sourceapi.Route_JSQ_33,
	"D93A30":        sourceapi.Route_NWK_WTC,
}

var panynjLabelToDirection = map[string]sourceapi.Direction{
	"TONY": sourceapi.Direction_TO_NY,
	"TONJ": sourceapi.Direction_TO_NJ,
}

// RidePathResponse contains information about all incoming trains at all stations
type RidePathResponse struct {
	Results []Result `json:"results"`
}

// Result represents all destinations for a given station
type Result struct {
	ConsideredStation string        `json:"consideredStation"`
	Destinations      []Destination `json:"destinations"`
}

// Destination is a labeled direction with associated trains
type Destination struct {
	Label    string    `json:"label"`
	Messages []Message `json:"messages"`
}

// Message contains information about a single train
type Message struct {
	Target             string `json:"target"`
	SecondsToArrival   string `json:"secondsToArrival"`
	ArrivalTimeMessage string `json:"arrivalTimeMessage"`
	LineColor          string `json:"lineColor"`
	HeadSign           string `json:"headSign"`
	LastUpdated        string `json:"lastUpdated"`
}

func NewPaNyNjSourceClient(httpClient HttpClient, clock clock.Clock) *PaNyNjClient {
	return &PaNyNjClient{httpClient: httpClient, clock: clock}
}

func (client *PaNyNjClient) GetTrainsAtStation(_ context.Context, station sourceapi.Station) ([]Train, error) {
	realtimeApiContent, err := client.getContent()
	if err != nil {
		return nil, err
	}
	response := RidePathResponse{}
	err = json.Unmarshal(realtimeApiContent, &response)
	if err != nil {
		return nil, err
	}
	var trains []Train
	for _, result := range response.Results {
		consideredStation := client.convertStationAsStringToStation(result.ConsideredStation)
		if consideredStation != station {
			continue
		}
		for _, destination := range result.Destinations {
			for _, message := range destination.Messages {
				lastUpdated, err := client.convertApiLastUpdatedTimeStringToTimestamp(message.LastUpdated)
				if err != nil {
					return nil, err
				}
				upcomingTrain := sourceapi.GetUpcomingTrainsResponse_UpcomingTrain{
					Route:            client.convertLineColorToRoute(message.LineColor),
					Direction:        client.convertDirectionAsStringToDirection(destination.Label),
					ProjectedArrival: client.convertApiSecondsToArrivalAsStringToTimestamp(lastUpdated, message.SecondsToArrival),
					LastUpdated:      lastUpdated,
				}
				trains = append(trains, &upcomingTrain)
			}
		}
	}

	return trains, nil
}

func (client *PaNyNjClient) GetStationToStopId(_ context.Context) (map[sourceapi.Station]string, error) {
	return sourceStationToGtfsStopId, nil
}

func (client *PaNyNjClient) GetRouteToRouteId(_ context.Context) (map[sourceapi.Route]string, error) {
	return sourceRouteToGtfsRouteId, nil
}

func (client *PaNyNjClient) convertDirectionAsStringToDirection(directionAsString string) sourceapi.Direction {
	direction, ok := panynjLabelToDirection[strings.ToUpper(directionAsString)]
	if !ok {
		return sourceapi.Direction_DIRECTION_UNSPECIFIED
	}
	return direction
}

func (client *PaNyNjClient) convertStationAsStringToStation(stationAsString string) sourceapi.Station {
	station, ok := panynjStationToSourceStation[stationAsString]
	if !ok {
		return sourceapi.Station_STATION_UNSPECIFIED
	}
	return station
}

func (client *PaNyNjClient) convertLineColorToRoute(lineColor string) sourceapi.Route {
	route, ok := panynjLineColorToRoute[strings.ToUpper(lineColor)]
	if !ok {
		return sourceapi.Route_ROUTE_UNSPECIFIED
	}
	return route
}

func (client *PaNyNjClient) convertApiLastUpdatedTimeStringToTimestamp(timeString string) (*timestamp.Timestamp, error) {
	const layout = "2006-01-02T15:04:05.999999-07:00"
	timeObj, err := time.Parse(layout, timeString)
	if err != nil {
		return nil, err
	}
	value := timestamp.Timestamp{Seconds: timeObj.Unix()}
	return &value, nil
}

func (client *PaNyNjClient) convertApiSecondsToArrivalAsStringToTimestamp(lastUpdated *timestamp.Timestamp, secondsToArrivalAsString string) *timestamp.Timestamp {
	secondsToArrival, err := strconv.ParseInt(secondsToArrivalAsString, 10, 64)
	if err != nil {
		return nil
	}
	return &timestamp.Timestamp{Seconds: lastUpdated.Seconds + secondsToArrival}
}

// Get the raw bytes from an endpoint in the API.
func (client *PaNyNjClient) getContent() (bytes []byte, err error) {
	client.mu.RLock()
	cachedData, err, ok := client.getCachedContent()
	if ok {
		defer client.mu.RUnlock()
		return cachedData, err
	}
	client.mu.RUnlock()

	client.mu.Lock()
	defer client.mu.Unlock()

	// Double check that the cache wasn't updated while we were waiting for the lock
	cachedData, err, ok = client.getCachedContent()
	if ok {
		return cachedData, err
	}

	url := attachTimestampToUrl(paNyNjApiUrl, client.clock)
	resp, err := client.httpClient.Get(url)
	if err != nil {
		client.cachedContent = &cachedContent{timestamp: client.clock.Now(), data: nil, error: err}
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		client.cachedContent = &cachedContent{timestamp: client.clock.Now(), data: nil, error: err}
		return nil, err
	}

	client.cachedContent = &cachedContent{timestamp: client.clock.Now(), data: data, error: nil}
	return data, nil
}

func (client *PaNyNjClient) getCachedContent() (cachedContent []byte, err error, ok bool) {
	if client.cachedContent != nil && client.clock.Now().Sub(client.cachedContent.timestamp) < cacheValidityTime {
		return client.cachedContent.data, client.cachedContent.error, true
	}
	return nil, nil, false
}

func attachTimestampToUrl(url string, clock clock.Clock) string {
	return url + "?timeStamp=" + strconv.FormatInt(clock.Now().Unix()*1000, 10)
}
