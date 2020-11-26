package monitoring

import "time"

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

func (m *Monitor) RecordUpdate(
	stopIDToErr map[string]error, builderErr error) {

}
