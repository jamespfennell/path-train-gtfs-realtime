package monitoring

import (
	"encoding/json"
	"fmt"
	sourceapi "github.com/jamespfennell/path-train-gtfs-realtime/feed/sourceapi"
	"html/template"
	"io"
	"sync"
	"time"
)

// TODO: create a fixed length cache
// Create an interface first and then implement

type cache struct {
	store []Updates

	// Index of the oldest element of the cache
	oldest int
	len    int
	mutex  sync.RWMutex
}

func newCache(size int) cache {
	return cache{
		store: make([]Updates, size),
	}
}

func (c *cache) newest() int {
	return (c.oldest + c.len - 1) % len(c.store)
}

func (c *cache) GetMostRecent() (Updates, bool) {
	c.mutex.RLock()
	defer func() {
		c.mutex.RUnlock()
	}()
	if c.len == 0 {
		return Updates{}, false
	}
	return c.store[c.newest()], true
}

func (c *cache) UpdateMostRecent(u Updates) {
	c.mutex.Lock()
	defer func() {
		c.mutex.Unlock()
	}()
	c.store[c.newest()] = u
}

func (c *cache) Add(u Updates) {
	c.mutex.Lock()
	defer func() {
		c.mutex.Unlock()
	}()
	if c.len < len(c.store) {
		c.store[c.len] = u
		c.len++
		return
	}
	c.store[c.oldest] = u
	c.oldest = (c.oldest + 1) % len(c.store)
}

func (c *cache) Copy() []Updates {
	result := make([]Updates, c.len) // len is constant
	c.mutex.RLock()
	defer func() {
		c.mutex.RUnlock()
	}()
	sourceIndex := c.newest()
	for targetIndex := 0; targetIndex < c.len; targetIndex++ {
		result[targetIndex] = c.store[sourceIndex]
		if sourceIndex == 0 {
			sourceIndex = c.len
		}
		sourceIndex--
	}
	return result
}

type Monitor struct {
	c   cache
	bus chan Updates
}

func NewMonitor(cacheSize int) *Monitor {
	m := Monitor{
		bus: make(chan Updates),
		c:   newCache(cacheSize),
	}
	go m.background()
	return &m
}

func (m *Monitor) background() {
	for u := range m.bus {
		// fmt.Println("~~~~~~Recieved ", u)
		// TODO: all of this needs to happen on a separate goroutine
		//  This function should just send a message down a channel
		// Should append or update existing?
		// If existing does not exist
		// If existing is errored
		// If current is errored
		// If existing is more than an hour old
		m.c.Add(u)
		// fmt.Println("Added")
	}
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
	StopIDToErr      map[sourceapi.Station]error
	StopIDToNumTrips map[sourceapi.Station]int
	BuilderErr       error
}

func (u *Updates) TimeDescription() string {
	return u.FirstTime.Format("2006-02-01 15:04:05")
}

func (u *Updates) StationDescription(station sourceapi.Station) string {
	if u.StopIDToErr[station] != nil {
		return "F"
	}
	return fmt.Sprintf("%d", u.StopIDToNumTrips[station])
}

func (u *Updates) hasErr() bool {
	if u.BuilderErr != nil {
		return true
	}
	for _, err := range u.StopIDToErr {
		if err != nil {
			return true
		}
	}
	return false
}

type StationUpdateResult struct {
	NumTrips int
	Err      error
}

func (m *Monitor) RecordUpdate(
	stopIDToErr map[sourceapi.Station]StationUpdateResult, builderErr error) {
	u := Updates{
		BuilderErr:       builderErr,
		StopIDToNumTrips: map[sourceapi.Station]int{},
		StopIDToErr:      map[sourceapi.Station]error{},
		FirstTime:        time.Now(),
		LastTime:         time.Now(),
	}
	for stopID, result := range stopIDToErr {
		u.StopIDToNumTrips[stopID] = result.NumTrips
		u.StopIDToErr[stopID] = result.Err
	}
	m.bus <- u
}

func (m *Monitor) WriteJson(w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(m.c.Copy())
}

func (m *Monitor) WriteHTML(w io.Writer) error {
	tmpl, err := template.New("index.html").Parse(statusHtmlTemplate)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, m.c.Copy())
}
