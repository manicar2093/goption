package goption

import (
	"fmt"
)

func (cb *Optional[T]) UnmarshalText(text []byte) error {
	return cb.UnmarshalJSON([]byte(fmt.Sprintf("\"%s\"", text)))

}
