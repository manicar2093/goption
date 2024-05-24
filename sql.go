package goption

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

// Scan assigns a value from a database driver.
//
// The src value will be of one of the following types:
//
//	int64
//	float64
//	bool
//	[]byte
//	string
//	time.Time
//	nil - for NULL values
//
// An error should be returned if the value cannot be stored
// without loss of information.
//
// Reference types such as []byte are only valid until the next call to Scan
// and should not be retained. Their underlying memory is owned by the driver.
// If retention is necessary, copy their values before the next call to Scan.
func (c *Optional[T]) Scan(src any) error {
	srcValue, isSrcValid := isValidData(src)
	c.isValidValue = isSrcValid
	if isSrcValid {
		destType := reflect.TypeOf(c.value)
		if !srcValue.Type().ConvertibleTo(destType) {
			return fmt.Errorf("interface conversion: interface {} is %s, not %s", srcValue.Type(), destType)
		}

		c.value = srcValue.Convert(destType).Interface().(T)
		return nil
	}
	return nil
}

func (c Optional[T]) Value() (driver.Value, error) {
	if !c.isValidValue {
		return nil, nil
	}
	return c.value, nil

}
