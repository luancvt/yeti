package yang_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestYang(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Yang Suite")
}
