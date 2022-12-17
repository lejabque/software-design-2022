package clock

import "time"

type FakeClock struct {
	currentTime time.Time
}

func NewFakeClock(t time.Time) *FakeClock {
	return &FakeClock{currentTime: t}
}

func (c *FakeClock) Now() time.Time {
	return c.currentTime
}

func (c *FakeClock) Add(d time.Duration) {
	c.currentTime = c.currentTime.Add(d)
}

func (c *FakeClock) Set(t time.Time) {
	c.currentTime = t
}
