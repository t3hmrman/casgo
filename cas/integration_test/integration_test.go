package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/matchers"
	 "github.com/sclevine/agouti"
)

var INTEGRATION_TEST_DATA map[string]string = map[string]string{
	"newUserEmail":         "testuser@testemail.com",
	"newUserPassword":      "testpassword",
	"fixtureUserEmail":     "test@test.com",
	"fixtureUserPassword":  "test",
	"fixtureAdminEmail":    "admin@test.com",
	"fixtureAdminPassword": "test",
}

var _ = Describe("CasGo", func() {
	var page *agouti.Page

	BeforeEach(func() {
		page, err := agoutiDriver.NewPage()
		Expect(err).NotTo(HaveOccurred())
		page.Navigate(testHTTPServer.URL)
		page.Size(640, 480)
	})

	AfterEach(func() {
		page.Destroy()
	})

	It("Finding the expected title on the index page", func() {
		Expect(page).To(HaveTitle("CasGo"))
	})

	It("Find the expected title on the login page", func() {
		page.Navigate(testHTTPServer.URL + "/login")
		expectedTitle := testCASServer.Config["companyName"] + " - Login"
		Expect(page.Find("#page-title")).To(HaveText(expectedTitle))
	})

	It("Find the expected title on the register page", func() {
		page.Navigate(testHTTPServer.URL + "/register")
		expectedTitle := testCASServer.Config["companyName"] + " - Register"
		Expect(page.Find("#page-title")).To(HaveText(expectedTitle))
		Expect(page.Find("#email")).To(BeFound())
		Expect(page.Find("#password")).To(BeFound())
	})

	It("Successfully register a new user", func() {
		StepRegisterUser(INTEGRATION_TEST_DATA["newUserEmail"], INTEGRATION_TEST_DATA["newUserPassword"], page)
	})

	It("Login with a user created by the users.json fixture", func() {
		StepLoginUser(INTEGRATION_TEST_DATA["fixtureUserEmail"], INTEGRATION_TEST_DATA["fixtureUserPassword"], page)
		StepLogoutUser(page)
	})

	It("Login and log out a user created by the users.json fixture", func() {
		StepLoginUser(INTEGRATION_TEST_DATA["fixtureUserEmail"], INTEGRATION_TEST_DATA["fixtureUserPassword"], page)
		StepLogoutUser(page)
	})

	It("The casgo SPA contains a reduced set of navigation options if logged in as a regular user(in users.json fixture)", func() {
		StepLoginUser(INTEGRATION_TEST_DATA["fixtureUserEmail"], INTEGRATION_TEST_DATA["fixtureUserPassword"], page)
		By("Ensure the casgo SPA shows more navigation options to the admin user", func() {
			page.Navigate(testHTTPServer.URL + "/")
			Expect(page.Find("#topnav-services-link")).To(BeFound())
			Expect(page.Find("#topnav-manage-link")).ToNot(BeFound())
			Expect(page.Find("#topnav-statistics-link")).ToNot(BeFound())
		})
		StepLogoutUser(page)
	})

	It("The casgo SPA contains extra navigation options if logged in as an admin user (in users.json fixture)", func() {
		StepLoginUser(INTEGRATION_TEST_DATA["fixtureAdminEmail"], INTEGRATION_TEST_DATA["fixtureAdminPassword"], page)
		By("Ensure the casgo SPA shows more navigation options to the admin user", func() {
			page.Navigate(testHTTPServer.URL + "/")
			Expect(page.Find("#topnav-services-link")).To(BeFound())
			Expect(page.Find("#topnav-manage-link")).To(BeFound())
		})
		StepLogoutUser(page)
	})
})

/** Reusable testing steps **/

// Steps to register a user
var StepRegisterUser func(string, string, *agouti.Page) = func(email, password string, page *agouti.Page) {
	By("Navigate to the register page", func() {
		page.Navigate(testHTTPServer.URL + "/register")
		Expect(page.Find("#email")).To(BeFound())
		Expect(page.Find("#password")).To(BeFound())
	})

	By("Fill out and submit the new user registration form", func() {
		page.Find("#email").Fill(email)
		page.Find("#password").Fill(password)
		page.Find("#frmRegister").Submit()
	})

	By("See alert telling you that you've successfully registered", func() {
		Expect(page.Find("div.alert.success")).To(BeFound())
		Expect(page.Find("div.alert.success")).To(HaveText("Registration successful!"))
	})
}

// Steps to simulate login
var StepLoginUser func(string, string, *agouti.Page) = func(email, password string, page *agouti.Page) {
	By("Navigate to the login page", func() {
		page.Navigate(testHTTPServer.URL + "/login")
		Expect(page.Find("#email")).To(BeFound())
		Expect(page.Find("#password")).To(BeFound())
	})

	By("Fill out and submit the user login form", func() {
		page.Find("#email").Fill(email)
		page.Find("#password").Fill(password)
		page.Find("#frmLogin").Submit()
	})

	By("See alert telling you that you've successfully registered", func() {
		Expect(page.Find("div.alert.success")).To(BeFound())
		Expect(page.Find("div.alert.success")).To(HaveText("Successful log in! Redirecting to services page..."))
	})

}

// Steps to simulate logout
var StepLogoutUser func(*agouti.Page) = func(page *agouti.Page) {
	page.Navigate(testHTTPServer.URL + "/logout")
	Expect(page.Find("div.alert.success")).To(BeFound())
	Expect(page.Find("div.alert.success")).To(HaveText("Successfully logged out"))
}
