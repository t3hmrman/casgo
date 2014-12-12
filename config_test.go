package cas_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/t3hmrman/casgo/cas"
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

})
