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
	config, err := NewCASServerConfig(map[string]string{"dbName": "casgo_test"})
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

// Test index page load
func TestHTTPIndexPageLoad(t *testing.T) {
	// Setup http test server
	server, httpTestServer := setupHTTPTest(t)
	defer httpTestServer.Close()

	doc, err := goquery.NewDocument(httpTestServer.URL)
	if err != nil {
		t.Error(err)
	}

	t.Log("docText:", doc.Text())

	// Ensure title of index page contains what we expect
	expectedText, actualText := "CasGo", doc.Find("title").Text()
	if expectedText != actualText {
		t.Errorf("Actual title text [%s] != expected title text [%s]", actualText, expectedText)
	}

	teardownDb(server, t)
}

// Test login page display (check for some expected elements)
func TestHTTPLoginPageLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test (in short mode).")
	}

	// Setup http test server
	server, httpTestServer := setupHTTPTest(t)
	defer httpTestServer.Close()

	doc, err := goquery.NewDocument(httpTestServer.URL + "/login")
	if err != nil {
		t.Error(err)
	}

	// Ensure actual title text matches what is expected
	expectedText := server.Config["companyName"] + " - Login"
	actualText := doc.Find("title").Text()
	if expectedText != actualText {
		t.Errorf("Actual title text [%s] != expected title text [%s]", actualText, expectedText)
	}

	teardownDb(server, t)
}
