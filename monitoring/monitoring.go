package monitoring

import (
	"encoding/json"
	"errors"
	"fmt"
	sourceapi "github.com/jamespfennell/path-train-gtfs-realtime/feed/sourceapi"
	"html/template"
	"io"
	"sync"
	"time"
)

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
	c                        cache
	bus                      chan Updates
	updatePeriod             time.Duration
	lastSuccessfulUpdateTime *time.Time
}

func NewMonitor(cacheSize int, updatePeriod time.Duration) *Monitor {
	m := Monitor{
		bus:          make(chan Updates),
		c:            newCache(cacheSize),
		updatePeriod: updatePeriod,
	}
	go m.background()
	return &m
}

func (m *Monitor) background() {
	for u := range m.bus {
		if m.lastSuccessfulUpdateTime == nil {
			u.SuccessLatency.Max = -1
			u.SuccessLatency.Mean = -1
		} else {
			u.SuccessLatency.Max = u.LastTime.Sub(*m.lastSuccessfulUpdateTime).Round(100 * time.Millisecond)
			u.SuccessLatency.Mean = u.LastTime.Sub(*m.lastSuccessfulUpdateTime).Round(100 * time.Millisecond)
		}
		if u.SuccessLatency.Mean >= 0 && u.SuccessLatency.Mean >= m.updatePeriod*3 {
			u.SuccessLatency.Err = errors.New(
				fmt.Sprintf(
					"time elapsed is more than 3 times the target update period (%s)",
					m.updatePeriod))
		}
		if !u.hasErr() {
			lastTime := u.LastTime
			m.lastSuccessfulUpdateTime = &lastTime
		}
		lastUpdate, exists := m.c.GetMostRecent()
		if !exists {
			m.c.Add(u)
			continue
		}
		if !shouldMergeUpdates(lastUpdate, u) {
			m.c.Add(u)
			continue
		}
		m.c.UpdateMostRecent(merge(lastUpdate, u))
	}
}

func shouldMergeUpdates(last, new Updates) bool {
	if last.hasErr() || new.hasErr() {
		return false
	}
	if last.SuccessLatency.Err != nil || new.SuccessLatency.Err != nil {
		return false
	}
	if new.FirstTime.Sub(last.FirstTime) >= time.Hour {
		return false
	}
	for stationID, _ := range new.StopIDToNumTrips {
		if last.StopIDToNumTrips[stationID] != new.StopIDToNumTrips[stationID] {
			return false
		}
	}
	return true
}

func merge(last, new Updates) Updates {
	last.LastTime = new.LastTime
	if last.SuccessLatency.Max < new.SuccessLatency.Max {
		last.SuccessLatency.Max = new.SuccessLatency.Max
	}
	if last.SuccessLatency.Mean < 0 {
		last.SuccessLatency.Mean = new.SuccessLatency.Mean
	} else {
		total := int64(last.SuccessLatency.Mean)*last.Count + int64(new.SuccessLatency.Mean)
		last.SuccessLatency.Mean = time.Duration(total / (last.Count + 1)).Round(100 * time.Millisecond)
	}
	last.Count = last.Count + 1
	return last
}

// Updates contains information on one or more updates
type Updates struct {
	Count          int64
	FirstTime      time.Time
	LastTime       time.Time
	SuccessLatency struct {
		Mean time.Duration
		Max  time.Duration
		Err  error
	}
	StopIDToErr      map[sourceapi.Station]error
	StopIDToNumTrips map[sourceapi.Station]int
	BuilderErr       error
}

func (u *Updates) TimeDescription() template.HTML {
	if u.Count == 1 {
		return template.HTML(u.LastTime.Format("2006-02-01 15:04:05"))
	}
	return template.HTML(fmt.Sprintf("%s<br />|<br />(%d updates)<br />|<br />%s",
		u.LastTime.Format("2006-02-01 15:04:05"),
		u.Count,
		u.FirstTime.Format("2006-02-01 15:04:05"),
	))
}

func (u *Updates) StationDescription(station sourceapi.Station) template.HTML {
	if u.StopIDToErr[station] != nil {
		return template.HTML(
			fmt.Sprintf(`<span title="%s" class="hover">%s</span>`, u.StopIDToErr[station], "F"))
	}
	return template.HTML(fmt.Sprintf("%d", u.StopIDToNumTrips[station]))
}

func (u *Updates) StationClass(station sourceapi.Station) string {
	if u.StopIDToErr[station] != nil {
		return "failure"
	}
	return "success"
}

func (u *Updates) LatencyDescription() template.HTML {
	if u.SuccessLatency.Mean < 0 {
		if u.hasErr() {
			return `<span title="no successful update yet" class="hover">N/A</span>`
		}
		return `<span title="this is the first successful update" class="hover">N/A</span>`
	}
	if u.SuccessLatency.Err != nil {
		return template.HTML(
			fmt.Sprintf(`<span title="%s" class="hover">%s</span>`, u.SuccessLatency.Err, u.SuccessLatency.Mean))
	}
	if u.Count == 1 {
		return template.HTML(fmt.Sprintf("%s", u.SuccessLatency.Mean))
	}
	return template.HTML(fmt.Sprintf("mean<br />%s<br /><br />max<br>%s",
		u.SuccessLatency.Mean, u.SuccessLatency.Max))
}
func (u *Updates) LatencyClass() template.HTML {
	if u.SuccessLatency.Mean < 0 {
		return ""
	}
	if u.SuccessLatency.Err != nil {
		return "failure"
	}
	return "success"
}

func (u *Updates) BuilderDescription() template.HTML {
	if u.BuilderErr != nil {
		return template.HTML(fmt.Sprintf(`<span title="%s" class="hover">F</span>`, u.BuilderErr))
	}
	return "S"
}

func (u *Updates) BuilderClass() template.HTML {
	if u.BuilderErr != nil {
		return "failure"
	}
	return "success"
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
	t := time.Now()
	u := Updates{
		BuilderErr:       builderErr,
		StopIDToNumTrips: map[sourceapi.Station]int{},
		StopIDToErr:      map[sourceapi.Station]error{},
		FirstTime:        t,
		LastTime:         t,
		Count:            1,
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
	data := struct {
		StationIDs   []sourceapi.Station
		StationNames []string
		Updates      []Updates
	}{
		StationIDs: []sourceapi.Station{10, 11, 12, 13, 9, 5, 4, 2, 8, 3, 1, 7, 6},
		StationNames: []string{
			"9th St",
			"14th St",
			"23rd St",
			"33rd St",
			"Christopher St",
			"Exchange Pl",
			"Grove St",
			"Harrison",
			"Hoboken",
			"Journal Sq",
			"Newark",
			"Newport",
			"WTC",
		},
		Updates: m.c.Copy(),
	}
	return tmpl.Execute(w, data)
}
