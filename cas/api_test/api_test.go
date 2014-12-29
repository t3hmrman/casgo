package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/t3hmrman/casgo/cas"
)

var _ = Describe("CasGo API", func() {

	Describe("#HookupAPIEndpoints", func() {
		It("Should hookup an endpoint for listing services (GET /api/services)", func() {})
		It("Should hookup an endpoint for creating services (POST /api/services)", func() {})
		It("Should hookup an endpoint for updating services (PUT /api/services/{servicename})", func() {})
		It("Should hookup an endpoint for deleting services (DELETE /api/services/{servicename})", func() {})
		It("Should hookup an endpoint for retrieving services for a user (GET /api/sessions/{userEmail}/services)", func() {})
		It("Should hookup an endpoint for retrieving logged in users's session (GET /api/sessions)", func() {})
	})

	Describe("Session handler", func() {
		It("Should return an error if there is no session", func() {})
	})

	Describe("Current user's services listing endpoint", func() {
		It("Should return an error if there is no user logged in", func() {})
	})

	Describe("#getSessionAndUser", func() {
		It("SHould fail and return an error if there is no session", func() {})
	})

	Describe("#GetServices (GET /services)", func() {
		It("Should list all services for an admin user", func() {})
		It("Should display an error for non-admin users", func() {})
	})

	Describe("#CreateService (POST /services)", func() {
		It("Should create a service for an admin user", func() {})
		It("Should display an error for non-admin users", func() {})
	})

	Describe("#RemoveService (DELETE /services)", func() {
		It("Should create a service for an admin user", func() {})
		It("Should display an error for non-admin users", func() {})
	})

	Describe("#UpdateService (PUT /services)", func() {
		It("Should create a service for an admin user", func() {})
		It("Should display an error for non-admin users", func() {})
	})

})
