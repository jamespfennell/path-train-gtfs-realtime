// Package pathgtfsrt contains a GTFS realtime feed generator for the PATH train.
package pathgtfsrt

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/benbjohnson/clock"
	gtfs "github.com/jamespfennell/path-train-gtfs-realtime/proto/gtfsrt"
	sourceapi "github.com/jamespfennell/path-train-gtfs-realtime/proto/sourceapi"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Train contains data about a PATH train at a specific station.
type Train *sourceapi.GetUpcomingTrainsResponse_UpcomingTrain

// SourceClient describes the methods that the feed generator requires from the source API in order to build the feed.
type SourceClient interface {
	// Return a map from source API station code to GTFS static stop ID
	GetStationToStopId(context.Context) (map[sourceapi.Station]string, error)
	// Return a map from source API route code to GTFS static route ID
	GetRouteToRouteId(context.Context) (map[sourceapi.Route]string, error)
	// List all upcoming trains at a station
	GetTrainsAtStation(context.Context, sourceapi.Station) ([]Train, error)
}

// Feed periodically generates GTFS Realtime data for the PATH train and makes
// it available through the `Get` method.
//
// Feed also satisfies the http.Handler interface, and simply responds to all requests with the most recent
// GTFS realtime data.
type Feed struct {
	gtfs  []byte
	mutex sync.RWMutex
}

// UpdateCallback is the type of callback that the feed runs after each update.
//
// The first argument is the GTFS realtime message that was just built.
// The second argument is the list of all errors that occured when getting realtime data
// from the source API.
type UpdateCallback func(msg *gtfs.FeedMessage, requestErrs []error)

// NewFeed creates a new feed.
//
// This function gets static and realtime data from the source API and creates the
// first version of the GTFS realtime feed before returning.
// It then, in the background, periodically updates the realtime data following the provided
// update period.
//
// After each update, including the first synchronous update, the provided callback is invoked.
func NewFeed(ctx context.Context, clock clock.Clock, updatePeriod time.Duration, sourceClient SourceClient, callback UpdateCallback) (*Feed, error) {
	f := Feed{}
	fmt.Println("Starting up")
	staticData, err := getStaticData(ctx, sourceClient)
	if err != nil {
		return nil, err
	}
	realtimeData := map[sourceapi.Station][]Train{}

	updateFunc := func() []error {
		fmt.Println("Updating GTFS Realtime feed.")
		requestErrs := updateRealtimeData(ctx, realtimeData, sourceClient, staticData)
		feedMessage := buildGtfsRealtimeFeedMessage(clock, staticData, realtimeData)
		out, err := proto.Marshal(feedMessage)
		if err != nil {
			panic(fmt.Sprintf("failed go generate realtime protobuf file: %s", err))
		}
		f.set(out)
		callback(feedMessage, requestErrs)
		fmt.Println("Finished updating")
		return requestErrs
	}

	errs := updateFunc()
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to initialize realtime data: %v", errs)
	}
	// We ensure the ticker is constructed before the function is returned; otherwise,
	// there is a race condition between initializing the ticker and incrementing the
	// time in the unit testing which results in a deadlock.
	ticker := clock.Ticker(updatePeriod)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				updateFunc()
			}
		}
	}()
	return &f, nil
}

// Get returns the most recent GTFS realtime data.
func (f *Feed) Get() []byte {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.gtfs
}

func (f *Feed) set(b []byte) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.gtfs = b
}

// ServeHTTP responds to all requests with the most recent GTFS realtime data.
func (f *Feed) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write(f.Get())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// A container for the static data retrieved at the start.
type staticData struct {
	stations        []sourceapi.Station
	stationToStopId map[sourceapi.Station]string
	routeToRouteId  map[sourceapi.Route]string
}

// Gets static data from the source API.
func getStaticData(ctx context.Context, sourceClient SourceClient) (staticData, error) {
	var s staticData
	var err error
	s.routeToRouteId, err = sourceClient.GetRouteToRouteId(ctx)
	if err != nil {
		return staticData{}, err
	}
	s.stationToStopId, err = sourceClient.GetStationToStopId(ctx)
	if err != nil {
		return staticData{}, err
	}
	for station := range s.stationToStopId {
		s.stations = append(s.stations, station)
	}
	sort.Slice(s.stations, func(i, j int) bool {
		return s.stations[i] < s.stations[j]
	})
	return s, nil
}

