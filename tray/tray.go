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
	svc           service.Service
	items         map[string]*systray.MenuItem
	menuTotalTime *systray.MenuItem
	menuStop      *systray.MenuItem
	menuQuit      *systray.MenuItem
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
	t.menuTotalTime = systray.AddMenuItem("0:00:00", "")
	t.menuTotalTime.Disable()
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
	t.menuStop = systray.AddMenuItem("Stop", "Stop the time tracking")
	go t.onStop(t.menuStop)
	t.menuQuit = systray.AddMenuItem("Quit", "Enough work done today!")
	go t.onQuit(t.menuQuit)
}

func (t *systemtray) onClick(name string, item *systray.MenuItem) {
	log.Printf("set onClick handler for %s", name)
	for {
		select {
		case <-item.ClickedCh:
			log.Printf("timer %s clicked", name)
			t.setTrayTimeOn(name)
			t.svc.StartTimer(name)
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
	var total time.Duration
	for name, menu := range t.items {
		timer := t.svc.GetTimer(name)
		durationStr := formatDuration(timer.TotalDuration)
		tracking := ""
		if timer.IsActive() {
			tracking = "...TRACKING..."
		}
		titleStr := fmt.Sprintf("%s\t%s\t%s", durationStr, name, tracking)
		menu.SetTitle(titleStr)
		total += timer.TotalDuration
	}
	durationStr := formatDuration(total)
	titleStr := fmt.Sprintf("%s\ttotal", durationStr)
	t.menuTotalTime.SetTitle(titleStr)
}

func formatDuration(duration time.Duration) string {
	total := int(duration.Seconds())
	hours := int(total / (60 * 60) % 24)
	minutes := int(total/60) % 60
	seconds := int(total % 60)
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
