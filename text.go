package goption

import (
	"fmt"
	"regexp"
)

func (c *Optional[T]) UnmarshalText(text []byte) error {
	//fse := string(text[0])
	//if slices.Contains([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}, fse) {
	//	return c.UnmarshalJSON(text)
	//}
	does, err := regexp.Match(`\d`, text)
	if err != nil {
		return err
	}
	if does {
		return c.UnmarshalJSON(text)
	}
	return c.UnmarshalJSON([]byte(fmt.Sprintf("\"%s\"", text)))
}
