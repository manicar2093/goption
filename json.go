package goption

import (
	"encoding/json"
	"reflect"
	"strconv"
)

func (c *Optional[T]) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	if unquoted == "null" {
		unquoted = ""
	}
	_, isValid := isValidData(unquoted)
	c.isValidValue = isValid
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
