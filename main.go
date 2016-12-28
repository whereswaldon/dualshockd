package main

import (
	"fmt"
	"github.com/pkg/profile"
	"github.com/whereswaldon/dualshockd/controllers"
	"os"
	"os/signal"
)

func main() {
	defer profile.Start(profile.MemProfile).Stop()
	done := make(chan struct{})
	defer close(done)
	m, err := controllers.NewMonitor(done)
	if err != nil {
		fmt.Println(err)
		return
	}
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)
	for c := range m.Controllers() {
		select {
		case <-interrupt:
			fmt.Println("exiting")
			return
		default:
			fmt.Println(c)
		}
	}
}
