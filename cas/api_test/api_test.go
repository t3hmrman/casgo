package api_test

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/t3hmrman/casgo/cas"
	"io/ioutil"
	"net/http"
)

var API_TEST_DATA map[string]string = map[string]string{
	"exampleAdminOnlyURI":          "/api/services",
	"exampleRegularUserURI":        "/api/sessions",
	"userApiKey":                   "userapikey",
	"userApiSecret":                "badsecret",
	"adminApiKey":                  "adminapikey",
	"adminApiSecret":               "badsecret",
}

// List of tuples that describe all endpoints hierarchically
var EXPECTED_API_ENDPOINTS map[string][]StringTuple = map[string][]StringTuple{
	"/api/services": []StringTuple{
		StringTuple{"GET", "/api/services"},
		StringTuple{"POST", "/api/services"},
		StringTuple{"GET", "/api/services"},
		StringTuple{"POST", "/api/services"},
		StringTuple{"PUT", "/api/services/{servicename}"},
		StringTuple{"DELETE", "/api/services/{servicename}"},
	},
	"/api/sessions": []StringTuple{
		StringTuple{"GET", "/api/sessions/{userEmail}/services"},
		StringTuple{"GET", "/api/sessions"},
	},
}

// Helper that checks to ensure unauthorized error response from performing an API request
func expectInsufficientPermissionsFromAPIRequest(req *http.Request) {
	// Perform request
	_, _, respJSON := jsonAPIRequestWithCustomHeaders(req)
	Expect(respJSON).NotTo(BeNil())
	Expect(respJSON["status"]).To(Equal("error"))
	Expect(respJSON["message"]).To(Equal(InsufficientPermissionsError.Msg))
}

// Function to be used with client creation to disallow redirects from API
func failRedirect(req *http.Request, via []*http.Request) error {
	Expect(req).To(BeNil())
	return errors.New("No redirects allowed")
}

// Utility function for performing JSON API requests
func jsonAPIRequestWithCustomHeaders(req *http.Request) (*http.Client, *http.Request, map[string]interface{}) {
	client := &http.Client{
		CheckRedirect: failRedirect,
	}

	// Perform request
	resp, err := client.Do(req)
	Expect(err).To(BeNil())

	// Read response body
	rawBody, err := ioutil.ReadAll(resp.Body)
	Expect(err).To(BeNil())

	// Parse response body into a map
	var respJSON map[string]interface{}
	err = json.Unmarshal(rawBody, &respJSON)
	Expect(err).To(BeNil())

	return client, req, respJSON
}

var _ = Describe("CasGo API", func() {

	Describe("#authenticateAPIUser", func() {
		It("Should fail for unauthenticated users", func() {
			// Craft a request request
			req, err := http.NewRequest("GET", testHTTPServer.URL+API_TEST_DATA["exampleRegularUserURI"], nil)
			Expect(err).To(BeNil())

			// Perform request
			_, _, respJSON := jsonAPIRequestWithCustomHeaders(req)
			Expect(respJSON["status"]).To(Equal("error"))
			Expect(respJSON["message"]).To(Equal(FailedToAuthenticateUserError.Msg))
		})

		It("Should properly authenticate a valid regular user's API key and secret to a non-admin-only endpoint", func() {
			// Craft request with regular user's API key
			req, err := http.NewRequest("GET", testHTTPServer.URL+API_TEST_DATA["exampleRegularUserURI"], nil)
			Expect(err).To(BeNil())
			req.Header.Add("X-Api-Key", API_TEST_DATA["userApiKey"])
			req.Header.Add("X-Api-Secret", API_TEST_DATA["userApiSecret"])

			// Perform request
			_, _, respJSON := jsonAPIRequestWithCustomHeaders(req)
			Expect(respJSON["status"]).To(Equal("success"))
		})

		It("Should properly authenticate a valid admin user's API key and secret to an admin-only endpoint", func() {
			// Craft request with admin user's API key
			req, err := http.NewRequest("GET", testHTTPServer.URL+API_TEST_DATA["exampleAdminOnlyURI"], nil)
			Expect(err).To(BeNil())
			req.Header.Add("X-Api-Key", API_TEST_DATA["adminApiKey"])
			req.Header.Add("X-Api-Secret", API_TEST_DATA["adminApiSecret"])

			// Perform JSON API request
			_, _, respJSON := jsonAPIRequestWithCustomHeaders(req)
			Expect(respJSON["status"]).To(Equal("success"))
		})

		It("Should fail to authenticate a regular user to an admin-only endpoint", func() {
			// Craft request with regular user's API key
			req, err := http.NewRequest("GET", testHTTPServer.URL+API_TEST_DATA["exampleAdminOnlyURI"], nil)
			Expect(err).To(BeNil())
			req.Header.Add("X-Api-Key", API_TEST_DATA["userApiKey"])
			req.Header.Add("X-Api-Secret", API_TEST_DATA["userApiSecret"])

			// Perform JSON API request
			expectInsufficientPermissionsFromAPIRequest(req)
		})

	})

	Describe("#HookupAPIEndpoints", func() {
		It("Should hookup all /api/services endpoints", func() {
			testMux := mux.NewRouter()
			api, err := NewCasgoFrontendAPI(nil)
			api.HookupAPIEndpoints(testMux)
			Expect(err).To(BeNil())

			// Check all expected endpoints below "/api/services"
			var routeMatch mux.RouteMatch
			for _, tuple := range EXPECTED_API_ENDPOINTS["/api/services"] {
				// Craft request for listing services (GET /api/services)
				req, err := http.NewRequest(tuple.First(), tuple.Second(), nil)
				Expect(err).To(BeNil())

				// Get pattern that was matched
				Expect(testMux.Match(req, &routeMatch)).To(BeTrue())
			}

		})
	})
	
})
