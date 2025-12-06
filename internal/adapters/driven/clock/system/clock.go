package system

import "time"

type SystemClock struct {}

func (s SystemClock) Current() time.Time {
	return time.Now()
}
