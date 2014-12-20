package cas_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/core"
	. "github.com/sclevine/agouti/dsl"
	. "github.com/sclevine/agouti/matchers"
)

var INTEGRATION_TEST_DATA map[string]string = map[string]string{
	"newUserEmail":        "testuser@testemail.com",
	"newUserPassword":     "testpassword",
	"fixtureUserEmail":    "test@test.com",
	"fixtureUserPassword": "test",
}

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
		expectedTitle := testCASServer.Config["companyName"] + " - Login"
		Expect(page.Find("#page-title")).To(HaveText(expectedTitle))
	})

	Scenario("Find the expected title on the register page", func() {
		page.Navigate(testHTTPServer.URL + "/register")
		expectedTitle := testCASServer.Config["companyName"] + " - Register"
		Expect(page.Find("#page-title")).To(HaveText(expectedTitle))
		Expect(page.Find("#email")).To(BeFound())
		Expect(page.Find("#password")).To(BeFound())
	})

	Scenario("Successfully register a new user", func() {
		StepRegisterUser(INTEGRATION_TEST_DATA["newUserEmail"], INTEGRATION_TEST_DATA["newUserPassword"], page)
	})

	Scenario("Login with a user created by the users.json fixture", func() {
		StepLoginUser(INTEGRATION_TEST_DATA["fixtureUserEmail"], INTEGRATION_TEST_DATA["fixtureUserPassword"], page)
	})

	Scenario("Login and log out a user created by the users.json fixture", func() {
		StepLoginUser(INTEGRATION_TEST_DATA["fixtureUserEmail"], INTEGRATION_TEST_DATA["fixtureUserPassword"], page)
		StepLogoutUser(page)
	})

})

// Reusable testing steps
var StepRegisterUser func(string, string, Page) = func(email, password string, page Page) {
	Step("Navigate to the register page", func() {
		page.Navigate(testHTTPServer.URL + "/register")
		Expect(page.Find("#email")).To(BeFound())
		Expect(page.Find("#password")).To(BeFound())
	})

	Step("Fill out and submit the new user registration form", func() {
		Fill(page.Find("#email"), email)
		Fill(page.Find("#password"), password)
		Submit(page.Find("#frmRegister"))
	})

	Step("See alert telling you that you've successfully registered", func() {
		Expect(page.Find("div.alert.success")).To(BeFound())
		Expect(page.Find("div.alert.success")).To(HaveText("Registration successful!"))
	})
}

var StepLoginUser func(string, string, Page) = func(email, password string, page Page) {
	Step("Navigate to the login page", func() {
		page.Navigate(testHTTPServer.URL + "/login")
		Expect(page.Find("#email")).To(BeFound())
		Expect(page.Find("#password")).To(BeFound())
	})

	Step("Fill out and submit the user login form", func() {
		Fill(page.Find("#email"), email)
		Fill(page.Find("#password"), password)
		Submit(page.Find("#frmLogin"))
	})

	Step("See alert telling you that you've successfully registered", func() {
		Expect(page.Find("div.alert.success")).To(BeFound())
		Expect(page.Find("div.alert.success")).To(HaveText("Successful log in! Redirecting to services page..."))
	})

}

var StepLogoutUser func(Page) = func(page Page) {
		page.Navigate(testHTTPServer.URL + "/logout")
		Expect(page.Find("div.alert.success")).To(BeFound())
		Expect(page.Find("div.alert.success")).To(HaveText("Successfully logged out"))
}
