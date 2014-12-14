package cas_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/t3hmrman/casgo/cas"
)

var _ = Describe("Cas", func() {
	Describe("CAS Server", func() {
		It("Should be creatable with nil configuration", func() {
			server, err := NewCASServer(nil)
			Expect(err).To(BeNil())
			Expect(server).ToNot(BeNil())
		})

		It("Should be creatable with non-nil configuration", func() {
			config, err := NewCASServerConfig(map[string]string{
				"companyName":        "Casgo Testing Company",
				"dbName":             "casgo_test",
				"templatesDirectory": "../templates",
			})
			Expect(err).To(BeNil())

			server, err := NewCASServer(config)
			Expect(err).To(BeNil())
			Expect(server).ToNot(BeNil())
		})

		It("Should produce a predicatable address from the GetAddr function", func() {
			config, err := NewCASServerConfig(nil)
			Expect(err).To(BeNil())

			server, err := NewCASServer(config)
			Expect(err).To(BeNil())
			Expect(server).ToNot(BeNil())

			expectedAddress := CONFIG_DEFAULTS["host"] + ":" + CONFIG_DEFAULTS["port"]
			actualAddress := server.GetAddr()
			Expect(actualAddress).To(Equal(expectedAddress))
		})
	})
})
