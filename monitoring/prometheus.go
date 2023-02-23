package monitoring

import (
	"net/http"

	gtfs "github.com/jamespfennell/path-train-gtfs-realtime/proto/gtfsrt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var numUpdatesCounter prometheus.Counter
var numRequestErrs prometheus.Counter
var lastUpdateGauge prometheus.Gauge
var numTripStopTimesGauge *prometheus.GaugeVec
var numRequestsCounter *prometheus.CounterVec

func init() {
	numUpdatesCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "path_train_gtfsrt_num_updates",
			Help: "Number of completed updates",
		},
	)
	numRequestErrs = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "path_train_gtfsrt_num_source_api_errors",
			Help: "Number of errors when retrieving realtime data from the source API",
		},
	)
	lastUpdateGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "path_train_gtfsrt_last_update",
			Help: "Time of the last completed update",
		},
	)
	numTripStopTimesGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "path_train_gtfsrt_num_trip_stop_times",
			Help: "Number of trip stop times per station and direction",
		},
		[]string{"stop_id", "direction"},
	)
	numRequestsCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "path_train_gtfsrt_num_requests",
			Help: "Number of times the GTFS-RT feed has been requested",
		},
		[]string{"code"},
	)
}

func RecordUpdate(msg *gtfs.FeedMessage, errs []error) {
	numTripStopTimesGauge.Reset()
	for _, entity := range msg.GetEntity() {
		directionID := "NY"
		if entity.GetTripUpdate().GetTrip().GetDirectionId() == 0 {
			directionID = "NJ"
		}
		for _, stopTimeUpdate := range entity.GetTripUpdate().GetStopTimeUpdate() {
			stopID := stopTimeUpdate.GetStopId()
			numTripStopTimesGauge.WithLabelValues(stopID, directionID).Inc()
		}
	}
	numUpdatesCounter.Inc()
	numRequestErrs.Add(float64(len(errs)))
	lastUpdateGauge.SetToCurrentTime()
}

func CountRequests(handler http.Handler) http.Handler {
	return promhttp.InstrumentHandlerCounter(numRequestsCounter, handler)
}

func HttpHandler() http.Handler {
	return promhttp.Handler()
}
