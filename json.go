package goption

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func (c *Optional[T]) UnmarshalJSON(data []byte) error {
	var (
		valuer struct {
			Value string `json:"value"`
		}
		asJsonString = strings.ReplaceAll(
			fmt.Sprintf(`{"value": %s}`, data),
			"\n",
			"\\n",
		)
		asJsonBytes = []byte(asJsonString)
	)

	if err := json.Unmarshal(asJsonBytes, &valuer); err != nil {
		return err
	}

	if valuer.Value == "null" {
		valuer.Value = ""
	}

	c.isValidValue = getIsValidDataBool(valuer.Value)
	if c.isValidValue {
		if err := json.Unmarshal(
			[]byte(strconv.Quote(valuer.Value)),
			&c.value,
		); err != nil {
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
