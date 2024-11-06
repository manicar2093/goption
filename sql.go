package goption

import (
	"database/sql"
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
	if !isSrcValid {
		return nil
	}

	destType := reflect.TypeOf(c.value)
	if srcValue.Type().ConvertibleTo(destType) {
		c.value = srcValue.Convert(destType).Interface().(T)
		return nil
	}

	destTypeP := reflect.TypeOf(&c.value)
	scannerType := reflect.TypeOf((*sql.Scanner)(nil))
	if destTypeP.Implements(scannerType.Elem()) {
		var s T
		if asScanner, ok := interface{}(&s).(sql.Scanner); ok {
			if err := asScanner.Scan(src); err != nil {
				return err
			}
			c.value = s
			return nil
		}
	}

	return fmt.Errorf("interface conversion: interface {} is %s, not %s nor implements sql.Scanner", srcValue.Type(), destType)
}

func (c Optional[T]) Value() (driver.Value, error) {
	if !c.isValidValue {
		return nil, nil
	}
	return c.value, nil

}
