package goption_test

import (
	"errors"
	"testing"

	"github.com/manicar2093/goption"
	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	t.Run("creates an empty optional", func(t *testing.T) {
		assert.NotPanics(t, func() {
			goption.Empty[string]()
		})
	})
}

func TestOf(t *testing.T) {
	t.Run("creates an optional with no pointer given value", func(t *testing.T) {
		var expectedValue string = "hello"

		assert.NotPanics(t, func() {
			goption.Of(expectedValue)
		})
	})
	t.Run("creates an optional with pointer given value", func(t *testing.T) {
		var expectedValue string = "hello"

		assert.NotPanics(t, func() {
			goption.Of(&expectedValue)
		})
	})
	t.Run("creates an optional with nil pointer as given value", func(t *testing.T) {
		var expectedValue *string = nil

		assert.NotPanics(t, func() {
			goption.Of(expectedValue)
		})
	})
}

func TestIsPresent(t *testing.T) {
	t.Run("when has a valid data returns true", func(t *testing.T) {
		assert.False(t, goption.Empty[string]().IsPresent())
	})
	t.Run("when has a nil value returns false", func(t *testing.T) {
		assert.True(t, goption.Of("Hello, I'm present").IsPresent())
	})
	t.Run("when has a valid pointer value returns false", func(t *testing.T) {
		var expectedValue = "Hello, I'm present"

		assert.True(t, goption.Of(&expectedValue).IsPresent())
	})
}

func TestGetFromOptional(t *testing.T) {
	t.Run("when optional has a valid value", func(t *testing.T) {
		var (
			expectedValue string = "a valid value"
			opt                  = goption.Of(expectedValue)
			got, _               = opt.Get()
		)

		assert.Equal(t, expectedValue, got)
	})
	t.Run("when optional has a nil value", func(t *testing.T) {
		got, err := goption.Empty[string]().Get()

		assert.Empty(t, got)
		assert.Equal(t, goption.ErrNoSuchElement, err)
	})
}

func TestOrElseError(t *testing.T) {
	t.Run("if optional is empty returns given error", func(t *testing.T) {
		var (
			opt           = goption.Empty[string]()
			expectedError = errors.New("expected error")
		)

		got, err := opt.OrElseError(expectedError)

		assert.Empty(t, got)
		assert.Equal(t, expectedError, err)
	})
	t.Run("if optional has value returns it", func(t *testing.T) {
		var (
			expectedValue = "expectedValue"
			opt           = goption.Of(expectedValue)
			expectedError = errors.New("expected error")
		)

		got, err := opt.OrElseError(expectedError)

		assert.Equal(t, expectedValue, got)
		assert.Nil(t, err)
	})
}

func TestOrElse(t *testing.T) {
	t.Run("if optional is empty returns given error", func(t *testing.T) {
		var (
			opt               = goption.Empty[string]()
			expectedOtherData = "other expected data"
		)

		got := opt.OrElse(expectedOtherData)

		assert.Equal(t, expectedOtherData, got)
	})
	t.Run("if optional has value returns it", func(t *testing.T) {
		var (
			expectedValue     = "expectedValue"
			opt               = goption.Of(expectedValue)
			expectedOtherData = "other expected data"
		)

		got := opt.OrElse(expectedOtherData)

		assert.Equal(t, expectedValue, got)
	})
}
