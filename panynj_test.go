package pathgtfsrt

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/go-cmp/cmp"
	sourceapi "github.com/jamespfennell/path-train-gtfs-realtime/proto/sourceapi"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestGetTrainsAtStation(t *testing.T) {
	c := clock.NewMock()
	for _, tc := range []struct {
		station sourceapi.Station
		trains  []Train
	}{
		{
			station: sourceapi.Station_NEWARK,
			trains: []Train{
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 232),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:42:07.827997-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 832),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:42:07.827997-05:00"),
				},
			},
		},
		{
			station: sourceapi.Station_HARRISON,
			trains: []Train{
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 217),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:27.941258-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 1088),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:27.941258-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 344),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:57.869032-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 944),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:57.869032-05:00"),
				},
			},
		},
		{
			station: sourceapi.Station_JOURNAL_SQUARE,
			trains: []Train{
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 408),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:47.905034-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 834),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:47.905034-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 367),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:52.813168-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 374),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:52.813168-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 1009),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:52.813168-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 1087),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:52.813168-05:00"),
				},
			},
		},
		{
			station: sourceapi.Station_GROVE_STREET,
			trains: []Train{
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 103),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:52.813168-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 446),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:52.813168-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 204),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:32.933609-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 669),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:32.933609-05:00"),
				},
			},
		},
		{
			station: sourceapi.Station_NEWPORT,
			trains: []Train{
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 231),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:27.941258-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_WTC,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 614),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:27.941258-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 44),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:27.941258-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_WTC,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 674),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:27.941258-05:00"),
				},
			},
		},
		{
			station: sourceapi.Station_EXCHANGE_PLACE,
			trains: []Train{
				{
					Route:            sourceapi.Route_HOB_WTC,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 299),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:42.854056-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 384),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:42.854056-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_WTC,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 239),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:42.854056-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 374),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:42.854056-05:00"),
				},
			},
		},
		{
			station: sourceapi.Station_HOBOKEN,
			trains: []Train{
				{
					Route:            sourceapi.Route_HOB_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 32),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:27.941258-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_WTC,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 452),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:27.941258-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 932),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:27.941258-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_WTC,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 1172),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:27.941258-05:00"),
				},
			},
		},
		{
			station: sourceapi.Station_WORLD_TRADE_CENTER,
			trains: []Train{
				{
					Route:            sourceapi.Route_HOB_WTC,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 47),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:42:12.868217-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 107),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:42:12.868217-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 707),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:42:12.868217-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_WTC,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 767),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:42:12.868217-05:00"),
				},
			},
		},
		{
			station: sourceapi.Station_CHRISTOPHER_STREET,
			trains: []Train{
				{
					Route:            sourceapi.Route_HOB_33,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 380),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:32.933609-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 465),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:32.933609-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 536),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:42.854056-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 621),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:42.854056-05:00"),
				},
			},
		},
		{
			station: sourceapi.Station_NINTH_STREET,
			trains: []Train{
				{
					Route:            sourceapi.Route_HOB_33,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 299),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:32.933609-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 384),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:32.933609-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 634),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:42:07.827997-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 719),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:42:07.827997-05:00"),
				},
			},
		},
		{
			station: sourceapi.Station_FOURTEENTH_STREET,
			trains: []Train{
				{
					Route:            sourceapi.Route_HOB_33,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 180),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:57.869032-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 265),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:57.869032-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33_HOB,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 2512),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:57.869032-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 22),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:52.813168-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 709),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:52.813168-05:00"),
				},
			},
		},
		{
			station: sourceapi.Station_TWENTY_THIRD_STREET,
			trains: []Train{
				{
					Route:            sourceapi.Route_HOB_33,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 121),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:27.941258-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 206),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:27.941258-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 80),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:47.905034-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromOffset(c, 165),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:47.905034-05:00"),
				},
			},
		},
		{
			station: sourceapi.Station_THIRTY_THIRD_STREET,
			trains: []Train{
				{
					Route:            sourceapi.Route_HOB_33,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 0),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:42:07.827997-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 0),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:42:07.827997-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_33,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 592),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:42:07.827997-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromOffset(c, 712),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:42:07.827997-05:00"),
				},
			},
		},
	} {
		client := NewClientWithMockedHttp(nil, c)
		ctx := context.Background()
		gotTrains, err := client.GetTrainsAtStation(ctx, tc.station)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if diff := cmp.Diff(&gotTrains, &tc.trains, protocmp.Transform()); diff != "" {
			t.Errorf("GTFS realtime feed got != want, diff=%s", diff)
		}

		// Ensure that the stationToStopId map contains the station
		stationToStopId, err := client.GetStationToStopId(ctx)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if _, ok := stationToStopId[tc.station]; !ok {
			t.Errorf("stationToStopId does not contain station %v", tc.station)
		}

		// Ensure that the routeToRouteId map contains the routes
		routeToRouteId, err := client.GetRouteToRouteId(ctx)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		for _, train := range tc.trains {
			if _, ok := routeToRouteId[train.Route]; !ok {
				t.Errorf("routeToRouteId does not contain route %v", train.Route)
			}
		}
	}
}

