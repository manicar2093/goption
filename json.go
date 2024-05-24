package goption

import (
	"encoding/json"
	"reflect"
)

func (c *Optional[T]) UnmarshalJSON(data []byte) error {
	asString := string(data)
	if asString == `"null"` || asString == `""` {
		c.isValidValue = false
		return nil
	}
	c.isValidValue = len(data) > 0
	if err := json.Unmarshal(data, &c.value); err != nil {
		return err
	}
	return nil
}

func (c Optional[T]) MarshalJSON() ([]byte, error) {
	if reflect.ValueOf(c.value).IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(c.value)
}
