package goption_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGoption(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Goption Suite")
}
