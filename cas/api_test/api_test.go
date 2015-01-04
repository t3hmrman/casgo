package api_test

import (
	"encoding/json"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/t3hmrman/casgo/cas"
	"io/ioutil"
	"net/http"
)

var API_TEST_DATA map[string]string = map[string]string{
	"exampleAdminOnlyURI":   "/api/services",
	"exampleRegularUserURI": "/api/sessions",
	"userApiKey":            "userapikey",
	"userApiSecret":         "badsecret",
	"adminApiKey":           "adminapikey",
	"adminApiSecret":        "badsecret",
}

func failRedirect(req *http.Request, via []*http.Request) error {
	Expect(req).To(BeNil())
	return errors.New("No redirects allowed")
}

var _ = Describe("CasGo API", func() {

	Describe("#authenticateAPIUser", func() {
		It("Should fail for unauthenticated users", func() {
			resp, err := http.Get(testHTTPServer.URL + API_TEST_DATA["exampleRegularUserURI"])
			Expect(err).To(BeNil())

			// Read response body
			rawBody, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			// Parse response body into a map
			var respJSON map[string]interface{}
			json.Unmarshal(rawBody, &respJSON)
			Expect(respJSON["status"]).To(Equal("error"))
			Expect(respJSON["message"]).To(Equal(FailedToAuthenticateUserError.Msg))
		})

		It("Should properly authenticate a valid regular user's API key and secret to a non-admin-only endpoint", func() {
			client := &http.Client{
				CheckRedirect: failRedirect,
			}

			// Craft a request with api key and secret
			req, err := http.NewRequest("GET", testHTTPServer.URL+API_TEST_DATA["exampleRegularUserURI"], nil)
			req.Header.Add("X-Api-Key", API_TEST_DATA["userApiKey"])
			req.Header.Add("X-Api-Secret", API_TEST_DATA["userApiSecret"])

			// Perform request
			resp, err := client.Do(req)
			Expect(err).To(BeNil())

			// Read response body
			rawBody, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			// Parse response body into a map
			var respJSON map[string]interface{}
			json.Unmarshal(rawBody, &respJSON)
			Expect(respJSON["status"]).To(Equal("success"))
		})

		It("Should properly authenticate a valid admin user's API key and secret to an admin-only endpoint", func() {
			client := &http.Client{
				CheckRedirect: failRedirect,
			}

			// Craft a request with api key and secret
			req, err := http.NewRequest("GET", testHTTPServer.URL+API_TEST_DATA["exampleAdminOnlyURI"], nil)
			req.Header.Add("X-Api-Key", API_TEST_DATA["adminApiKey"])
			req.Header.Add("X-Api-Secret", API_TEST_DATA["adminApiSecret"])

			// Perform request
			resp, err := client.Do(req)
			Expect(err).To(BeNil())

			// Read response body
			rawBody, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			// Parse response body into a map
			var respJSON map[string]interface{}
			json.Unmarshal(rawBody, &respJSON)
			Expect(respJSON["status"]).To(Equal("success"))
		})

		It("Should fail to authenticate a regular user to an admin-only endpoint", func() {
			client := &http.Client{
				CheckRedirect: failRedirect,
			}

			// Craft a request with api key and secret
			req, err := http.NewRequest("GET", testHTTPServer.URL+API_TEST_DATA["exampleAdminOnlyURI"], nil)
			req.Header.Add("X-Api-Key", API_TEST_DATA["userApiKey"])
			req.Header.Add("X-Api-Secret", API_TEST_DATA["userApiSecret"])

			// Perform request
			resp, err := client.Do(req)
			Expect(err).To(BeNil())

			// Read response body
			rawBody, err := ioutil.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			// Parse response body into a map
			var respJSON map[string]interface{}
			json.Unmarshal(rawBody, &respJSON)
			Expect(respJSON["status"]).To(Equal("error"))
		})

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
