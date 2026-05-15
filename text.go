package goption

import (
	"bytes"
	"fmt"
	"regexp"
)

func (c *Optional[T]) UnmarshalText(text []byte) error {
	var isString bool
	switch any(c).(type) {
	case Optional[string], *Optional[string]:
		isString = true
	}
	isNumber, err := regexp.Match(`^\d+(\.\d+)?$`, text)
	if err != nil {
		return err
	}

	isArray := bytes.HasPrefix(text, []byte("["))

	isBool, _ := c.isValueBoolTypeAndPointer()
	if (isNumber && !isBool) && !isString || isArray {
		return c.UnmarshalJSON(text)
	}

	return c.UnmarshalJSON([]byte(fmt.Sprintf("\"%s\"", text)))
}
