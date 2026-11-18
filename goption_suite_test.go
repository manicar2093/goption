package goption_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func toPointer[T any](value T) *T {
	var pointer T = value
	return &pointer
}

func TestGoption(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Goption Suite")
}
