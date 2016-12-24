package controllers

// LinuxMonitor fulfills the Monitor interface on Linux
// by watching sysfs with libudev
type LinuxMonitor struct {
}

// NewMonitor creates a monitor that watches for Controllers
// to be added to the system and creates Controller structs
// for them when they are.
func NewMonitor() Monitor {
	return &LinuxMonitor{}
}

// Controllers returns a channel along which each new Controller
// connected to the system will be passed.
func (m *LinuxMonitor) Controllers() <-chan Controller {
	return nil
}

var _ Monitor = &LinuxMonitor{}
