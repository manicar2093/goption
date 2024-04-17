package goption

import (
	"errors"
	"reflect"
)

var (
	ErrNoSuchElement = errors.New("no value in this optional")
)

type Optional[T any] struct {
	value        T
	isValidValue bool
}

// Empty returns an empty Optional instance.
func Empty[T any]() Optional[T] {
	return Optional[T]{}
}

// Of returns an Optional with the specified present value. It does not matters if value is nil
func Of[T any](value T) Optional[T] {
	return Optional[T]{value: value, isValidValue: isValidData(value)}
}

// Get when a value is present returns the value, otherwise throws ErrNoSuchElement.
func (c Optional[T]) Get() (T, error) {
	if !c.isValidValue {
		return c.value, ErrNoSuchElement
	}
	return c.value, nil
}

// IsPresent returns true if there is a value present, otherwise false.
func (c Optional[T]) IsPresent() bool {
	return c.isValidValue
}

// OrElseError return the contained value, if present, otherwise returns the given error.
func (c Optional[T]) OrElseError(err error) (T, error) {
	if !c.isValidValue {
		return c.value, err
	}
	return c.value, nil
}

// OrElse returns the value if present, otherwise return other.
func (c Optional[T]) OrElse(other T) T {
	if !c.isValidValue {
		return other
	}
	return c.value
}

// MustGet retrieves only a valid value. If is not present it panics with ErrNoSuchElement
func (c Optional[T]) MustGet() T {
	val, err := c.Get()
	if err != nil {
		panic(err)
	}
	return val
}

func isValidData[T any](value T) bool {
	typeOfValue := reflect.TypeOf(value)
	if typeOfValue == nil {
		return false
	}
	kind := typeOfValue.Kind()
	val := reflect.ValueOf(value)

	switch kind {
	case reflect.Pointer:
		return !val.IsNil()
	case reflect.Slice:
		return val.Len() != 0
	default:
		return !val.IsZero()
	}
}
