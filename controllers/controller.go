package controllers

// Controller represents a single ps3 controller connected to the
// current system.
type Controller interface {
	Charge() (uint, error)
	Name() string
}
