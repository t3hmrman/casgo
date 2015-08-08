package cas_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/t3hmrman/casgo/cas"
)

var _ = Describe("CasGo", func() {
	Describe("CAS Server", func() {
		It("Should be creatable with nil configuration", func() {
			server, err := NewCASServer(nil)
			Expect(err).To(BeNil())
			Expect(server).ToNot(BeNil())
		})

		It("Should be creatable with non-nil configuration", func() {
			config, err := NewCASServerConfig("")
			config["companyName"] = "Casgo Testing Company"
			config["dbName"] = "casgo_test"
			config["templatesDirectory"] = "../templates"
			Expect(err).To(BeNil())

			server, err := NewCASServer(config)
			Expect(err).To(BeNil())
			Expect(server).ToNot(BeNil())
		})

		It("Should produce a predicatable address from the GetAddr function", func() {
			config, err := NewCASServerConfig("")
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
