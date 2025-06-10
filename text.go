package goption

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
)

func (c *Optional[T]) UnmarshalText(text []byte) error {
	isNumber, err := regexp.Match(`^\d+(\.\d+)?$`, text)
	isArray := bytes.HasPrefix(text, []byte("["))

	if err != nil {
		return err
	}
	if (isNumber && reflect.ValueOf(c.value).Kind() != reflect.Bool) || isArray {
		return c.UnmarshalJSON(text)
	}

	return c.UnmarshalJSON([]byte(fmt.Sprintf("\"%s\"", text)))
}
