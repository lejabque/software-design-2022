package eventstats

import (
	"fmt"
	"time"

	"github.com/lejabque/software-design-2022/clock/internal/clock"
)

type Stats struct {
	// requests per minute for last hour
	counters map[string][]time.Time
	timer    clock.Clock
}

func NewStats(timer clock.Clock) *Stats {
	return &Stats{
		counters: make(map[string][]time.Time),
		timer:    timer,
	}
}

func (s *Stats) IncEvent(name string) {
	s.counters[name] = append(s.counters[name], s.timer.Now())
	// cleanup old events:
	now := s.timer.Now()
	for i, t := range s.counters[name] {
		if now.Sub(t) > time.Hour {
			s.counters[name] = s.counters[name][i:]
			break
		}
	}
}

func (s *Stats) GetEventStatsByName(name string) float64 {
	now := s.timer.Now()
	count := 0
	for _, t := range s.counters[name] {
		if now.Sub(t) < time.Hour {
			count++
		}
	}
	// return requests per minute
	return float64(count) / 60.0
}

func (s *Stats) GetAllEventStats() map[string]float64 {
	stats := make(map[string]float64)
	for name := range s.counters {
		stats[name] = s.GetEventStatsByName(name)
	}
	return stats
}

func (s *Stats) PrintStatistic() {
	// print statistic for event with name
	stats := s.GetAllEventStats()
	for name, count := range stats {
		fmt.Printf("%s: %f requests per minute\n", name, count)
	}
}
