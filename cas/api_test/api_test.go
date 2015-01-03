package api_test

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"log"
	"net/http"
)

var _ = Describe("CasGo API", func() {

	Describe("#authenticateAPIUser", func() {
		It("Should fail for unauthenticated users", func() {
			resp, err := http.Get(testHTTPServer.URL + "/api/services")
			Expect(err).To(BeNil())

			rawBody, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			var respJSON map[string]interface{}
			json.Unmarshal(rawBody, &respJSON)
			Expect(respJSON["status"]).To(Equal("error"))

			log.Println("Response: %v", respJSON)
		})
		// It("Should authenticate a user without a session who has passed an API key", func() {
		//	http.Get(testHTTPServer.URL + "/api/services")
		// })
	})

	// Describe("#HookupAPIEndpoints", func() {
	//	It("Should hookup an endpoint for listing services (GET /api/services)", func() {})
	//	It("Should hookup an endpoint for creating services (POST /api/services)", func() {})
	//	It("Should hookup an endpoint for updating services (PUT /api/services/{servicename})", func() {})
	//	It("Should hookup an endpoint for deleting services (DELETE /api/services/{servicename})", func() {})
	//	It("Should hookup an endpoint for retrieving services for a user (GET /api/sessions/{userEmail}/services)", func() {})
	//	It("Should hookup an endpoint for retrieving logged in users's session (GET /api/sessions)", func() {})
	// })

	// Describe("Session handler", func() {
	//	It("Should return an error if there is no session", func() {})
	// })

	// Describe("Current user's services listing endpoint", func() {
	//	It("Should return an error if there is no user logged in", func() {})
	// })

	// Describe("#getSessionAndUser", func() {
	//	It("SHould fail and return an error if there is no session", func() {})
	// })

	// Describe("#GetServices (GET /services)", func() {
	//	It("Should list all services for an admin user", func() {})
	//	It("Should display an error for non-admin users", func() {})
	// })

	// Describe("#CreateService (POST /services)", func() {
	//	It("Should create a service for an admin user", func() {})
	//	It("Should display an error for non-admin users", func() {})
	// })

	// Describe("#RemoveService (DELETE /services)", func() {
	//	It("Should create a service for an admin user", func() {})
	//	It("Should display an error for non-admin users", func() {})
	// })

	// Describe("#UpdateService (PUT /services)", func() {
	//	It("Should create a service for an admin user", func() {})
	//	It("Should display an error for non-admin users", func() {})
	// })

})
