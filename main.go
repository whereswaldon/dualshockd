package main

import (
	"fmt"
	"github.com/0xAX/notificator"
	"github.com/pkg/profile"
	"github.com/whereswaldon/dualshockd/controllers"
	"log"
	"os"
	"os/signal"
	"time"
)

var notify *notificator.Notificator = notificator.New(notificator.Options{
	AppName: "Dualshockd",
})

// watchController observes the battery status of a controller
// and emits updates on it until the done channel is closed.
func watchController(done <-chan struct{}, c controllers.Controller) {
	var lastCharge uint
	for range time.NewTicker(time.Minute).C {
		select {
		case <-done:
			return
		default:
			charge, err := c.Charge()
			if err != nil {
				log.Println(err)
			} else {
				if charge != lastCharge {
					notify.Push(c.Name(), fmt.Sprintf("Battery is now at %d%%", charge), "", notificator.UR_NORMAL)
				}
				log.Printf("%s: Battery %d%%\n", c.Name(), charge)
			}
			lastCharge = charge
		}
	}
}

func main() {
	defer profile.Start(profile.MemProfile).Stop()
	done := make(chan struct{})
	defer close(done)
	m, err := controllers.NewMonitor(done)
	if err != nil {
		log.Println(err)
		return
	}
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)
	for c := range m.Controllers() {
		select {
		case <-interrupt:
			log.Println("Recieved SIGINT, exiting...")
			return
		default:
			log.Printf("Found controller %s\n", c.Name())
			go watchController(done, c)
		}
	}
}
