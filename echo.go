package goption

import (
	"fmt"
)

func (c *Optional[T]) UnmarshalText(text []byte) error {
	return c.UnmarshalJSON([]byte(fmt.Sprintf("\"%s\"", text)))

}
