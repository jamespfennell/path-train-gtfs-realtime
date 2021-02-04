package monitoring

import (
	"github.com/jamespfennell/path-train-gtfs-realtime/feed"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var numUpdatesCounter *prometheus.CounterVec
var lastUpdateGauge *prometheus.GaugeVec
var successfulUpdateLatencyGauge prometheus.Gauge
var numTripStopTimesGauge *prometheus.GaugeVec
var numRequestsCounter *prometheus.CounterVec

var lastSuccessfulUpdateTime *time.Time

func init() {
	numUpdatesCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "path_train_gtfsrt_num_updates",
			Help: "The total number of updates that have occurred",
		},
		[]string{"has_error"},
	)
	lastUpdateGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "path_train_gtfsrt_last_update",
			Help: "The time of the last update",
		},
		[]string{"has_error"},
	)
	successfulUpdateLatencyGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "path_train_gtfsrt_successful_update_latency",
			Help: "The time since the last successful update",
		})
	numTripStopTimesGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "path_train_gtfsrt_num_trip_stop_times",
			Help: "The number of trip stop times per station and direction",
		},
		[]string{"station", "direction"},
	)
	numRequestsCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "path_train_gtfsrt_num_requests",
			Help: "The number of times the GTFS-RT feed is requested",
		},
		[]string{"code"},
	)
}

func Listen(c <-chan feed.UpdateResult) {
	for result := range c {
		hasErr := result.GtfsBuilderErr != nil
		for station, stationResult := range result.StationsToResult {
			numTripStopTimesGauge.WithLabelValues(station.String(), "NY").Set(float64(stationResult.NumTripsNY))
			numTripStopTimesGauge.WithLabelValues(station.String(), "NJ").Set(float64(stationResult.NumTripsNJ))
			hasErr = hasErr || stationResult.Err != nil
		}

		if lastSuccessfulUpdateTime != nil {
			d := time.Now().Sub(*lastSuccessfulUpdateTime)
			successfulUpdateLatencyGauge.Set(d.Seconds())
		}
		if !hasErr {
			if lastSuccessfulUpdateTime == nil {
				lastSuccessfulUpdateTime = &time.Time{}
			}
			*lastSuccessfulUpdateTime = time.Now()
		}

		if hasErr {
			numUpdatesCounter.WithLabelValues("true").Inc()
			lastUpdateGauge.WithLabelValues("true").SetToCurrentTime()
		} else {
			numUpdatesCounter.WithLabelValues("false").Inc()
			lastUpdateGauge.WithLabelValues("false").SetToCurrentTime()
		}
	}
}

func CountRequests(handler http.Handler) http.Handler {
	return promhttp.InstrumentHandlerCounter(numRequestsCounter, handler)
}

func HttpHandler() http.Handler {
	return promhttp.Handler()
}
