package clock

import "time"

type FakeClock struct {
	NowTime time.Time
}

func (c *FakeClock) Now() time.Time {
	return c.NowTime
}

func (c *FakeClock) Add(d time.Duration) {
	c.NowTime = c.NowTime.Add(d)
}

func (c *FakeClock) Set(t time.Time) {
	c.NowTime = t
}

func NewFakeClock(t time.Time) *FakeClock {
	return &FakeClock{NowTime: t}
}
