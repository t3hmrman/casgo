package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/t3hmrman/casgo/cas"
	"net/http"
	"net/url"
	"strings"
)

var API_SERVICE_TEST_DATA map[string]string = map[string]string{
	"fakeServiceName":                 "test_service_created_by_test",
	"fakeServiceUrl":                  "localhost:9999/validateCASLogin",
	"fakeServiceAdminEmail":           "admin@test.com",
	"nameOfFixtureServiceToDelete":    "test_service_2",
	"nameOfFixtureServiceToUpdate":    "test_service_3",
	"updatedFixtureServiceUrl":        "localhost:3002/validateCASLogin",
	"updatedFixtureServiceAdminEmail": "updated@test.com",
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
	Expect(newService["name"]).To(Equal(API_SERVICE_TEST_DATA["fakeServiceName"]))
	Expect(newService["url"]).To(Equal(API_SERVICE_TEST_DATA["fakeServiceUrl"]))
	Expect(newService["adminEmail"]).To(Equal(API_SERVICE_TEST_DATA["fakeServiceAdminEmail"]))

	return newService
}

var _ = Describe("CasGo /api/services API", func() {
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
						"name":       {API_SERVICE_TEST_DATA["fakeServiceName"]},
						"url":        {API_SERVICE_TEST_DATA["fakeServiceUrl"]},
						"adminEmail": {API_SERVICE_TEST_DATA["fakeServiceAdminEmail"]},
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
				Name:       API_SERVICE_TEST_DATA["fakeServiceName"],
				Url:        API_SERVICE_TEST_DATA["fakeServiceUrl"],
				AdminEmail: API_SERVICE_TEST_DATA["fakeServiceAdminEmail"],
			})
		})
	})

	Describe("#RemoveService (DELETE /services)", func() {
		It("Should properly delete a service for an admin user", func() {
			// Craft request with admin user's API key
			req, err := http.NewRequest(
				"DELETE",
				testHTTPServer.URL+"/api/services/"+API_SERVICE_TEST_DATA["nameOfFixtureServiceToDelete"],
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
			Expect(respJSON["data"]).To(Equal(API_SERVICE_TEST_DATA["nameOfFixtureServiceToDelete"]))
		})

		It("Should display an error for non-admin users", func() {
			// Craft request with regular user's API key
			req, err := http.NewRequest(
				"DELETE",
				testHTTPServer.URL+"/api/services/"+API_SERVICE_TEST_DATA["fakeServiceName"],
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

	Describe("#UpdateService (PUT /api/services)", func() {
		It("Should update a service with valid input from an admin user", func() {
			// Craft request with admin user's API key
			req, err := http.NewRequest(
				"PUT",
				testHTTPServer.URL+"/api/services/"+API_SERVICE_TEST_DATA["nameOfFixtureServiceToUpdate"],
				strings.NewReader(
					url.Values{
						"url":        {API_SERVICE_TEST_DATA["updatedFixtureServiceUrl"]},
						"adminEmail": {API_SERVICE_TEST_DATA["updatedFixtureServiceAdminEmail"]},
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
		})

		It("Should fail to update a service with valid input from a non-admin user", func() {
			// Craft request with admin user's API key
			req, err := http.NewRequest(
				"PUT",
				testHTTPServer.URL+"/api/services/"+API_SERVICE_TEST_DATA["nameOfFixtureServiceToUpdate"],
				strings.NewReader(
					url.Values{
						"name":       {API_SERVICE_TEST_DATA["nameOfFixtureServiceToUpdate"]},
						"url":        {API_SERVICE_TEST_DATA["updatedFixtureServiceUrl"]},
						"adminEmail": {API_SERVICE_TEST_DATA["updatedFixtureServiceAdminEmail"]},
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
	})

})
