package goption

import "database/sql/driver"

// Scan assigns a value from a database driver.
//
// The src value will be of one of the following types:
//
//    int64
//    float64
//    bool
//    []byte
//    string
//    time.Time
//    nil - for NULL values
//
// An error should be returned if the value cannot be stored
// without loss of information.
//
// Reference types such as []byte are only valid until the next call to Scan
// and should not be retained. Their underlying memory is owned by the driver.
// If retention is necessary, copy their values before the next call to Scan.
func (c *Optional[T]) Scan(src any) error {
	isSrcNil := checkIsNil(src)
	c.isValueNil = isSrcNil
	if !isSrcNil {
		c.value = src.(T)
	}
	return nil
}

func (c Optional[T]) Value() (driver.Value, error) {
	if c.isValueNil {
		return nil, nil
	}
	return c.value, nil

}
