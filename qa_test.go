package pathgtfsrt

import (
	"github.com/google/go-cmp/cmp"
	sourceapi "github.com/jamespfennell/path-train-gtfs-realtime/proto/sourceapi"
	"google.golang.org/protobuf/testing/protocmp"
	"testing"
)

func TestRouteQa(t *testing.T) {
	for _, tc := range []struct {
		name  string
		train *sourceapi.GetUpcomingTrainsResponse_UpcomingTrain
		want  *sourceapi.GetUpcomingTrainsResponse_UpcomingTrain
	}{
		{
			name: "no change",
			train: &sourceapi.GetUpcomingTrainsResponse_UpcomingTrain{
				Route:    sourceapi.Route_JSQ_33_HOB,
				LineName: "33rd Street via Hoboken",
			},
			want: &sourceapi.GetUpcomingTrainsResponse_UpcomingTrain{
				Route:    sourceapi.Route_JSQ_33_HOB,
				LineName: "33rd Street via Hoboken",
			},
		},
		{
			name: "no change, lower case",
			train: &sourceapi.GetUpcomingTrainsResponse_UpcomingTrain{
				Route:    sourceapi.Route_JSQ_33_HOB,
				LineName: "33rd street via hoboken",
			},
			want: &sourceapi.GetUpcomingTrainsResponse_UpcomingTrain{
				Route:    sourceapi.Route_JSQ_33_HOB,
				LineName: "33rd street via hoboken",
			},
		},
		{
			name: "change",
			train: &sourceapi.GetUpcomingTrainsResponse_UpcomingTrain{
				Route:    sourceapi.Route_JSQ_33_HOB,
				LineName: "33rd Street",
			},
			want: &sourceapi.GetUpcomingTrainsResponse_UpcomingTrain{
				Route:    sourceapi.Route_JSQ_33,
				LineName: "33rd Street",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			applyRouteQaToTrain(tc.train)
			if diff := cmp.Diff(tc.want, tc.train, protocmp.Transform()); diff != "" {
				t.Errorf("applyRouteQaToTrain() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
