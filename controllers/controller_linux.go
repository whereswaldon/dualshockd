package controllers

import (
	udev "github.com/jochenvg/go-udev"
	"path/filepath"
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
	return filepath.Base(c.device.Syspath())
}

func (c *LinuxController) String() string {
	var properties string
	for key, value := range c.device.Properties() {
		properties += key + ":\t" + value + "\n"
	}
	return c.Name() + "\n" + properties
}

var _ Controller = &LinuxController{}
