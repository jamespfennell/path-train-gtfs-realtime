package feed

import s "github.com/jamespfennell/path-train-gtfs-realtime/feed/sourceapi"
import "testing"

var flagtests = []struct {
	in  string
	out s.Direction
}{
	{"TO_NY", s.Direction_TO_NY},
	{"TO_NJ", s.Direction_TO_NJ},
	{"DIRECTION_UNSPECIFIED", s.Direction_DIRECTION_UNSPECIFIED},
	{"random", s.Direction_DIRECTION_UNSPECIFIED},
}

func TestConvertDirectionAsStringToDirection(t *testing.T) {
	client := httpApiClient{}

	for _, testCase := range flagtests {
		actual := client.convertDirectionAsStringToDirection(testCase.in)
		if actual != testCase.out {
			t.Errorf("Input %s, expected %s, actual %s", testCase.in, testCase.out, actual)
		}
	}
}
