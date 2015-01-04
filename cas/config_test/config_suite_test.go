package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestCasgoConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CasGo Config Suite")
}
