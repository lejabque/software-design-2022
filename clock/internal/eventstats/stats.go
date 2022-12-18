package eventstats

import (
	"fmt"
	"time"

	"github.com/lejabque/software-design-2022/clock/internal/clock"
)

type event struct {
	name string
	time time.Time
}

type EventStats struct {
	counters map[string]int64
	events   []event
	timer    clock.Clock
}

func NewEventStats(timer clock.Clock) *EventStats {
	return &EventStats{
		counters: make(map[string]int64),
		timer:    timer,
	}
}

func (s *EventStats) IncEvent(name string) {
	now := s.timer.Now()
	s.counters[name]++
	s.events = append(s.events, event{name: name, time: now})
	s.cleanup(now)
}

func (s *EventStats) GetEventStatsByName(name string) float64 {
	s.cleanup(s.timer.Now())
	return float64(s.counters[name]) / 60
}

func (s *EventStats) GetAllEventStats() map[string]float64 {
	s.cleanup(s.timer.Now())
	stats := make(map[string]float64)
	for name, count := range s.counters {
		stats[name] = float64(count) / 60
	}
	return stats
}

func (s *EventStats) PrintStatistic() {
	stats := s.GetAllEventStats()
	for name, count := range stats {
		fmt.Printf("%s: %f requests per minute\n", name, count)
	}
}

func (s *EventStats) cleanup(now time.Time) {
	removed := 0
	for removed < len(s.events) && now.Sub(s.events[removed].time) > time.Hour {
		name := s.events[removed].name
		s.counters[name]--
		if s.counters[name] == 0 {
			delete(s.counters, name)
		}
		removed++
	}
	// we assume that slicing works in O(1) time
	s.events = s.events[removed:]
}
