package controllers

import (
	udev "github.com/jochenvg/go-udev"
)

// NewController creates a controller from a linux udev
// device.
func NewController(device *udev.Device) Controller {
	disconnected := make(chan struct{})
	battery := make(chan uint)
	return &LinuxController{
		device:        device,
		disconnected:  disconnected,
		batteryChange: battery,
	}

}

// LinuxController wraps a udev device with convenience
// methods.
type LinuxController struct {
	device        *udev.Device
	disconnected  chan struct{}
	batteryChange chan uint
}

// Disconnected returns a channel that will be closed when
// the controller is disconnected.
func (c *LinuxController) Disconnected() <-chan struct{} {
	return c.disconnected
}

// BatteryChanges returns a channel that will receive an
// update each time the battery status of the controller
// changes.
func (c *LinuxController) BatteryChanges() <-chan uint {
	return c.batteryChange
}

// Name returns a unique identifier for this controller.
func (c *LinuxController) Name() string {
	return ""
}

var _ Controller = &LinuxController{}
