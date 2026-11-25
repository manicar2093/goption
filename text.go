package goption

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func parseRelaxedBool(text string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(text)) {
	case "true", "1", "on", "yes":
		return true, nil
	case "false", "0", "off", "no", "":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean text %q", text)
	}
}

// MarshalText implements encoding.TextMarshaler for Optional[T].
func (o Optional[T]) MarshalText() ([]byte, error) {
	if !o.valid {
		return []byte(""), nil
	}
	if tm, ok := any(o.value).(encoding.TextMarshaler); ok {
		return tm.MarshalText()
	}
	return []byte(fmt.Sprint(o.value)), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for Optional[T].
func (o *Optional[T]) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		o.valid = false
		var zero T
		o.value = zero
		return nil
	}

	if tu, ok := any(&o.value).(encoding.TextUnmarshaler); ok {
		if err := tu.UnmarshalText(text); err != nil {
			return err
		}
		o.valid = true
		return nil
	}

	if _, ok := any(o.value).(*bool); ok {
		b, err := parseRelaxedBool(str)
		if err != nil {
			return err
		}
		o.value = any(&b).(T)
		o.valid = true
		return nil
	}

	// Handle primitive values
	v := reflect.ValueOf(&o.value).Elem()
	switch v.Kind() {
	case reflect.Bool:
		b, err := parseRelaxedBool(str)
		if err != nil {
			return err
		}
		v.SetBool(b)
		o.valid = true
		return nil
	case reflect.String:
		v.SetString(str)
		o.valid = true
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
		o.valid = true
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(u)
		o.valid = true
		return nil
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		v.SetFloat(f)
		o.valid = true
		return nil
	}

	return o.UnmarshalJSON(text)
}

// MarshalText implements encoding.TextMarshaler for Zero[T].
func (z Zero[T]) MarshalText() ([]byte, error) {
	if !z.present {
		return []byte(""), nil
	}
	if tm, ok := any(z.value).(encoding.TextMarshaler); ok {
		return tm.MarshalText()
	}
	return []byte(fmt.Sprint(z.value)), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for Zero[T].
func (z *Zero[T]) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		z.present = false
		var zero T
		z.value = zero
		return nil
	}

	if tu, ok := any(&z.value).(encoding.TextUnmarshaler); ok {
		if err := tu.UnmarshalText(text); err != nil {
			return err
		}
		z.present = !isZero(z.value)
		return nil
	}

	if _, ok := any(z.value).(*bool); ok {
		b, err := parseRelaxedBool(str)
		if err != nil {
			return err
		}
		z.value = any(&b).(T)
		z.present = !isZero(z.value)
		return nil
	}

	v := reflect.ValueOf(&z.value).Elem()
	switch v.Kind() {
	case reflect.Bool:
		b, err := parseRelaxedBool(str)
		if err != nil {
			return err
		}
		v.SetBool(b)
		z.present = !isZero(z.value)
		return nil
	case reflect.String:
		v.SetString(str)
		z.present = !isZero(z.value)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
		z.present = !isZero(z.value)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(u)
		z.present = !isZero(z.value)
		return nil
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		v.SetFloat(f)
		z.present = !isZero(z.value)
		return nil
	}

	if err := z.UnmarshalJSON(text); err != nil {
		return err
	}
	z.present = !isZero(z.value)
	return nil
}
