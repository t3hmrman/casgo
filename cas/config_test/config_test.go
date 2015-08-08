package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/t3hmrman/casgo/cas"
	"log"
	"os"
	"path/filepath"
)

var _ = Describe("CasGo Config", func() {
	Describe("Config creation", func() {
		It("Should work (pick sensible defaults) with nil passed in", func() {
			config, err := NewCASServerConfig("")
			Expect(err).To(BeNil())
			Expect(config).ToNot(BeNil())
		})
	})

	Describe("Config defaults", func() {
		It("should match defaults as specified in the package", func() {
			config, err := NewCASServerConfig("")
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

	Describe("Config ENV overriding", func() {
		It("Should properly override host", func() {
			// Save current ENV value
			previousEnvValue := os.Getenv("CASGO_HOST")
			host := "TESTHOST"
			err := os.Setenv("CASGO_HOST", host)
			Expect(err).To(BeNil())

			config, err := NewCASServerConfig("")
			Expect(err).To(BeNil())
			Expect(config).ToNot(BeNil())
			Expect(config["host"]).To(Equal(host))

			// Reset ENV
			os.Setenv("CASGO_HOST", previousEnvValue)
		})
	})

	Describe("Config from file", func() {
		It("Should fail if file is specified but missing", func() {
			config, err := NewCASServerConfig("missing.json")
			Expect(err).ToNot(BeNil())
			Expect(config).To(BeNil())
		})

		It("Should fail if file is present but has invalid JSON", func() {
			config, err := NewCASServerConfig("../../fixtures/invalid-json-config.json")
			Expect(err).ToNot(BeNil())
			Expect(config).To(BeNil())
		})

		It("Should succeed on proper config, and update values", func() {
			config, err := NewCASServerConfig("../../fixtures/valid-config.json")
			Expect(err).To(BeNil())
			Expect(config).ToNot(BeNil())
			log.Printf("%v", config)
			Expect(config["dbName"]).To(Equal("TEST_DB_NAME"))
		})

	})

})
