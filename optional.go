package goption

import "reflect"

// Optional represents an optional value where zero values (0, "", false, empty struct) are considered valid.
// It is intended for SQL and JSON nullability where null/absence must be distinguished from zero values.
type Optional[T any] struct {
	value T
	valid bool
}

func isNil[T any](val T) bool {
	v := reflect.ValueOf(val)
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface, reflect.Slice, reflect.Map, reflect.Func, reflect.Chan, reflect.UnsafePointer:
		return v.IsNil()
	default:
		return false
	}
}

// Some creates an Optional[T] containing val.
// If val is nil (such as a nil pointer, slice, map, channel, or interface), the returned Optional[T] is marked invalid/null.
// All other values, including zero values like 0, "", false, or empty structs, are considered valid present values.
func Some[T any](val T) Optional[T] {
	if isNil(val) {
		return Optional[T]{
			value: val,
			valid: false,
		}
	}
	return Optional[T]{
		value: val,
		valid: true,
	}
}

// IsValid returns true if the value is valid (non-null).
func (o Optional[T]) IsValid() bool {
	return o.valid
}

// IsPresent returns true if the value is present/valid.
func (o Optional[T]) IsPresent() bool {
	return o.valid
}

// IsZero returns true if the value is absent/null.
// This allows Optional[T] to work seamlessly with the Go 1.24 omitzero JSON tag.
func (o Optional[T]) IsZero() bool {
	return !o.valid
}

// Get returns the contained value and true if valid, or the zero value of T and false if null.
func (o Optional[T]) Get() (T, bool) {
	if !o.valid {
		var zero T
		return zero, false
	}
	return o.value, true
}

// OrElse returns the contained value if valid, or fallback if null.
func (o Optional[T]) OrElse(fallback T) T {
	if !o.valid {
		return fallback
	}
	return o.value
}

// OrElseZero returns the contained value if valid, or the zero value of T if null.
func (o Optional[T]) OrElseZero() T {
	if !o.valid {
		var zero T
		return zero
	}
	return o.value
}

// MustGet returns the contained value if valid, or panics with ErrNoSuchElement if null.
func (o Optional[T]) MustGet() T {
	if !o.valid {
		panic(ErrNoSuchElement)
	}
	return o.value
}

// Filter returns o if it is valid and satisfies predicate; otherwise returns an invalid Optional[T]{}.
func (o Optional[T]) Filter(predicate func(T) bool) Optional[T] {
	return Filter(o, predicate)
}

// IfPresent executes action with the contained value if o is valid.
func (o Optional[T]) IfPresent(action func(T)) {
	IfPresent(o, action)
}

// IfPresentOrElse executes action if o is valid, otherwise executes emptyAction.
func (o Optional[T]) IfPresentOrElse(action func(T), emptyAction func()) {
	IfPresentOrElse(o, action, emptyAction)
}

// Or returns o if valid, otherwise returns fallback.
func (o Optional[T]) Or(fallback Optional[T]) Optional[T] {
	return Or(o, fallback)
}
