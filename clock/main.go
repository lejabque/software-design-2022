package main

import (
	"github.com/lejabque/software-design-2022/clock/internal/clock"
	"github.com/lejabque/software-design-2022/clock/internal/eventstats"
)

func main() {
	stats := eventstats.NewEventStats(clock.RealClock{})
	stats.IncEvent("event1")
	stats.IncEvent("event2")
	stats.IncEvent("event1")
	stats.PrintStatistic()
}
