package goption

func (c *Optional[T]) UnmarshalText(text []byte) error {
	return c.UnmarshalJSON(text)
}
