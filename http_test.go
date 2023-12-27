package pathgtfsrt

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/go-cmp/cmp"
	sourceapi "github.com/jamespfennell/path-train-gtfs-realtime/proto/sourceapi"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestSourceHttpGetTrainsAtStation(t *testing.T) {
	for _, tc := range []struct {
		testName     string
		station      sourceapi.Station
		jsonFilePath string
		trains       []Train
	}{
		{
			testName:     "Hoboken",
			station:      sourceapi.Station_HOBOKEN,
			jsonFilePath: "mock_data/source_http_hoboken.json",
			trains: []Train{
				{
					Route:            sourceapi.Route_JSQ_33_HOB,
					Direction:        sourceapi.Direction_TO_NY,
					LineName:         "33rd Street via Hoboken",
					ProjectedArrival: mkTimestampFromRfc3339("2023-12-23T05:36:15Z"),
					LastUpdated:      mkTimestampFromRfc3339("2023-12-23T05:35:44Z"),
				},
				{
					Route:            sourceapi.Route_JSQ_33_HOB,
					Direction:        sourceapi.Direction_TO_NY,
					LineName:         "33rd Street via Hoboken",
					ProjectedArrival: mkTimestampFromRfc3339("2023-12-23T06:01:30Z"),
					LastUpdated:      mkTimestampFromRfc3339("2023-12-23T05:35:44Z"),
				},
				{
					Route:            sourceapi.Route_JSQ_33_HOB,
					Direction:        sourceapi.Direction_TO_NJ,
					LineName:         "Journal Square via Hoboken",
					ProjectedArrival: mkTimestampFromRfc3339("2023-12-23T05:36:15Z"),
					LastUpdated:      mkTimestampFromRfc3339("2023-12-23T05:35:44Z"),
				},
				{
					Route:            sourceapi.Route_JSQ_33_HOB,
					Direction:        sourceapi.Direction_TO_NJ,
					LineName:         "Journal Square via Hoboken",
					ProjectedArrival: mkTimestampFromRfc3339("2023-12-23T06:02:44Z"),
					LastUpdated:      mkTimestampFromRfc3339("2023-12-23T05:35:44Z"),
				},
			},
		},
		{
			// See: https://github.com/mrazza/path-data/issues/22
			testName:     "Handle incorrect route mapping of JSQ_33 route",
			station:      sourceapi.Station_NEWPORT,
			jsonFilePath: "mock_data/source_http_newport.json",
			trains: []Train{
				{
					Route:            sourceapi.Route_JSQ_33,
					Direction:        sourceapi.Direction_TO_NY,
					LineName:         "33rd Street",
					ProjectedArrival: mkTimestampFromRfc3339("2023-12-27T00:09:21Z"),
					LastUpdated:      mkTimestampFromRfc3339("2023-12-27T00:01:24Z"),
				},
			},
		},
	} {
		mockHttpClient := MockSourceHTTPClient{StationToJsonFilePath: map[sourceapi.Station]string{tc.station: tc.jsonFilePath}}
		client := NewHttpSourceClient(mockHttpClient)
		ctx := context.Background()
		gotTrains, err := client.GetTrainsAtStation(ctx, tc.station)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if diff := cmp.Diff(&gotTrains, &tc.trains, protocmp.Transform()); diff != "" {
			t.Errorf("GTFS realtime feed got != want, diff=%s", diff)
		}
	}
}

func mkTimestampFromRfc3339(timeString string) *timestamp.Timestamp {
	timeObj, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return nil
	}
	value := timestamp.Timestamp{Seconds: timeObj.Unix()}
	return &value
}

func getStationNameFromURL(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	segments := strings.Split(parsedURL.Path, "/")
	// The segments slice will contain ["", "v1", "stations", "newark", "realtime", ""]
	// So, the station name is the fourth element (index 3)
	if len(segments) > 4 {
		return segments[3], nil
	}
	return "", fmt.Errorf("invalid URL format")
}

type MockSourceHTTPClient struct {
	JSONFilePath          string
	StationToJsonFilePath map[sourceapi.Station]string
}

func (m MockSourceHTTPClient) Get(reqUrl string) (*http.Response, error) {
	statioName, err := getStationNameFromURL(reqUrl)
	if err != nil {
		return nil, err
	}

	station := sourceapi.Station(sourceapi.Station_value[statioName])
	jsonFilePath, ok := m.StationToJsonFilePath[station]
	if !ok {
		return nil, fmt.Errorf("no JSON file path found for station %s", station)
	}

	// Read data from the specified JSON file
	file, err := os.Open(jsonFilePath)
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
