package business

import "errors"

var (
	// ErrServerNotFound is returned when a server is not found
	ErrServerNotFound = errors.New("server not found")
)
