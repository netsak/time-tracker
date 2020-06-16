package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/netsak/time-tracker/timer"
)

// Service for managing the timers
type Service interface {
	AddTimer(names ...string) error
	ListTimer() []string
	ActivateTimer(name string) error
	OnEvent() chan Event
	StopCurrentTimer()
	GetTimer(name string) *timer.Timer
}

// Event posted when a timer changes
type Event struct {
	TimerName     string
	TimerDuration time.Duration
}

// Duration prints the timer event duration as HH:MM:SS
// func (evt Event) Duration() string {
// 	total := int(evt.TimerDuration.Seconds())
// 	hours := int(total / (60 * 60) % 24)
// 	minutes := int(total/60) % 60
// 	seconds := int(total % 60)
// 	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
// }

type timerservice struct {
	list         map[string]*timer.Timer
	current      *timer.Timer
	eventChannel chan Event
}

// New creates a new service
func New() (Service, error) {
	svc := timerservice{
		list:         make(map[string]*timer.Timer),
		eventChannel: make(chan Event, 1),
	}
	return &svc, nil
}

// AddTimer adds a new timer to the service
func (svc *timerservice) AddTimer(names ...string) error {
	for _, name := range names {
		_, found := svc.list[name]
		if found {
			return fmt.Errorf("timer %s already defined", name)
		}
		t := timer.Timer{
			Name: name,
		}
		svc.list[name] = &t
	}
	return nil
}

// ListTimer returns a list of all timer names
func (svc *timerservice) ListTimer() []string {
	ret := []string{}
	for key := range svc.list {
		ret = append(ret, key)
	}
	sort.Strings(ret)
	return ret
}

// ActivateTimer activates a timer by name
func (svc *timerservice) ActivateTimer(name string) error {
	// TODO: rename to StartTimer
	if svc.current != nil {
		svc.current.Stop()
	}
	t, found := svc.list[name]
	if !found {
		return fmt.Errorf("timer %s not found", name)
	}
	ticker := t.Start()
	svc.current = t
	go func() {
		for {
			select {
			case t := <-ticker:
				evt := Event{
					TimerName:     name,
					TimerDuration: t,
				}
				svc.eventChannel <- evt
			}
		}
	}()
	return nil
}

// StopCurrentTimer stops the current timer if one is running
func (svc *timerservice) StopCurrentTimer() {
	if svc.current != nil {
		svc.current.Stop()
		svc.current = nil
	}
	return
}

// OnEvent returns a channel with change events
func (svc *timerservice) OnEvent() chan Event {
	return svc.eventChannel
}

// GetTimer gets a timer by name
func (svc *timerservice) GetTimer(name string) *timer.Timer {
	return svc.list[name]
}
