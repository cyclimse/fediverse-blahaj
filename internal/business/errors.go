package business

import "errors"

var (
	// ErrInstanceNotFound is returned when a instance is not found
	ErrInstanceNotFound = errors.New("instance not found")
)
