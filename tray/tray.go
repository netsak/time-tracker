package tray

import (
	"fmt"
	"io/ioutil"
	"log"

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
	return &t, nil
}

func (t *systemtray) Run() {
	systray.Run(t.onReady, t.onExit)
}

func (t *systemtray) onReady() {
	log.Println("system tray is ready, setting up menus and actions...")

	systray.SetIcon(getIcon("assets/time-off.png"))
	systray.SetTooltip("Time Tracker")
	systray.SetTitle("Time Tracker")

	menuTotalTime := systray.AddMenuItem("0:00:00", "")
	menuTotalTime.Disable()
	systray.AddSeparator()

	// create dynamic tasks menu
	for _, name := range t.svc.ListTimer() {
		item := systray.AddMenuItem(fmt.Sprintf("%s 00:00:00", name), "name")
		log.Printf("add item %s: %+v", name, item)
		t.items[name] = item
		go t.onClick(item)
	}

	// add quit
	systray.AddSeparator()
	menuQuit := systray.AddMenuItem("Quit", "Enough work done today!")
	go t.onQuit(menuQuit)
}

func (t *systemtray) onClick(item *systray.MenuItem) {
	log.Printf("set onClick handler for %s", item)
	for {
		select {
		case <-item.ClickedCh:
			log.Printf("item %s clicked", item)
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
