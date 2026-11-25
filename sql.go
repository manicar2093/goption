package goption

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
)

// Scan implements database/sql.Scanner for Optional[T].
// Only a nil src is treated as SQL NULL.
func (o *Optional[T]) Scan(src any) error {
	if src == nil {
		o.valid = false
		var zero T
		o.value = zero
		return nil
	}

	if scanner, ok := any(&o.value).(sql.Scanner); ok {
		if err := scanner.Scan(src); err != nil {
			return err
		}
		o.valid = true
		return nil
	}

	srcVal := reflect.ValueOf(src)
	destType := reflect.TypeOf(o.value)
	if destType != nil && srcVal.Type().ConvertibleTo(destType) {
		o.value = srcVal.Convert(destType).Interface().(T)
		o.valid = true
		return nil
	}

	if destType == nil {
		o.valid = true
		return nil
	}

	return fmt.Errorf("cannot scan %T into %T", src, o.value)
}

// Value implements driver.Valuer for Optional[T].
func (o Optional[T]) Value() (driver.Value, error) {
	if !o.valid {
		return nil, nil
	}
	if valuer, ok := any(o.value).(driver.Valuer); ok {
		return valuer.Value()
	}
	return driver.DefaultParameterConverter.ConvertValue(o.value)
}

// Scan implements database/sql.Scanner for Zero[T].
func (z *Zero[T]) Scan(src any) error {
	if src == nil {
		z.present = false
		var zero T
		z.value = zero
		return nil
	}

	if scanner, ok := any(&z.value).(sql.Scanner); ok {
		if err := scanner.Scan(src); err != nil {
			return err
		}
		z.present = !isZero(z.value)
		return nil
	}

	srcVal := reflect.ValueOf(src)
	destType := reflect.TypeOf(z.value)
	if destType != nil && srcVal.Type().ConvertibleTo(destType) {
		z.value = srcVal.Convert(destType).Interface().(T)
		z.present = !isZero(z.value)
		return nil
	}

	if destType == nil {
		z.present = false
		return nil
	}

	return fmt.Errorf("cannot scan %T into %T", src, z.value)
}

// Value implements driver.Valuer for Zero[T].
func (z Zero[T]) Value() (driver.Value, error) {
	if !z.present {
		return nil, nil
	}
	if valuer, ok := any(z.value).(driver.Valuer); ok {
		return valuer.Value()
	}
	return driver.DefaultParameterConverter.ConvertValue(z.value)
}
