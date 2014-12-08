package cas_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCas(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cas Suite")
}
