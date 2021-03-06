package timing_test

import (
	"fmt"
	"math"
	"testing"
	"time"
"github.com/bxy09/timing"
	"github.com/Sirupsen/logrus"
)

const BJ_TZ time.Duration = time.Hour*8

func TestTiming(t *testing.T) {
	t.Log("In Timing")
	stop := make(chan bool, 1)
	start:= time.Now()
	hour, min, second := start.Clock()
	const offset int = 7
	fmt.Println(start, hour, min, second)
	go timing.TimingAtClock(func() {
		t.Log(time.Now())
		current := time.Now()
		if math.Abs(current.Sub(start).Seconds() - float64(offset)) > 2 {
			fmt.Println(current.Sub(start).Seconds() - float64(offset))
			t.Fatal("time is not right!", current, start);
		}
		close(stop)
	}, stop, time.Hour*24,
		time.Hour*time.Duration(hour)+time.Minute*time.Duration(min)+time.Second*time.Duration(second+offset),
		BJ_TZ)
	timer := time.NewTimer(10*time.Second)
	select {
	case <-stop:
		timer.Stop()
	case <-timer.C:
		t.Fatal("Don't run the function on time");
	}
}

func TestSecond(t *testing.T) {
	t.Log("In TestSecond")
	times := 0
	var last time.Time

	stop := make(chan bool, 1)
	go timing.Timing(func() {
		t.Log(time.Now())
		if times == 0 {
			last = time.Now()
		} else {
			current := time.Now()
			if math.Abs(current.Sub(last.Add(-time.Microsecond)).Seconds()-1) > 0.01 {
				t.Fatal("interval is not a second")
			}
			last = current
		}
		times++
		if times == 10 {
			close(stop)
		}
	}, stop, time.Second, 0)
	<-stop
}

func TestSecondInClass(t *testing.T) {
	t.Log("In TestSecondInClass")
	logrus.SetLevel(logrus.DebugLevel)
	timer := timing.New(time.Second, 0, time.Hour*8)
	limit := 8
	for <-timer.Signal() {
		limit--
		if limit == 0 {
			timer.Stop()
		}
	}
}