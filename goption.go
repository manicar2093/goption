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
	return Optional[T]{value: value, isValidValue: getIsValidDataBool(value)}
}

// Get when a value is present returns the value, otherwise throws ErrNoSuchElement.
func (c Optional[T]) Get() (T, error) {
	if !c.isValidValue {
		return c.value, ErrNoSuchElement
	}
	return c.value, nil
}

// IsPresent returns true if there is a value present, otherwise false. It recognizes an empty slice as not present so it returns false
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

func (c Optional[T]) IsZero() bool {
	return !c.isValidValue
}

func (c Optional[T]) GetValue() T {
	return c.value
}

func isValidData[T any](value T) (reflect.Value, bool) {
	typeOfValue := reflect.TypeOf(value)
	if typeOfValue == nil {
		return reflect.Value{}, false
	}

	val := reflect.ValueOf(value)

	switch typeOfValue.Kind() {
	case reflect.Pointer:
		return val, !val.IsNil()
	case reflect.Slice:
		return val, val.Len() != 0
	default:
		return val, !val.IsZero()
	}
}

func getIsValidDataBool[T any](value T) bool {
	_, is := isValidData(value)
	return is
}
