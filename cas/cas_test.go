package cas_test

import (
	. "github.com/t3hmrman/casgo/cas"

	"github.com/PuerkitoBio/goquery"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

	Describe("CAS Website", func() {
		It("Should have an working index page", func() {
			// Visit index endpoint
			doc, err := goquery.NewDocument(testHTTPServer.URL)
			Expect(err).To(BeNil())

			// Ensure title of index page (endpoint "") contains what we expect
			expectedText, actualText := "CasGo", doc.Find("title").Text()
			Expect(actualText).To(Equal(expectedText))
		})

		It("Should have a working login page", func() {

			// Visit login endpoint
			doc, err := goquery.NewDocument(testHTTPServer.URL + "/login")
			Expect(err).To(BeNil())

			expectedText := testCASServer.Config["companyName"] + " CasGo Login"
			actualText := doc.Find("title").Text()
			Expect(actualText).To(Equal(expectedText))
		})

	})
})
