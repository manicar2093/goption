package goption

import (
	"errors"
	"reflect"
)

var (
	ErrNoSuchElement = errors.New("no value in this optional")
)

type Goptional[T any] struct {
	value      T
	isValueNil bool
}

// Empty returns an empty Optional instance.
func Empty[T any]() Goptional[T] {
	return Goptional[T]{isValueNil: true}
}

// Of returns an Optional with the specified present value. It does not matters if value is nil
func Of[T any](value T) Goptional[T] {
	return Goptional[T]{value: value, isValueNil: checkIsNil(value)}
}

// Get when a value is present returns the value, otherwise throws ErrNoSuchElement.
func (c Goptional[T]) Get() (T, error) {
	if c.isValueNil {
		return c.value, ErrNoSuchElement
	}
	return c.value, nil
}

// IsPresent returns true if there is a value present, otherwise false.
func (c Goptional[T]) IsPresent() bool {
	return !c.isValueNil
}

// OrElseError return the contained value, if present, otherwise returns the given error.
func (c Goptional[T]) OrElseError(err error) (T, error) {
	if c.isValueNil {
		return c.value, err
	}
	return c.value, nil
}

// OrElse returns the value if present, otherwise return other.
func (c Goptional[T]) OrElse(other T) T {
	if c.isValueNil {
		return other
	}
	return c.value
}

func checkIsNil[T any](value T) bool {
	kind := reflect.TypeOf(value).Kind()
	val := reflect.ValueOf(value)
	if kind == reflect.Pointer {
		return val.IsNil()
	}
	return false
}
