package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/t3hmrman/casgo/cas"
)

var _ = Describe("CasGo API", func() {
	Describe("CAS Server", func() {
		It("Should be creatable with nil configuration", func() {
			server, err := NewCASServer(nil)
			Expect(err).To(BeNil())
			Expect(server).ToNot(BeNil())
		})
	})
})
