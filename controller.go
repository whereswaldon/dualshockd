package main

import "fmt"

// Controller represents a single ps3 controller connected to the
// current system.
//
// Disconnected returns a channel that will be closed when the controller
// is disconnected from the system.
//
// BatteryChanges returns a channel of battery levels. A new level will
// be sent each time the battery status changes.
//
// Name returns a unique identifier for this controller.
type Controller interface {
	Disconnected() <-chan struct{}
	BatteryChanges() <-chan uint
	Name() string
}
