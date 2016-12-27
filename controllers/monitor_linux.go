package controllers

import (
	udev "github.com/jochenvg/go-udev"
	"github.com/pkg/errors"
)

// LinuxMonitor fulfills the Monitor interface on Linux
// by watching sysfs with libudev
type LinuxMonitor struct {
	controllers <-chan Controller
}

// NewMonitor creates a monitor that watches for Controllers
// to be added to the system and creates Controller structs
// for them when they are. When the provided done channel is
// closed, the monitor will stop watching.
func NewMonitor(done <-chan struct{}) (Monitor, error) {
	controllers := make(chan Controller)
	u := &udev.Udev{}
	m := u.NewMonitorFromNetlink("udev") //handle events after kernel
	m.FilterAddMatchSubsystem("input")
	events, err := m.DeviceChan(done)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create LinuxMonitor")
	}
	go func() {
		// look at each input device event
		for event := range events {
			// if the event is about a joystick
			if dev.Properties()["ID_INPUT_JOYSTICK"] == "1" {
				// find the "hid" subsystem parent and pass that on.
				for parent := dev.Parent(); parent != nil; parent = parent.Parent() {
					if parent.Subsystem() == "hid" {
						controllers <- NewController(parent)
					}
				}

			}
		}
	}()
	return &LinuxMonitor{
		controllers: controllers,
	}, nil
}

// Controllers returns a channel along which each new Controller
// connected to the system will be passed.
func (m *LinuxMonitor) Controllers() <-chan Controller {
	return nil
}

var _ Monitor = &LinuxMonitor{}
