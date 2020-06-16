package tray

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/getlantern/systray"
	"github.com/netsak/time-tracker/service"
)

// Tray functions
type Tray interface {
	Run()
}

type systemtray struct {
	svc   service.Service
	items map[string]*systray.MenuItem
}

// New creates a new system tray for time tracking
func New(timerService service.Service) (Tray, error) {
	t := systemtray{
		svc:   timerService,
		items: make(map[string]*systray.MenuItem),
	}
	go func() {
		eventChannel := t.svc.OnEvent()
		for {
			select {
			case evt := <-eventChannel:
				log.Printf("event: %+v", evt)
				t.updateTime(evt.TimerName, evt.TimerDuration)
				t.updateMenu()
				// systray.SetTitle(evt.Duration())
			}
		}
	}()
	return &t, nil
}

func (t *systemtray) Run() {
	systray.Run(t.onReady, t.onExit)
}

// onReady setup the system tray and creates the menu linked with actions
func (t *systemtray) onReady() {
	log.Println("system tray is ready, setting up menus and actions...")
	// top setting visible only on startup
	systray.SetIcon(getIcon("assets/time-off.png"))
	systray.SetTooltip("Time Tracker is off")
	systray.SetTitle("Time Tracker")
	// total time
	menuTotalTime := systray.AddMenuItem("0:00:00", "")
	menuTotalTime.Disable()
	systray.AddSeparator()
	// create dynamic tasks menu
	for _, name := range t.svc.ListTimer() {
		item := systray.AddMenuItem(fmt.Sprintf("%s 00:00:00", name), "name")
		log.Printf("add item %s: %+v", name, item)
		t.items[name] = item
		go t.onClick(name, item)
	}
	// add pause and quit
	systray.AddSeparator()
	menuStop := systray.AddMenuItem("Stop", "Stop the time tracking")
	go t.onStop(menuStop)
	menuQuit := systray.AddMenuItem("Quit", "Enough work done today!")
	go t.onQuit(menuQuit)
}

func (t *systemtray) onClick(name string, item *systray.MenuItem) {
	log.Printf("set onClick handler for %s", name)
	for {
		select {
		case <-item.ClickedCh:
			log.Printf("timer %s clicked", name)
			t.setTrayTimeOn(name)
			t.svc.ActivateTimer(name)
		}
	}
}

func (t *systemtray) onStop(item *systray.MenuItem) {
	log.Printf("set onStop handler for %s", item)
	for {
		select {
		case <-item.ClickedCh:
			log.Printf("item %s clicked", item)
			t.svc.StopCurrentTimer()
			t.setTrayTimeOff()
		}
	}
}

func (t *systemtray) onQuit(item *systray.MenuItem) {
	log.Printf("set onQuit handler for %s", item)
	for {
		select {
		case <-item.ClickedCh:
			systray.Quit()
		}
	}
}

func (t *systemtray) onExit() {
	log.Println("exiting time tracker...")
}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func (t *systemtray) setTrayTimeOff() {
	systray.SetIcon(getIcon("assets/time-off.png"))
	systray.SetTitle("Time Tracker")
	systray.SetTooltip("Time Tracker is off")
}

func (t *systemtray) setTrayTimeOn(name string) {
	systray.SetIcon(getIcon("assets/time-on.png"))
	systray.SetTitle("00:00:00")
	systray.SetTooltip(fmt.Sprintf("Tracking time for %s", name))
}

func (t *systemtray) updateTime(name string, duration time.Duration) {
	systray.SetIcon(getIcon("assets/time-on.png"))
	durationStr := formatDuration(duration)
	systray.SetTitle(durationStr)
	systray.SetTooltip(fmt.Sprintf("Tracking time for %s", name))
}

func (t *systemtray) updateMenu() {
	for name, menu := range t.items {
		timer := t.svc.GetTimer(name)
		durationStr := formatDuration(timer.TotalDuration)
		titleStr := fmt.Sprintf("%s\t%s", durationStr, name)
		menu.SetTitle(titleStr)
	}
}

func formatDuration(duration time.Duration) string {
	total := int(duration.Seconds())
	hours := int(total / (60 * 60) % 24)
	minutes := int(total/60) % 60
	seconds := int(total % 60)
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
