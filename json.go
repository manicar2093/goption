package goption

import (
	"encoding/json"
	"log"
)

func (c *Optional[T]) UnmarshalJSON(data []byte) error {
	log.Println(string(data) == `"null"`, string(data), "null")
	if string(data) == `"null"` {
		c.isValueNil = true
		return nil
	}
	c.isValueNil = len(data) <= 0
	if err := json.Unmarshal(data, &c.value); err != nil {
		return err
	}
	return nil
}
