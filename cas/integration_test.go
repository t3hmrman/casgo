package cas_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/core"
	. "github.com/sclevine/agouti/dsl"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Feature("CASGO", func() {
	var page Page

	Background(func() {
		page = CustomPage(Use().With("handlesAlerts"))
		page.Navigate(testHTTPServer.URL)
		page.Size(640, 480)
	})

	AfterEach(func() {
		page.Destroy()
	})

	Scenario("Finding the expected title on the index page", func() {
		Expect(page).To(HaveTitle("CasGo"))
	})

	Scenario("Find the expected title on the login page", func() {
		page.Navigate(testHTTPServer.URL + "/login")
		expectedTitle := testCASServer.Config["companyName"] + " CasGo Login"
		Expect(page).To(HaveTitle(expectedTitle))
	})

	Scenario("Find the expected title on the register page", func() {
		page.Navigate(testHTTPServer.URL + "/register")
		expectedTitle := testCASServer.Config["companyName"] + " CasGo Register"
		Expect(page).To(HaveTitle(expectedTitle))
	})

})