package goption

import (
	"bytes"
	"encoding/json"
)

// MarshalJSON implements json.Marshaler for Optional[T].
// It returns "null" if the value is invalid/null, or serializes the contained value.
func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.valid {
		return []byte("null"), nil
	}
	return json.Marshal(o.value)
}

// UnmarshalJSON implements json.Unmarshaler for Optional[T].
// Only a literal JSON "null" marks the container invalid/absent.
// Zero values like 0, "", false, or empty objects/arrays remain valid present values.
func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		o.valid = false
		var zero T
		o.value = zero
		return nil
	}
	if err := json.Unmarshal(data, &o.value); err != nil {
		return err
	}
	o.valid = true
	return nil
}

// MarshalJSON implements json.Marshaler for Zero[T].
// It returns "null" if the value is absent or zero, or serializes the contained value.
func (z Zero[T]) MarshalJSON() ([]byte, error) {
	if !z.present {
		return []byte("null"), nil
	}
	return json.Marshal(z.value)
}

// UnmarshalJSON implements json.Unmarshaler for Zero[T].
// If data is "null" or unmarshals into a zero value, present is set to false.
func (z *Zero[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		z.present = false
		var zero T
		z.value = zero
		return nil
	}
	if err := json.Unmarshal(data, &z.value); err != nil {
		return err
	}
	z.present = !isZero(z.value)
	return nil
}
