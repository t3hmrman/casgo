package cas

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
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

// Utility function for tearing down and setting up the database
func tearDownAndSetupDB(server *CAS, t *testing.T) {
	teardownErr := server.dbAdapter.TeardownDB()
	if teardownErr != nil {
		t.Errorf("Failed to tear down database: %s", *teardownErr.err)
	}

	setupErr := server.dbAdapter.SetupDB()
	if setupErr != nil {
		t.Errorf("Failed to set down database: %s", *setupErr.err)
	}
}

// Test setup and tear down of database (with utility function)
func TestDBSetupAndTeardown(t *testing.T) {
	server := setupCASServer(t)
	tearDownAndSetupDB(server, t)
}

// Test retrieving service from DB
func TestGetServiceFn(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping DB-involved test (in short mode).")
	}

	server := setupCASServer(t)
	tearDownAndSetupDB(server, t)
	server.dbAdapter.ImportTableDataFromFile("services", "id", "fixtures/services.json")

}

// Test making new tickets in DB
func TestMakeNewTicketForServiceFn(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping DB-involved test (in short mode).")
	}

}

// Login page tests
func TestLoginPage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test (in short mode).")
	}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	})

	loginServer := httptest.NewServer(testHandler)
	defer loginServer.Close()

	res, err := http.Get(loginServer.URL)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	if doc == nil {
		t.Error("Doc is nil!")
	}
}
