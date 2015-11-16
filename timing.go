package timing

import (
	"time"
	"github.com/Sirupsen/logrus"
)

//running timely with offset and timezone, stop by stop channel
func TimingAtClock(do func(), stop chan bool, interval, localeClock, timeZoneOffset time.Duration) {
	now := time.Now()
	offset := (localeClock-timeZoneOffset%(time.Hour*24))%interval
	next := now.Truncate(interval).Add(offset)
	if next.Before(now) {
		next = next.Add(interval)
	}
	logrus.Debug("Timing At Clock next:",next)
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
		logrus.Debug("Timing At Clock next:",next)
		timer.Reset(next.Sub(now))
	}
}

//running timely with offset, stop by stop channel
func Timing(do func(), stop chan bool, interval time.Duration, utcOffset time.Duration) {
	TimingAtClock(do, stop, interval, utcOffset, 0)
}

//New 新建Timer
func New(interval,localeClock,timeZoneOffset time.Duration) *Timer {
	output := make(chan bool, 10)
	stop := make(chan bool, 0)
	go TimingAtClock(func(){
		output<-true
	},stop,interval,localeClock,timeZoneOffset)
	return &Timer{output:output,stop:stop}
}

//Timer 计时器
type Timer struct{
	output chan bool
	stop chan bool
}

//Stop 停止计时器
func (t Timer) Stop() {
	close(t.stop)
	close(t.output)
}

//Signal 给予定时信号
func (t Timer) Signal() chan bool {
	return t.output
}

