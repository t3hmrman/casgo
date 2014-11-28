package cas

import (
	"testing"
)

// Utility functions for tearing down/setting up the database
func setup(server *CAS, t *testing.T) {
	setupErr := server.dbAdapter.Setup()
	if setupErr != nil {
		if setupErr != nil { t.Errorf("Failed to set up database: %s", *setupErr.err) }
	}
}

func teardown(server *CAS, t *testing.T) {
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
	setup(casServer, t)
	teardown(casServer, t)
}

// Test database import functionality (implemented LoadJSONFixture function should not fail to load fixtures)
func TestLoadJSONFixture(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping DB-involved test (in short mode).")
	}

	server := setupCASServer(t)
	setup(server, t)
	importErr := server.dbAdapter.LoadJSONFixture("casgo_test", "services", "fixtures/services.json")
	if importErr != nil {
		if importErr.err != nil {
			t.Log("DB error: %s", importErr.msg)
			t.Errorf("Failed to import data into database: %s", *importErr.err)
		}
	}

	teardown(server, t)
}
