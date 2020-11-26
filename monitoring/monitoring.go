package monitoring

import (
	"fmt"
	sourceapi "github.com/jamespfennell/path-train-gtfs-realtime/feed/sourceapi"
	"time"
)

type Monitor struct{}

func NewMonitor() *Monitor {
	return &Monitor{}
}

// Updates contains information on one or more updates
type Updates struct {
	Count          int
	FirstTime      time.Time
	LastTime       time.Time
	SuccessLatency struct {
		Mean float64
		Max  float64
	}
	StopIDToErr map[string]error
	BuilderErr  error
}

type StationUpdateResult struct {
	NumTrips int
	Err      error
}

func (m *Monitor) RecordUpdate(
	stopIDToErr map[sourceapi.Station]StationUpdateResult, builderErr error) {
	fmt.Println("Received monitoring update", stopIDToErr)
}
