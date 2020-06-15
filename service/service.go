package service

import (
	"fmt"
	"sort"

	"github.com/netsak/time-tracker/timer"
)

// Service for managing the timers
type Service interface {
	AddTimer(names ...string) error
	ListTimer() []string
}

type timerservice struct {
	list map[string]timer.Timer
}

// New creates a new service
func New() (Service, error) {
	svc := timerservice{
		list: map[string]timer.Timer{},
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
		t := timer.Timer{}
		svc.list[name] = t
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
