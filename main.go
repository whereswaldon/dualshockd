package main

import (
	"fmt"
	"github.com/whereswaldon/dualshockd/controllers"
)

func main() {
	done := make(chan struct{})
	defer close(done)
	m, err := controllers.NewMonitor(done)
	if err != nil {
		fmt.Println(err)
		return
	}
	for c := range m.Controllers() {
		fmt.Println(c)
	}
}
