package goption

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

func (c *Optional[T]) UnmarshalJSON(data []byte) error {
	var (
		err      error
		asString = string(data)
	)
	if strings.HasPrefix(asString, "\"") {
		asString, err = strconv.Unquote(asString)
		if err != nil {
			return err
		}
	}
	if asString == "null" {
		asString = ""
	}
	c.isValidValue = getIsValidDataBool(asString)
	if c.isValidValue {
		if err := json.Unmarshal(data, &c.value); err != nil {
			return err
		}
	}

	return nil
}

func (c Optional[T]) MarshalJSON() ([]byte, error) {
	if reflect.ValueOf(c.value).IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(c.value)
}
