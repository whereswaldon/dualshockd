package main

import (
	"github.com/pkg/profile"
	"github.com/whereswaldon/dualshockd/controllers"
	"log"
	"os"
	"os/signal"
	"time"
)

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
			go func(c controllers.Controller) {
				for range time.NewTicker(time.Minute).C {
					charge, err := c.Charge()
					if err != nil {
						log.Println(err)
					} else {
						log.Printf("%s: Battery %d%%\n", c.Name(), charge)
					}
				}
			}(c)
		}
	}
}
