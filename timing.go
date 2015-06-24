package timing

import (
	"time"
)

//running timely with offset and timezone, stop by stop channel
func TimingAtClock(do func(), stop chan bool, interval, localeClock, timeZoneOffset time.Duration) {
	now := time.Now()
	offset := localeClock-timeZoneOffset%(time.Hour*24)
	next := now.Truncate(interval).Add(offset)
	if next.Before(now) {
		next = next.Add(interval)
	}
	timer := time.NewTimer(next.Sub(now))
	stopped := false
	for {
		select {
		case <-timer.C:
		case <-stop:
			timer.Stop()
			stopped = true
		}
		if stopped {
			break
		}
		do()
		now = time.Now()
		next = now.Truncate(interval).Add(offset)
		if next.Before(now) {
			next = next.Add(interval)
		}
		timer.Reset(next.Sub(now))
	}
}

//running timely with offset, stop by stop channel
func Timing(do func(), stop chan bool, interval time.Duration, utcOffset time.Duration) {
	TimingAtClock(do, stop, interval, utcOffset, 0)
}


