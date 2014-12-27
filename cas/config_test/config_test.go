package cas_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/t3hmrman/casgo/cas"
	"os"
	"path/filepath"
)

var _ = Describe("CAS Config", func() {
	Describe("Config creation", func() {
		It("Should work (pick sensible defaults) with nil passed in", func() {
			config, err := NewCASServerConfig(nil)
			Expect(err).To(BeNil())
			Expect(config).ToNot(BeNil())
		})
	})

	Describe("Config defaults", func() {
		It("should match defaults as specified in the package", func() {
			config, err := NewCASServerConfig(nil)
			Expect(err).To(BeNil())
			Expect(config).ToNot(BeNil())

			// Get the absolute path of the config default directory
			expectedTemplatesDirectory, err := filepath.Abs(CONFIG_DEFAULTS["templatesDirectory"])
			Expect(err).To(BeNil())
			CONFIG_DEFAULTS["templatesDirectory"] = expectedTemplatesDirectory

			// Ensure all configurations are at defaults (or expected modified defaults)
			for k, v := range config {
				Expect(v).To(Equal(CONFIG_DEFAULTS[k]))
			}
		})
	})

	Describe("Config user overriding", func() {
		It("Should properly override host", func() {
			host := "fake-host-string"
			config, err := NewCASServerConfig(map[string]string{"host": host})
			Expect(err).To(BeNil())
			Expect(config["host"]).To(Equal(host))
		})
	})

	Describe("Config ENV overriding", func() {
		It("Should properly override host", func() {
			// Save current ENV value
			previousEnvValue := os.Getenv("CASGO_HOST")
			host := "TESTHOST"
			err := os.Setenv("CASGO_HOST", host)
			Expect(err).To(BeNil())

			config, err := NewCASServerConfig(map[string]string{"host":"SHOULD NOT BE THIS"})
			Expect(err).To(BeNil())
			Expect(config).ToNot(BeNil())
			Expect(config["host"]).To(Equal(host))
			
			// Reset ENV
			os.Setenv("CASGO_HOST", previousEnvValue)
		})
	})
})
