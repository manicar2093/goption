package goption

import "errors"

var (
	// ErrNoSuchElement is returned or panicked with when attempting to access a value from an empty Optional/Zero.
	ErrNoSuchElement = errors.New("no value present")
)
