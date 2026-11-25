package goption

import "reflect"

// Zero represents an optional value where zero values (0, "", false, nil) are treated as empty/absent.
// It is intended for presence semantics where default zero values indicate absence.
type Zero[T any] struct {
	value   T
	present bool
}

func isZero[T any](val T) bool {
	v := reflect.ValueOf(val)
	if !v.IsValid() {
		return true
	}
	return v.IsZero()
}

// SomeZero creates a Zero[T] containing val.
// If val is a zero value (such as 0, "", false, or nil), the returned Zero[T] is marked empty (present = false).
// If val is non-zero, the returned Zero[T] is marked present (present = true).
func SomeZero[T any](val T) Zero[T] {
	if isZero(val) {
		return Zero[T]{
			value:   val,
			present: false,
		}
	}
	return Zero[T]{
		value:   val,
		present: true,
	}
}

// IsPresent returns true if the value is non-zero and present.
func (z Zero[T]) IsPresent() bool {
	return z.present
}

// IsZero returns true if the value is empty or zero.
// This allows Zero[T] to work seamlessly with the Go 1.24 omitzero JSON tag.
func (z Zero[T]) IsZero() bool {
	return !z.present
}

// Get returns the contained value and true if present, or the zero value of T and false if empty.
func (z Zero[T]) Get() (T, bool) {
	if !z.present {
		var zero T
		return zero, false
	}
	return z.value, true
}

// OrElse returns the contained value if present, or fallback if empty.
func (z Zero[T]) OrElse(fallback T) T {
	if !z.present {
		return fallback
	}
	return z.value
}

// OrElseZero returns the contained value if present, or the zero value of T if empty.
func (z Zero[T]) OrElseZero() T {
	if !z.present {
		var zero T
		return zero
	}
	return z.value
}

// MustGet returns the contained value if present, or panics with ErrNoSuchElement if empty.
func (z Zero[T]) MustGet() T {
	if !z.present {
		panic(ErrNoSuchElement)
	}
	return z.value
}

// Filter returns z if it is present and satisfies predicate; otherwise returns an empty Zero[T]{}.
func (z Zero[T]) Filter(predicate func(T) bool) Zero[T] {
	return FilterZero(z, predicate)
}

// IfPresent executes action with the contained value if z is present.
func (z Zero[T]) IfPresent(action func(T)) {
	IfPresentZero(z, action)
}

// IfPresentOrElse executes action if z is present, otherwise executes emptyAction.
func (z Zero[T]) IfPresentOrElse(action func(T), emptyAction func()) {
	IfPresentOrElseZero(z, action, emptyAction)
}

// Or returns z if present, otherwise returns fallback.
func (z Zero[T]) Or(fallback Zero[T]) Zero[T] {
	return OrZero(z, fallback)
}
