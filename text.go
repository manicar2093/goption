package goption

import (
	"fmt"
	"regexp"
)

func (c *Optional[T]) UnmarshalText(text []byte) error {
	does, err := regexp.Match(`^\d+(\.\d+)?$`, text)
	if err != nil {
		return err
	}
	if does {
		return c.UnmarshalJSON(text)
	}
	return c.UnmarshalJSON([]byte(fmt.Sprintf("\"%s\"", text)))
}
