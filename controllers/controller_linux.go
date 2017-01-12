package controllers

import (
	udev "github.com/jochenvg/go-udev"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// NewController creates a controller from a linux udev
// device.
func NewController(done <-chan struct{}, device *udev.Device) (Controller, error) {
	return &LinuxController{
		device: device,
	}, nil
}

// LinuxController wraps a udev device with convenience
// methods.
type LinuxController struct {
	device *udev.Device
}

// Charge returns the current charge in this controller's
// battery expressed as a percentage.
func (c *LinuxController) Charge() (uint, error) {
	charge, err := os.Open(c.device.Syspath() + "/power_supply")
	if err != nil {
		return 0, errors.Wrapf(err, "Unable to open device syspath")
	}
	contents, err := charge.Readdirnames(1)
	if err != nil {
		return 0, errors.Wrapf(err, "Unable to read device syspath contents")
	}
	charge, err = os.Open(c.device.Syspath() + "/power_supply/" + contents[0] + "/capacity")
	if err != nil {
		return 0, errors.Wrapf(err, "Unable to open device battery file")
	}
	chars := make([]byte, 3)
	numberRead, err := charge.Read(chars[:])
	if err != nil {
		return 0, errors.Wrapf(err, "Unable to read battery file")
	}
	stringified := strings.Replace(string(chars[:numberRead]), "\n", "", -1)
	value, err := strconv.Atoi(stringified)
	if err != nil {
		return uint(value), errors.Wrapf(err, "Unable to convert battery level to intger")
	}
	return uint(value), nil
}

// Name returns a unique identifier for this controller.
func (c *LinuxController) Name() string {
	return filepath.Base(c.device.Syspath())
}

// String serializes many of the properties of the controller
// for use in debugging. Its output varies by system and
// connection type.
func (c *LinuxController) String() string {
	var properties string
	for key, value := range c.device.Properties() {
		properties += key + ":\t" + value + "\n"
	}
	return c.Name() + "\n" + properties
}

var _ Controller = &LinuxController{}
