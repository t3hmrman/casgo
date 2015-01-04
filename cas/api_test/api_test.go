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
	"net/url"
	"strings"
)

var API_TEST_DATA map[string]string = map[string]string{
	"exampleAdminOnlyURI":          "/api/services",
	"exampleRegularUserURI":        "/api/sessions",
	"userApiKey":                   "userapikey",
	"userApiSecret":                "badsecret",
	"adminApiKey":                  "adminapikey",
	"adminApiSecret":               "badsecret",
	"fakeServiceName":              "test_service_3",
	"fakeServiceUrl":               "localhost:3002/validateCASLogin",
	"fakeServiceAdminEmail":        "admin@test.com",
	"nameOfFixtureServiceToDelete": "test_service_2",
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

// Helper function to create new fake service
func createNewFakeService(service CASService) map[string]interface{} {
	// Craft request with admin user's API key
	req, err := http.NewRequest(
		"POST",
		testHTTPServer.URL+"/api/services",
		strings.NewReader(
			url.Values{
				"name":       {service.Name},
				"url":        {service.Url},
				"adminEmail": {service.AdminEmail},
			}.Encode()),
	)
	Expect(err).To(BeNil())

	// Set header for request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Api-Key", API_TEST_DATA["adminApiKey"])
	req.Header.Add("X-Api-Secret", API_TEST_DATA["adminApiSecret"])

	// Perform request
	_, _, respJSON := jsonAPIRequestWithCustomHeaders(req)
	Expect(respJSON).NotTo(BeNil())
	Expect(respJSON["status"]).To(Equal("success"))

	// Get the list of services that was returned (dependent on fixture)
	var newService map[string]interface{}
	Expect(respJSON["data"]).To(BeAssignableToTypeOf(newService))
	newService = respJSON["data"].(map[string]interface{})
	Expect(newService["name"]).To(Equal(API_TEST_DATA["fakeServiceName"]))
	Expect(newService["url"]).To(Equal(API_TEST_DATA["fakeServiceUrl"]))
	Expect(newService["adminEmail"]).To(Equal(API_TEST_DATA["fakeServiceAdminEmail"]))

	return newService
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

	Describe("#GetServices (GET /services)", func() {
		It("Should list all services for an admin user", func() {
			// Craft request with admin user's API key
			req, err := http.NewRequest("GET", testHTTPServer.URL+"/api/services", nil)
			Expect(err).To(BeNil())
			req.Header.Add("X-Api-Key", API_TEST_DATA["adminApiKey"])
			req.Header.Add("X-Api-Secret", API_TEST_DATA["adminApiSecret"])

			_, _, respJSON := jsonAPIRequestWithCustomHeaders(req)
			Expect(respJSON["status"]).To(Equal("success"))

			// Get the list of services that was returned
			var rawServicesList []interface{}
			Expect(respJSON["data"]).To(BeAssignableToTypeOf(rawServicesList))
			rawServicesList = respJSON["data"].([]interface{})
			Expect(len(rawServicesList)).To(BeNumerically(">", 0))

			// Check the list of services for map that represents a service we expect to exist
			foundExpectedService := false
			var serviceMap map[string]interface{}
			for _, rawMap := range rawServicesList {
				Expect(rawMap).To(BeAssignableToTypeOf(serviceMap))
				serviceMap = rawMap.(map[string]interface{})
				if name, ok := serviceMap["name"]; ok && name == "test_service" {
					foundExpectedService = true
					Expect(serviceMap["name"]).To(Equal("test_service"))
					Expect(serviceMap["url"]).To(Equal("localhost:3000/validateCASLogin"))
					Expect(serviceMap["adminEmail"]).To(Equal("admin@test.com"))
				}
			}
			Expect(foundExpectedService).To(BeTrue())
		})
	})

	Describe("#CreateService (POST /services)", func() {
		It("Should fail to create a service given invalid input from an admin user", func() {
			// Craft request with admin user's API key
			req, err := http.NewRequest(
				"POST",
				testHTTPServer.URL+"/api/services",
				strings.NewReader(url.Values{"nope": {"nope"}}.Encode()),
			)
			Expect(err).To(BeNil())

			// Set header
			req.Header.Add("X-Api-Key", API_TEST_DATA["adminApiKey"])
			req.Header.Add("X-Api-Secret", API_TEST_DATA["adminApiSecret"])

			// Perform request
			_, _, respJSON := jsonAPIRequestWithCustomHeaders(req)
			Expect(respJSON).NotTo(BeNil())
			Expect(respJSON["status"]).To(Equal("error"))
			Expect(respJSON["message"]).To(Equal(InvalidServiceError.Msg))
		})

		It("Should display an error for non-admin users", func() {
			// Craft request with regular user's API key
			req, err := http.NewRequest(
				"POST",
				testHTTPServer.URL+"/api/services",
				strings.NewReader(
					url.Values{
						"name":       {API_TEST_DATA["fakeServiceName"]},
						"url":        {API_TEST_DATA["fakeServiceUrl"]},
						"adminEmail": {API_TEST_DATA["fakeServiceAdminEmail"]},
					}.Encode()),
			)
			Expect(err).To(BeNil())

			// Set header for request
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Add("X-Api-Key", API_TEST_DATA["userApiKey"])
			req.Header.Add("X-Api-Secret", API_TEST_DATA["userApiSecret"])

			// Perform request
			expectInsufficientPermissionsFromAPIRequest(req)
		})

		It("Should successfully create a service given valid input from an admin user", func() {
			createNewFakeService(CASService{
				Name:       API_TEST_DATA["fakeServiceName"],
				Url:        API_TEST_DATA["fakeServiceUrl"],
				AdminEmail: API_TEST_DATA["fakeServiceAdminEmail"],
			})
		})
	})

	Describe("#RemoveService (DELETE /services)", func() {
		It("Should properly delete a service for an admin user", func() {
			// Craft request with admin user's API key
			req, err := http.NewRequest(
				"DELETE",
				testHTTPServer.URL+"/api/services/"+API_TEST_DATA["nameOfFixtureServiceToDelete"],
				nil,
			)
			Expect(err).To(BeNil())

			// Set header for request
			req.Header.Add("X-Api-Key", API_TEST_DATA["adminApiKey"])
			req.Header.Add("X-Api-Secret", API_TEST_DATA["adminApiSecret"])

			// Perform request
			_, _, respJSON := jsonAPIRequestWithCustomHeaders(req)
			Expect(respJSON).NotTo(BeNil())
			Expect(respJSON["status"]).To(Equal("success"))
			Expect(respJSON["data"]).To(Equal(API_TEST_DATA["nameOfFixtureServiceToDelete"]))
		})

		It("Should display an error for non-admin users", func() {
			// Craft request with regular user's API key
			req, err := http.NewRequest(
				"DELETE",
				testHTTPServer.URL+"/api/services/"+API_TEST_DATA["fakeServiceName"],
				nil,
			)
			Expect(err).To(BeNil())

			// Set header for request
			req.Header.Add("X-Api-Key", API_TEST_DATA["userApiKey"])
			req.Header.Add("X-Api-Secret", API_TEST_DATA["userApiSecret"])

			// Perform request
			expectInsufficientPermissionsFromAPIRequest(req)
		})
	})

	// Describe("#UpdateService (PUT /services)", func() {
	//	It("Should create a service for an admin user", func() {})
	//	It("Should display an error for non-admin users", func() {})
	// })

})
