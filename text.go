package goption

import (
	"bytes"
	"fmt"
	"regexp"
)

func (c *Optional[T]) UnmarshalText(text []byte) error {
	isNumber, err := regexp.Match(`^\d+(\.\d+)?$`, text)
	isArray := bytes.HasPrefix(text, []byte("["))
	
	if err != nil {
		return err
	}
	if isNumber || isArray {
		return c.UnmarshalJSON(text)
	}

	return c.UnmarshalJSON([]byte(fmt.Sprintf("\"%s\"", text)))
}
