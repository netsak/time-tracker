package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/getlantern/systray"
)

func main() {
	fmt.Println("hello")
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(getIcon("assets/time-off.png"))
	systray.SetTitle("I'm alive!")
	systray.SetTooltip("Look at me, I'm a tooltip!")

	systray.AddMenuItem("0:00:00", "....")
	systray.AddSeparator()
	meeting := systray.AddMenuItem("Meeting", "I'm in a meeting")
	meeting.SetIcon(getIcon("assets/time-off.png")) // only works on osx :-(
	systray.AddMenuItem("Task", "I'm working on a task")
	systray.AddMenuItem("Bug", "I'm debugging...")
	systray.AddMenuItem("Other", "I'm distracted :-(")
	systray.AddSeparator()
	quit := systray.AddMenuItem("Quit", "Enough work done today!")

	go func() {
		for {
			systray.SetTitle(getClockTime())
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			select {
			case <-quit.ClickedCh:
				systray.Quit()
			case <-meeting.ClickedCh:
				meeting.Check()
			}
		}
	}()
}

func onExit() {
	// Cleaning stuff here.
	fmt.Println("bye")
}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}
	return b
}

func getClockTime() string {
	t := time.Now()

	hour, min, sec := t.Clock()
	return ItoaTwoDigits(hour) + ":" + ItoaTwoDigits(min) + ":" + ItoaTwoDigits(sec)
}

// ItoaTwoDigits time.Clock returns one digit on values, so we make sure to convert to two digits
func ItoaTwoDigits(i int) string {
	b := "0" + strconv.Itoa(i)
	return b[len(b)-2:]
}
