package cas

import (
	"testing"
)

// Utility functions for tearing down/setting up the database
func setupDb(server *CAS, t *testing.T) {
	setupErr := server.dbAdapter.Setup()
	if setupErr != nil {
		if setupErr != nil { t.Errorf("Failed to set up database: %s", *setupErr.err) }
	}
}

func teardownDb(server *CAS, t *testing.T) {
	teardownErr := server.dbAdapter.Teardown()
	if teardownErr != nil {
		if teardownErr.err != nil { t.Errorf("Failed to tear down database: %s", *teardownErr.err) }
	}
}

// Test setup and tear down of database (with utility function)
func TestDBSetup(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping DB-involved test (in short mode).")
	}

	casServer := setupCASServer(t)
	setupDb(casServer, t)
	teardownDb(casServer, t)
}

// Test database import functionality (implemented LoadJSONFixture function should not fail to load fixtures)
func TestLoadJSONFixture(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping DB-involved test (in short mode).")
	}

	s := setupCASServer(t)
	setupDb(s, t)

	// Setup the services table for importing into
	err := s.dbAdapter.SetupServicesTable()
	if err != nil {
		t.Errorf("Failed to setup services table:", err)
	}

	// Import into the services table
	importErr := s.dbAdapter.LoadJSONFixture(s.dbAdapter.getDbName(), s.dbAdapter.getServicesTableName(), "fixtures/services.json")
	if importErr != nil {
		if importErr.err != nil {
			t.Log("DB error: %s", importErr.msg)
			t.Errorf("Failed to import data into database: %s", *importErr.err)
		}
	}

	teardownDb(s, t)
}

// Test getting service by URL
func TestGetServiceByUrl(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping DB-involved test (in short mode).")
	}

	s := setupCASServer(t)
	setupDb(s, t)

	// Setup & load services table
	err := s.dbAdapter.SetupServicesTable()
	if err != nil {
		t.Errorf("Failed to setup services table", err)
	}
	casErr := s.dbAdapter.LoadJSONFixture(s.dbAdapter.getDbName(), s.dbAdapter.getServicesTableName(), "fixtures/services.json")
	if casErr != nil {
		if casErr.err != nil {
			t.Log("DB error: %s", casErr.msg)
			t.Errorf("Failed to import data into database: %s", *casErr.err)
		}
	}

	// Attempt to get a service by name
	returnedService, casErr := s.dbAdapter.GetServiceByUrl("localhost:9090/validateCASLogin")
	if casErr != nil {
		if casErr.err != nil {
			t.Log("Internal Error: %v", *casErr.err)
		}
		t.Errorf("Failed to find service that was expected to be present. err: %v", err)
	}

	// Inspect the error, it should have properties we expect
	expectedService := &CASService{
		Name: "test_service",
		Url: "localhost:9090/validateCASLogin",
		AdminEmail: "noone@nowhere.com",
	}

	// Ensure received data matches expected
	if *returnedService != *expectedService {
		t.Errorf("Returned service %v is not equal to expected service %v", returnedService, expectedService)
	}

	teardownDb(s, t)
}
