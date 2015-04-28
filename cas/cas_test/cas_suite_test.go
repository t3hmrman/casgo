package cas_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestCasgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Casgo Suite")
}
