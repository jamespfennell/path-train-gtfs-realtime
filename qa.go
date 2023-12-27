package pathgtfsrt

import (
	"strings"

	sourceapi "github.com/jamespfennell/path-train-gtfs-realtime/proto/sourceapi"
)

const (
	ViaHobokenSuffix = "via hoboken"
)

func applyRouteQaToTrain(train *sourceapi.GetUpcomingTrainsResponse_UpcomingTrain) {
	normalizedLineName := strings.ToLower(train.LineName)
	if train.Route == sourceapi.Route_JSQ_33_HOB && !strings.HasSuffix(normalizedLineName, ViaHobokenSuffix) {
		train.Route = sourceapi.Route_JSQ_33
	}
}
