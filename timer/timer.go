package timer

import (
	"fmt"
	"log"
	"time"
)

// Timer for the logging
type Timer struct {
	Name          string
	StartTime     time.Time
	EndTime       time.Time
	TotalDuration time.Duration
	ticker        *time.Ticker
	IsCurrent     bool
}

// String returns a string representation of the current state
func (t Timer) String() string {
	return fmt.Sprintf("%s[%s %s-%s", t.Name, t.TotalDuration, t.StartTime, t.EndTime)
}

// Start the timer and posts the duration to a channel
func (t *Timer) Start() chan time.Duration {
	ret := make(chan time.Duration, 1)
	t.StartTime = time.Now()
	t.EndTime = time.Time{}
	log.Printf("timer started: %s", t)
	t.ticker = time.NewTicker(1 * time.Second)
	go func() {
		for now := range t.ticker.C {
			diff := now.Sub(t.StartTime)
			ret <- diff
		}
	}()
	return ret
}

// Stop the timer
func (t *Timer) Stop() {
	t.ticker.Stop()
	t.EndTime = time.Now()
	diff := t.EndTime.Sub(t.StartTime)
	t.TotalDuration += diff
	log.Printf("timer stopped: %s diff=%s", t, diff)
}