func TestGetStationToStopId(t *testing.T) {
	client := NewClientWithMockedHttp(nil, clock.New())
	ctx := context.Background()
	stationToStopId, err := client.GetStationToStopId(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !MapsEqual(sourceStationToGtfsStopId, stationToStopId) {
		t.Errorf("stationToStopId got=%v, want=%v", sourceStationToGtfsStopId, stationToStopId)
	}
	// A station should map to a unique stop ID
	seenStopIds := make(map[string]bool)
	for _, stopId := range stationToStopId {
		if seenStopIds[stopId] {
			t.Errorf("stationToStopId has duplicate stop ID %s", stopId)
		}
		seenStopIds[stopId] = true
	}

}

func TestGetRouteToRouteId(t *testing.T) {
	client := NewClientWithMockedHttp(nil, clock.New())
	ctx := context.Background()
	routeToRouteId, err := client.GetRouteToRouteId(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !MapsEqual(sourceRouteToGtfsRouteId, routeToRouteId) {
		t.Errorf("routeToRouteId got=%v, want=%v", sourceRouteToGtfsRouteId, routeToRouteId)
	}
	// A route should map to a unique route ID
	seenRouteIds := make(map[string]bool)
	for _, routeId := range routeToRouteId {
		if seenRouteIds[routeId] {
			t.Errorf("routeToRouteId has duplicate route ID %s", routeId)
		}
		seenRouteIds[routeId] = true
	}
}

func NewClientWithMockedHttp(jsonFilePath *string, clock clock.Clock) *PaNyNjClient {
	var jsonFilePathString string
	if jsonFilePath == nil {
		jsonFilePathString = "mock_data/ridepath.json"
	} else {
		jsonFilePathString = *jsonFilePath
	}
	mockHttp := MockHTTPClient{
		JSONFilePath: jsonFilePathString,
		clock:        clock,
	}
	return NewPaNyNjSourceClient(mockHttp, clock)
}

type MockHTTPClient struct {
	JSONFilePath string
	clock        clock.Clock
}

func (m MockHTTPClient) Get(reqUrl string) (*http.Response, error) {
	// Validate URL and timeStamp query parameter
	parsedURL, err := url.Parse(reqUrl)
	if err != nil {
		return nil, err
	}
	params := parsedURL.Query()
	if params.Get("timeStamp") == "" {
		return nil, errors.New("timeStamp query parameter is required")
	}
	if !isValidMillisecondUnixTimestamp(params.Get("timeStamp"), m.clock) {
		return nil, errors.New("timeStamp query parameter must be a valid millisecond unix timestamp")
	}

	// Read data from the specified JSON file
	file, err := os.Open(m.JSONFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Mock the response
	r := ioutil.NopCloser(bytes.NewReader(data))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil
}

func MapsEqual[K comparable, V comparable](map1, map2 map[K]V) bool {
	if len(map1) != len(map2) {
		return false
	}

	for key, value1 := range map1 {
		value2, ok := map2[key]
		if !ok || value1 != value2 {
			return false
		}
	}

	return true
}

func isValidMillisecondUnixTimestamp(str string, clock clock.Clock) bool {
	// Parse the string as an integer
	timestamp, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return false
	}

	// Convert the timestamp to a time.Time
	timeFromTimestamp := time.UnixMilli(timestamp)

	// Get the current time
	now := clock.Now()

	// Check if the timestamp is less than or equal to the current time
	// and not more than 10 seconds before now
	return !timeFromTimestamp.After(now) && now.Sub(timeFromTimestamp) <= 10*time.Second
}

func mkTimestampFromOffset(c clock.Clock, offsetSeconds int64) *timestamp.Timestamp {
	return &timestamp.Timestamp{Seconds: c.Now().Unix() + offsetSeconds}
}

func mkTimestampFromIso8601(timeString string) *timestamp.Timestamp {
	const layout = "2006-01-02T15:04:05.000000-07:00"
	timeObj, err := time.Parse(layout, timeString)
	if err != nil {
		return nil
	}
	value := timestamp.Timestamp{Seconds: timeObj.Unix()}
	return &value
}
