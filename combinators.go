package goption

// Map applies fn to the value inside Optional[T] if valid, returning an Optional[U].
// If o is invalid/null, it returns an invalid Optional[U]{}.
func Map[T, U any](o Optional[T], fn func(T) U) Optional[U] {
	if !o.valid {
		return Optional[U]{}
	}
	return Some(fn(o.value))
}

// FlatMap applies fn to the value inside Optional[T] if valid, returning the resulting Optional[U].
// If o is invalid/null, it returns an invalid Optional[U]{}.
func FlatMap[T, U any](o Optional[T], fn func(T) Optional[U]) Optional[U] {
	if !o.valid {
		return Optional[U]{}
	}
	return fn(o.value)
}

// Filter returns o if it is valid and satisfies predicate; otherwise returns an invalid Optional[T]{}.
func Filter[T any](o Optional[T], predicate func(T) bool) Optional[T] {
	if !o.valid || !predicate(o.value) {
		return Optional[T]{}
	}
	return o
}

// IfPresent executes action with the contained value if o is valid.
func IfPresent[T any](o Optional[T], action func(T)) {
	if o.valid {
		action(o.value)
	}
}

// IfPresentOrElse executes action if o is valid, otherwise executes emptyAction.
func IfPresentOrElse[T any](o Optional[T], action func(T), emptyAction func()) {
	if o.valid {
		action(o.value)
	} else {
		emptyAction()
	}
}

// Or returns o if valid, otherwise returns fallback.
func Or[T any](o Optional[T], fallback Optional[T]) Optional[T] {
	if o.valid {
		return o
	}
	return fallback
}

// MapZero applies fn to the value inside Zero[T] if present, returning a Zero[U].
// If z is empty/absent, it returns an empty Zero[U]{}.
func MapZero[T, U any](z Zero[T], fn func(T) U) Zero[U] {
	if !z.present {
		return Zero[U]{}
	}
	return SomeZero(fn(z.value))
}

// FlatMapZero applies fn to the value inside Zero[T] if present, returning the resulting Zero[U].
// If z is empty/absent, it returns an empty Zero[U]{}.
func FlatMapZero[T, U any](z Zero[T], fn func(T) Zero[U]) Zero[U] {
	if !z.present {
		return Zero[U]{}
	}
	return fn(z.value)
}

// FilterZero returns z if it is present and satisfies predicate; otherwise returns an empty Zero[T]{}.
func FilterZero[T any](z Zero[T], predicate func(T) bool) Zero[T] {
	if !z.present || !predicate(z.value) {
		return Zero[T]{}
	}
	return z
}

// IfPresentZero executes action with the contained value if z is present.
func IfPresentZero[T any](z Zero[T], action func(T)) {
	if z.present {
		action(z.value)
	}
}

// IfPresentOrElseZero executes action if z is present, otherwise executes emptyAction.
func IfPresentOrElseZero[T any](z Zero[T], action func(T), emptyAction func()) {
	if z.present {
		action(z.value)
	} else {
		emptyAction()
	}
}

// OrZero returns z if present, otherwise returns fallback.
func OrZero[T any](z Zero[T], fallback Zero[T]) Zero[T] {
	if z.present {
		return z
	}
	return fallback
}
