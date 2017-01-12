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

// stringSet manipulates a set of strings
type stringSet map[string]struct{}

// Add adds a new value to the set.
func (s stringSet) Add(value string) {
	s[value] = struct{}{}
}

// Remove takes a value out of the set.
func (s stringSet) Remove(value string) {
	delete(s, value)
}

// Contains returns true if the value is in the set.
func (s stringSet) Contains(value string) bool {
	_, ok := s[value]
	return ok
}

// isValidJoystick returns whether the device in question is an initialized joystick.
func isValidJoystick(device *udev.Device) bool {
	return device != nil && device.IsInitialized() && device.Properties()["ID_INPUT_JOYSTICK"] == "1"
}

// getHIDParent returns the parent of the provided device that is in the HID subsystem,
// or nil if there is no parent in that subsystem.
func getHIDParent(device *udev.Device) *udev.Device {
	for parent := device.Parent(); parent != nil; parent = parent.Parent() {
		if parent.Subsystem() == "hid" {
			return parent
		}
	}
	return nil
}

func stringifyDevice(device *udev.Device) string {
	var result string
	result += device.Devpath() + "\n"
	result += device.Syspath() + "\n"
	result += device.Subsystem() + "\n"
	for key, value := range device.Properties() {
		result += key + ":\t" + value + "\n"
	}
	if isValidJoystick(device) {
		result += "validJS\n"
	}
	if p := getHIDParent(device); p != nil {
		result += "PARENT:\t" + stringifyDevice(p)
	}
	return result
}

// NewMonitor creates a monitor that watches for Controllers
// to be added to the system and creates Controller structs
// for them when they are. When the provided done channel is
// closed, the monitor will stop watching.
func NewMonitor(done <-chan struct{}) (Monitor, error) {
	controllers := make(chan Controller)
	passDevices := make(chan *udev.Device)
	monitored := make(stringSet)
	passIfValidJoystick := func(device *udev.Device) {
		// if the device is about a joystick
		if isValidJoystick(device) {
			// find the "hid" subsystem parent and pass that on.
			hidParent := getHIDParent(device)
			if hidParent != nil && !monitored.Contains(hidParent.Properties()["HID_UNIQ"]) {
				controller, err := NewController(done, hidParent)
				if err != nil {
					//TODO: handle error?
				}
				controllers <- controller
				monitored.Add(hidParent.Properties()["HID_UNIQ"])
			}
		}
	}

	u := &udev.Udev{}
	m := u.NewMonitorFromNetlink("udev") //handle events after kernel
	m.FilterAddMatchSubsystem("input")
	newDevices, err := m.DeviceChan(done)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create LinuxMonitor")
	}

	// watch for new controllers
	go func() {
		defer close(controllers)
		// look at each input device event
		for {
			select {
			case <-done:
				return
			case device := <-passDevices:
				if device == nil {
					// stop reading this channel
					passDevices = nil
				}
				passIfValidJoystick(device)
			case device := <-newDevices:
				passIfValidJoystick(device)
			}
		}
	}()

	// find existing devices
	e := u.NewEnumerate()
	e.AddMatchSubsystem("input")
	e.AddMatchIsInitialized()
	oldDevices, err := e.Devices()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create LinuxMonitor")
	}
	go func() {
		defer close(passDevices)
		for _, device := range oldDevices {
			passDevices <- device
		}
	}()

	return &LinuxMonitor{
		controllers: controllers,
	}, nil
}

// Controllers returns a channel along which each new Controller
// connected to the system will be passed.
func (m *LinuxMonitor) Controllers() <-chan Controller {
	return m.controllers
}

var _ Monitor = &LinuxMonitor{}
