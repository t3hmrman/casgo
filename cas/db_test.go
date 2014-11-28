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

// Test DB property getters with default cas server
func TestAdapterGetters(t *testing.T) {
	s := setupCASServer(t)

	actual, expected := s.dbAdapter.getDbName(), s.Config["dbName"]
	if actual != expected {
		t.Errorf("Expected GetDbName to return [%s], returned [%s]", actual, expected)
		return
	}

	actual, expected = s.dbAdapter.getUsersTableName(), "users"
	if actual != expected {
		t.Errorf("Expected getUsersTableName to return [%s], returned [%s]", actual, expected)
		return
	}

	actual, expected = s.dbAdapter.getServicesTableName(), "services"
	if actual != expected {
		t.Errorf("Expected getServicesTableName to return [%s], returned [%s]", actual, expected)
		return
	}

	actual, expected = s.dbAdapter.getTicketsTableName(), "tickets"
	if actual != expected {
		t.Errorf("Expected getTicketsTableName to return [%s], returned [%s]", actual, expected)
		return
	}

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

	// Create the service we expect to find in the fixture
	expectedService := &CASService{
		Name: "test_service",
		Url: "localhost:9090/validateCASLogin",
		AdminEmail: "noone@nowhere.com",
	}


	// Attempt to get a service by name
	returnedService, casErr := s.dbAdapter.FindServiceByUrl(expectedService.Url)
	if casErr != nil {
		if casErr.err != nil {
			t.Log("Internal Error: %v", *casErr.err)
		}
		t.Errorf("Failed to find service that was expected to be present. err: %v", err)
	}

	// Ensure received data matches expected
	if *returnedService != *expectedService {
		t.Errorf("Returned service %v is not equal to expected service %v", returnedService, expectedService)
	}

	teardownDb(s, t)
}

// Test getting a user by email
func TestFindUserByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping DB-involved test (in short mode).")
	}

	s := setupCASServer(t)
	setupDb(s, t)

	// Setup & load users table
	err := s.dbAdapter.SetupUsersTable()
	if err != nil {
		t.Errorf("Failed to setup users table", err)
	}
	casErr := s.dbAdapter.LoadJSONFixture(s.dbAdapter.getDbName(), s.dbAdapter.getUsersTableName(), "fixtures/users.json")
	if casErr != nil {
		if casErr.err != nil {
			t.Errorf("Failed to import data into database: %s", *casErr.err)
		}
			t.Log("Failed to load JSON fixture: %s", casErr.msg)
	}

	// Inspect the error, it should have properties we expect
	expectedUser := &User{
		Email: "test@test.com",
		Password: "thisisnotarealpassword",
	}

	// Attempt to get a user by name
	returnedUser, casErr := s.dbAdapter.FindUserByEmail(expectedUser.Email)
	if casErr != nil {
		if casErr.err != nil {
			t.Log("Internal Error: %v", *casErr.err)
		}
		t.Errorf("Failed to find user that was expected to be present. err: %v", err)
	}

	// Ensure received data matches expected
	if returnedUser != nil && *returnedUser != *expectedUser {
		t.Errorf("Returned user %v is not equal to expected user %v", returnedUser, expectedUser)
	}

	teardownDb(s, t)
}

// Test adding new user
func TestAddNewUser(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping DB-involved test (in short mode).")
	}

	s := setupCASServer(t)
	setupDb(s, t)

	// Setup & load users table
	err := s.dbAdapter.SetupUsersTable()
	if err != nil {
		t.Errorf("Failed to setup users table", err)
	}

	// Add the user
	newUser, casErr := s.dbAdapter.AddNewUser("test_user@test.com", "randompassword")
	if casErr != nil {
		t.Errorf("Failed to add new user", casErr)
	}

	// Find the user by email
	returnedUser, casErr := s.dbAdapter.FindUserByEmail(newUser.Email)
	if casErr != nil {
		t.Errorf("Failed to find created temporary user", casErr)
	}

	if newUser.Email != returnedUser.Email {
		t.Errorf("Newly created user and returned user's emails don't match")
	}
}
