package controllers

// Monitor watches the system for controllers and manages their creation
// and status.
type Monitor interface {
	Controllers() <-chan Controller
}
