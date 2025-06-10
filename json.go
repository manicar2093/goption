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
		asString = string(data)
	)
	if strings.HasPrefix(asString, "\"") {
		return c.stringUnmarshall(asString)
	}

	return c.numberUnmarshal(asString)
}

func (c Optional[T]) MarshalJSON() ([]byte, error) {
	if reflect.ValueOf(c.value).IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(c.value)
}

func isNull(data string) bool {
	return data == "null"
}

func (c *Optional[T]) stringUnmarshall(data string) error {
	var (
		valuer struct {
			Value string `json:"value"`
		}
		cleanables = []struct {
			find        string
			replacement string
		}{
			{find: "\n", replacement: "\\n"},
			{find: "\r", replacement: "\\r"},
		}
		asJsonString = fmt.Sprintf(`{"value": %s}`, data)
	)
	for _, item := range cleanables {
		asJsonString = strings.ReplaceAll(asJsonString, item.find, item.replacement)
	}

	if err := json.Unmarshal([]byte(asJsonString), &valuer); err != nil {
		return err
	}

	if isNull(valuer.Value) {
		valuer.Value = ""
	}

	if reflect.ValueOf(c.value).Kind() == reflect.Bool {
		switch valuer.Value {
		case "true", "1", "on", "yes":
			valuer.Value = "true"
		case "false", "0", "off", "no":
			valuer.Value = "false"
		}
		c.isValidValue = getIsValidDataBool(valuer.Value)
		return c.unmarshallIntoValueIfValid([]byte(valuer.Value))
	}

	c.isValidValue = getIsValidDataBool(valuer.Value)
	return c.unmarshallIntoValueIfValid([]byte(strconv.Quote(valuer.Value)))
}

func (c *Optional[T]) numberUnmarshal(data string) error {
	if isNull(data) {
		data = ""
	}

	c.isValidValue = getIsValidDataBool(data)
	return c.unmarshallIntoValueIfValid([]byte(data))
}

func (c *Optional[T]) unmarshallIntoValueIfValid(data []byte) error {
	if !c.isValidValue {
		return nil
	}
	return json.Unmarshal(
		data,
		&c.value,
	)
}
