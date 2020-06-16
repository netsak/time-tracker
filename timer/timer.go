package timer

import (
	"fmt"
	"log"
	"time"
)

// Timer for the logging
type Timer struct {
	Name            string
	StartTime       time.Time
	EndTime         time.Time
	TotalDuration   time.Duration
	CurrentDuration time.Duration
	ticker          *time.Ticker
	ticking         bool
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
	t.ticking = true
	go func() {
		for now := range t.ticker.C {
			t.CurrentDuration = now.Sub(t.StartTime)
			ret <- t.CurrentDuration
		}
	}()
	return ret
}

// Stop the timer
func (t *Timer) Stop() {
	t.ticker.Stop()
	t.ticking = false
	t.EndTime = time.Now()
	diff := t.EndTime.Sub(t.StartTime)
	t.TotalDuration += diff
	log.Printf("timer stopped: %s diff=%s", t, diff)
}

// IsActive checks if the timer is currently ticking
func (t Timer) IsActive() bool {
	return t.ticking
}
