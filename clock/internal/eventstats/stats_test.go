package eventstats

import (
	"testing"
	"time"

	"github.com/lejabque/software-design-2022/clock/internal/clock"
	"github.com/stretchr/testify/suite"
)

const eqDelta = 1e-9
const rpm = float64(1) / 60

type StatsTestSuite struct {
	suite.Suite
	timer clock.FakeClock
	stats *EventStats
}

func (s *StatsTestSuite) SetupTest() {
	s.stats = NewEventStats(&s.timer)
}

func (s *StatsTestSuite) TestEmpty() {
	s.InDelta(0, s.stats.GetEventStatsByName("ev1"), eqDelta)
	s.Equal(0, len(s.stats.GetAllEventStats()))
}

func (s *StatsTestSuite) TestIncAndGet() {
	s.timer.Set(time.Unix(123, 0))
	s.stats.IncEvent("ev1")
	s.stats.IncEvent("ev2")
	s.timer.Add(5 * time.Minute)
	s.stats.IncEvent("ev1")

	all := s.stats.GetAllEventStats()
	s.InDelta(2*rpm, s.stats.GetEventStatsByName("ev1"), eqDelta)
	s.InDelta(2*rpm, all["ev1"], eqDelta)
	s.InDelta(rpm, s.stats.GetEventStatsByName("ev2"), eqDelta)
	s.InDelta(rpm, all["ev2"], eqDelta)

	s.timer.Add(56 * time.Minute)
	s.stats.IncEvent("ev3")

	all = s.stats.GetAllEventStats()
	s.Equal(2, len(all))
	s.InDelta(rpm, s.stats.GetEventStatsByName("ev1"), eqDelta)
	s.InDelta(rpm, all["ev1"], eqDelta)
	s.InDelta(0, s.stats.GetEventStatsByName("ev2"), eqDelta)
	s.InDelta(0, all["ev2"], eqDelta)
	s.InDelta(rpm, s.stats.GetEventStatsByName("ev3"), eqDelta)
	s.InDelta(rpm, all["ev3"], eqDelta)

	s.timer.Add(5 * time.Minute)
	all = s.stats.GetAllEventStats()
	s.Equal(1, len(all))
	s.InDelta(rpm, s.stats.GetEventStatsByName("ev3"), eqDelta)
	s.InDelta(rpm, all["ev3"], eqDelta)

	s.timer.Add(60 * time.Minute)
	all = s.stats.GetAllEventStats()
	s.Equal(0, len(all))
}

func Test(t *testing.T) {
	suite.Run(t, new(StatsTestSuite))
}
