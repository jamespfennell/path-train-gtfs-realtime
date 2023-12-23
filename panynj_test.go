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
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950359),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:42:07.827997-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950959),
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
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950304),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:27.941258-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702951175),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:27.941258-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950461),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:57.869032-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702951061),
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
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950515),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:47.905034-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950941),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:47.905034-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950479),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:52.813168-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950486),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:52.813168-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702951121),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:52.813168-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702951199),
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
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950215),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:52.813168-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950558),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:52.813168-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950296),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:32.933609-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950761),
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
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950318),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:27.941258-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_WTC,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950701),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:27.941258-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950131),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:27.941258-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_WTC,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950761),
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
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950401),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:42.854056-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950486),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:42.854056-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_WTC,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950341),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:42.854056-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950476),
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
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950119),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:27.941258-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_WTC,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950539),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:27.941258-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702951019),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:27.941258-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_WTC,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702951259),
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
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950179),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:42:12.868217-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950239),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:42:12.868217-05:00"),
				},
				{
					Route:            sourceapi.Route_NWK_WTC,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950839),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:42:12.868217-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_WTC,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950899),
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
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950472),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:32.933609-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950557),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:32.933609-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950638),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:42.854056-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950723),
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
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950391),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:32.933609-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950476),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:32.933609-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950761),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:42:07.827997-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950846),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:42:07.827997-05:00"),
				},
			},
		},
		{
			station: sourceapi.Station_FOURTEENTH_STREET,
			trains:  GetFourteenthStreetTrains(c, 0),
		},
		{
			station: sourceapi.Station_TWENTY_THIRD_STREET,
			trains: []Train{
				{
					Route:            sourceapi.Route_HOB_33,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950208),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:27.941258-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950293),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:27.941258-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950187),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:41:47.905034-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_33,
					Direction:        sourceapi.Direction_TO_NY,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950272),
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
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950127),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:42:07.827997-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950127),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:42:07.827997-05:00"),
				},
				{
					Route:            sourceapi.Route_HOB_33,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950719),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:42:07.827997-05:00"),
				},
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NJ,
					ProjectedArrival: mkTimestampFromUnixSeconds(1702950839),
					LastUpdated:      mkTimestampFromIso8601("2023-12-18T20:42:07.827997-05:00"),
				},
			},
		},
	} {
		client, _ := NewClientWithMockedHttp(nil, c)
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

func TestResponseCaching(t *testing.T) {
	c := clock.NewMock()
	client, mockHttpClient := NewClientWithMockedHttp(nil, c)
	ctx := context.Background()
	expectedTrains := GetFourteenthStreetTrains(c, 0)
	gotTrains, err := client.GetTrainsAtStation(ctx, sourceapi.Station_FOURTEENTH_STREET)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if diff := cmp.Diff(&gotTrains, &expectedTrains, protocmp.Transform()); diff != "" {
		t.Errorf("GTFS realtime feed got != want, diff=%s", diff)
	}

	// Update response on the server
	mockHttpClient.JSONFilePath = "mock_data/ridepath_02.json"

	// Advance the clock by 1 second, cache should still be valid
	c.Add(1 * time.Second)
	gotTrains, err = client.GetTrainsAtStation(ctx, sourceapi.Station_FOURTEENTH_STREET)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if diff := cmp.Diff(&gotTrains, &expectedTrains, protocmp.Transform()); diff != "" {
		t.Errorf("GTFS realtime feed got != want, diff=%s", diff)
	}

	// Advance the clock by 10 seconds, cache should now be invalid
	c.Add(10 * time.Second)
	expectedTrains = GetFourteenthStreetTrains(c, 5)
	gotTrains, err = client.GetTrainsAtStation(ctx, sourceapi.Station_FOURTEENTH_STREET)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if diff := cmp.Diff(&gotTrains, &expectedTrains, protocmp.Transform()); diff != "" {
		t.Errorf("GTFS realtime feed got != want, diff=%s", diff)
	}
}

func TestGetStationToStopId(t *testing.T) {
	client, _ := NewClientWithMockedHttp(nil, clock.New())
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

// This test ensures that the client can handle fractional second precision
// Of 5 and 6 digits. The file 'mock_data/ridepath_03.json' contains a real
// response from the PATH API with 5 digit fractional second precision that
// was causing the client to crash previously.
func TestDataWith5DigitFractionalSecondPrecision(t *testing.T) {
	jsonPath := "mock_data/ridepath_03.json"
	client, _ := NewClientWithMockedHttp(&jsonPath, clock.New())
	ctx := context.Background()
	_, err := client.GetTrainsAtStation(ctx, sourceapi.Station_FOURTEENTH_STREET)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestGetRouteToRouteId(t *testing.T) {
	client, _ := NewClientWithMockedHttp(nil, clock.New())
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

func GetFourteenthStreetTrains(c clock.Clock, offset int64) []Train {
	return []Train{
		{
			Route:            sourceapi.Route_HOB_33,
			Direction:        sourceapi.Direction_TO_NJ,
			ProjectedArrival: mkTimestampFromUnixSeconds(1702950297 + offset),
			LastUpdated:      mkTimestampFromIso8601WithOffset("2023-12-18T20:41:57.869032-05:00", offset),
		},
		{
			Route:            sourceapi.Route_JSQ_33,
			Direction:        sourceapi.Direction_TO_NJ,
			ProjectedArrival: mkTimestampFromUnixSeconds(1702950382 + offset),
			LastUpdated:      mkTimestampFromIso8601WithOffset("2023-12-18T20:41:57.869032-05:00", offset),
		},
		{
			Route:            sourceapi.Route_JSQ_33_HOB,
			Direction:        sourceapi.Direction_TO_NJ,
			ProjectedArrival: mkTimestampFromUnixSeconds(1702952629 + offset),
			LastUpdated:      mkTimestampFromIso8601WithOffset("2023-12-18T20:41:57.869032-05:00", offset),
		},
		{
			Route:            sourceapi.Route_HOB_33,
			Direction:        sourceapi.Direction_TO_NY,
			ProjectedArrival: mkTimestampFromUnixSeconds(1702950134 + offset),
			LastUpdated:      mkTimestampFromIso8601WithOffset("2023-12-18T20:41:52.813168-05:00", offset),
		},
		{
			Route:            sourceapi.Route_HOB_33,
			Direction:        sourceapi.Direction_TO_NY,
			ProjectedArrival: mkTimestampFromUnixSeconds(1702950821 + offset),
			LastUpdated:      mkTimestampFromIso8601WithOffset("2023-12-18T20:41:52.813168-05:00", offset),
		},
	}
}

func NewClientWithMockedHttp(jsonFilePath *string, clock clock.Clock) (*PaNyNjClient, *MockHTTPClient) {
	var jsonFilePathString string
	if jsonFilePath == nil {
		jsonFilePathString = "mock_data/ridepath_01.json"
	} else {
		jsonFilePathString = *jsonFilePath
	}
	mockHttp := MockHTTPClient{
		JSONFilePath: jsonFilePathString,
		Clock:        clock,
	}
	return NewPaNyNjSourceClient(&mockHttp, clock), &mockHttp
}

type MockHTTPClient struct {
	JSONFilePath string
	Clock        clock.Clock
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
	if !isValidMillisecondUnixTimestamp(params.Get("timeStamp"), m.Clock) {
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

func mkTimestampFromUnixSeconds(seconds int64) *timestamp.Timestamp {
	return &timestamp.Timestamp{Seconds: seconds}
}

func mkTimestampFromIso8601(timeString string) *timestamp.Timestamp {
	return mkTimestampFromIso8601WithOffset(timeString, 0)
}

func mkTimestampFromIso8601WithOffset(timeString string, offset int64) *timestamp.Timestamp {
	const layout = "2006-01-02T15:04:05.000000-07:00"
	timeObj, err := time.Parse(layout, timeString)
	if err != nil {
		return nil
	}
	value := timestamp.Timestamp{Seconds: timeObj.Unix() + offset}
	return &value
}
