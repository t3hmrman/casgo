package cas

import (
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

// CAS Server creation should fail if no configuration is provided
func TestNilConfigCASServerCreation(t *testing.T) {
	server, err := NewCASServer(nil)
	if server != nil || err == nil {
		t.Error("CAS Server should have failed with nil configuration object")
	}
}

// Utility function for setting up CAS Server
func setupCASServer(t *testing.T) *CAS {
	config, err := NewCASServerConfig(map[string]string{
		"companyName":        "Casgo Testing Company",
		"dbName":             "casgo_test",
		"templatesDirectory": "../templates",
		"logLevel":           "INFO",
	})
	if err != nil {
		t.Error("Error creating config:", err)
	}

	server, err := NewCASServer(config)
	if server == nil || err != nil {
		t.Error("Server creation failed:", err)
	}

	return server
}

// CAS Server creation should succeed if default configuration is made
func TestDefaultConfigCASServerCreation(t *testing.T) { setupCASServer(t) }

// CAS Server init should properly attach handler functions to expected addresses
func TestCASGetAddrFn(t *testing.T) {
	server := setupCASServer(t)
	expectedAddress := CONFIG_DEFAULTS["host"] + ":" + CONFIG_DEFAULTS["port"]
	actualAddress := server.GetAddr()

	if actualAddress != expectedAddress {
		t.Error("Expected address [%s], got [%s]", expectedAddress, actualAddress)
	}
}

// Test retrieving service from DB
func TestGetServiceFn(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping DB-involved test (in short mode).")
	}

	//teardownDB(server, t)
}

// Test making new tickets in DB
func TestMakeNewTicketForServiceFn(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping DB-involved test (in short mode).")
	}

}

// Utility function for setting up necessary things for http test
func setupHTTPTest(t *testing.T) (*CAS, *httptest.Server) {
	// Setup CAS server and DB
	server := setupCASServer(t)
	setupDb(server, t)

	httpTestServer := httptest.NewServer(server.serveMux)
	return server, httpTestServer
}

// Utility test harness for HTTP tests
func httpTestHarness(t *testing.T, endpoint string, testFunc func(*testing.T, *CAS, *httptest.Server, *goquery.Document)) {
	if testing.Short() {
		t.Skip("Skipping integration test (in short mode).")
	}

	// Setup http test server
	server, httpTestServer := setupHTTPTest(t)
	defer httpTestServer.Close()
	defer teardownDb(server, t)

	// Visit specified endpoint
	doc, err := goquery.NewDocument(httpTestServer.URL + endpoint)
	if err != nil {
		t.Error(err)
	}

	// Run the function we were given
	testFunc(t, server, httpTestServer, doc)
}

// Test index page load
func TestHTTPIndexPageLoad(t *testing.T) {
	httpTestHarness(t, "", func(t *testing.T, _ *CAS, httpTestServer *httptest.Server, doc *goquery.Document) {
		// Ensure title of index page (endpoint "") contains what we expect
		expectedText, actualText := "CasGo", doc.Find("title").Text()
		if expectedText != actualText {
			t.Errorf("Actual title text [%s] != expected title text [%s]", actualText, expectedText)
		}
	})
}

// Test login page display (check for some expected elements)
func TestHTTPLoginPageLoad(t *testing.T) {
	httpTestHarness(t, "/login", func(t *testing.T, server *CAS, httpTestServer *httptest.Server, doc *goquery.Document) {
		// Ensure actual title text of login page (endpoint "/login") matches what is expected
		expectedText := server.Config["companyName"] + " CasGo Login"
		actualText := doc.Find("title").Text()
		if expectedText != actualText {
			t.Errorf("Actual title text [%s] != expected title text [%s]", actualText, expectedText)
		}

	})
}

// Test register page load
func TestHTTPRegisterPageLoad(t *testing.T) {
	httpTestHarness(t, "/register", func(t *testing.T, server *CAS, httpTestServer *httptest.Server, doc *goquery.Document) {
		expectedText := server.Config["companyName"] + " CasGo Register"
		actualText := doc.Find("title").Text()
		if expectedText != actualText {
			t.Errorf("Actual title text [%s] != expected title text [%s]", actualText, expectedText)
		}
	})
}