// Updates the realtime data using the source API.
//
// If data for one or more stations cannot be retrieved, the pre-existing realtime data is conservered
// and corresponding number of errors are returned.
func updateRealtimeData(ctx context.Context, data map[sourceapi.Station][]Train, sourceClient SourceClient, staticData staticData) []error {
	type trainsAtStation struct {
		Station sourceapi.Station
		Trains  []Train
		Err     error
	}
	allTrainsAtStations := make(chan trainsAtStation, len(staticData.stationToStopId))
	for station := range staticData.stationToStopId {
		station := station
		go func() {
			r := trainsAtStation{Station: station}
			r.Trains, r.Err = sourceClient.GetTrainsAtStation(ctx, station)
			allTrainsAtStations <- r
		}()
	}
	var errs []error
	for range staticData.stationToStopId {
		trainsAtStation := <-allTrainsAtStations
		if trainsAtStation.Err != nil {
			errs = append(errs, trainsAtStation.Err)
			fmt.Println("There was an error when retrieving data for station",
				staticData.stationToStopId[trainsAtStation.Station])
			continue
		}
		data[trainsAtStation.Station] = trainsAtStation.Trains
	}
	return errs
}

// Build a GTFS Realtime message from a snapshot of the current data.
func buildGtfsRealtimeFeedMessage(clock clock.Clock, staticData staticData, realtimeData map[sourceapi.Station][]Train) *gtfs.FeedMessage {
	directionToBoolean := func(direction sourceapi.Direction) *uint32 {
		var result uint32
		if direction == sourceapi.Direction_TO_NY {
			result = 1
		} else if direction == sourceapi.Direction_TO_NJ {
			result = 0
		}
		return &result
	}
	timestamppbToInt64 := func(t *timestamppb.Timestamp) *int64 {
		if t != nil {
			return ptr(t.Seconds)
		}
		return nil
	}
	timestamppbToUint64 := func(t *timestamppb.Timestamp) *uint64 {
		if t != nil {
			return ptr(uint64(t.Seconds))
		}
		return nil
	}
	var entities []*gtfs.FeedEntity
	for _, apiStationId := range staticData.stations {
		trains := realtimeData[apiStationId]
		for _, train := range trains {
			routeID, ok := staticData.routeToRouteId[train.Route]
			if !ok {
				continue
			}
			if train.Direction != sourceapi.Direction_TO_NJ && train.Direction != sourceapi.Direction_TO_NY {
				continue
			}
			if train.ProjectedArrival == nil {
				continue
			}
			if train.LastUpdated == nil {
				continue
			}
			update := &gtfs.TripUpdate{
				Trip: &gtfs.TripDescriptor{
					RouteId:     &routeID,
					DirectionId: directionToBoolean(train.Direction),
				},
				StopTimeUpdate: []*gtfs.TripUpdate_StopTimeUpdate{
					{
						StopId: ptr(staticData.stationToStopId[apiStationId]),
						Arrival: &gtfs.TripUpdate_StopTimeEvent{
							Time: timestamppbToInt64(train.ProjectedArrival),
						},
					},
				},
				Timestamp: timestamppbToUint64(train.LastUpdated),
			}
			b, err := json.Marshal(update)
			if err != nil {
				panic(err)
			}
			update.Trip.TripId = ptr(fmt.Sprintf("%x", md5.Sum(b)))
			entities = append(entities, &gtfs.FeedEntity{
				Id:         update.Trip.TripId,
				TripUpdate: update,
			})
		}
	}
	return &gtfs.FeedMessage{
		Header: &gtfs.FeedHeader{
			GtfsRealtimeVersion: ptr("0.2"),
			Incrementality:      gtfs.FeedHeader_FULL_DATASET.Enum(),
			Timestamp:           ptr(uint64(clock.Now().Unix())),
		},
		Entity: entities,
	}
}

func ptr[T any](t T) *T {
	return &t
}
