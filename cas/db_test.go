package cas

import (
	"testing"
)

// Utility functions for setting up the database
func setupDb(server *CAS, t *testing.T) *CASServerError {
	casErr := server.dbAdapter.Setup()
	if casErr != nil {
		return casErr
	}
	return nil
}

// Utility function for tearing down the database
func teardownDb(server *CAS, t *testing.T) *CASServerError {
	casErr := server.dbAdapter.Teardown()
	if casErr != nil {
		return casErr
	}
	return nil
}

// Test setup and tear down of database (with utility function)
func TestDBSetup(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping DB-involved test (in short mode).")
	}

	casServer := setupTestCASServer(t)

	// If the database already

	setupDb(casServer, t)
	teardownDb(casServer, t)
}

// Test Database checking function
func TestDbExists(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping DB-involved test (in short mode).")
	}

	s := setupTestCASServer(t)
	setupDb(s, t)
	defer teardownDb(s, t)

	exists, casErr := s.dbAdapter.DbExists()
	if casErr != nil {
		if casErr.err != nil {
			t.Logf("INTERNAL ERROR: %v", *casErr.err)
		}
		t.Errorf("Failed to check if database exists, err:", casErr)
	}
	
	if !exists {
		t.Error("Database [%s] should have been set up by setupDb, but does not exist", s.dbAdapter.getDbName())
	}
}

// Test DB property getters with default cas server
func TestAdapterGetters(t *testing.T) {
	s := setupTestCASServer(t)

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

	s := setupTestCASServer(t)
	setupDb(s, t)

	// Setup the services table for importing into
	casErr := s.dbAdapter.SetupServicesTable()
	if casErr != nil {
		if casErr.err != nil {
			t.Logf("DB error: %v", *casErr.err)
		}
		t.Errorf("Failed to setup services table: %v", casErr)
		return
	}

	// Import into the services table
	importErr := s.dbAdapter.LoadJSONFixture(s.dbAdapter.getDbName(), s.dbAdapter.getServicesTableName(), "fixtures/services.json")
	if importErr != nil {
		if importErr.err != nil {
			t.Logf("DB error: %s", importErr.msg)
			t.Errorf("Failed to import data into database: %s", *importErr.err)
			return
		}
	}

	teardownDb(s, t)
}

func dbTestHarness(t *testing.T, fixturesToLoad []StringTuple, dbTestFunc func(*testing.T, *CAS)) {
	if testing.Short() {
		t.Skip("Skipping DB-involved test (in short mode).")
	}

	// Setup CAS Server & DB modified for test
	s := setupTestCASServer(t)
	
	// Setup the test database if it doesn't already exists
	exists, err := s.dbAdapter.DbExists()
	if err != nil {
		t.Errorf("DB Existence check failed, err: %v", err)
		return
	}

	// Setup the DB if it doesn't already exist
	if !exists {
		setupDb(s, t)
	}

	// Possibly load fixtures
	for _, fixtureTuple := range fixturesToLoad {
		tableName, fixturePath := fixtureTuple.First(), fixtureTuple.Second()

		// Setup table in the tuple
		err := s.dbAdapter.SetupTable(tableName)
		if err != nil {
			t.Errorf("Failed to setup %s table: %v", tableName, err)
			return
		}

		// Load fixture into table
		casErr := s.dbAdapter.LoadJSONFixture(s.dbAdapter.getDbName(), tableName, fixturePath)
		if casErr != nil {
			if casErr.err != nil {
				t.Logf("DB error: %s", casErr.msg)
				t.Errorf("Failed to import data into database: %s", *casErr.err)
				return
			}
		}
	}

	// Run test function
	dbTestFunc(t, s)

	// Tear down tables that were loaded
	for _, fixtureTuple := range fixturesToLoad {
		tableName := fixtureTuple.First()
		s.dbAdapter.TeardownTable(tableName)
	}
}

// Test getting service by URL
func TestFindServiceByUrl(t *testing.T) {
	fixturesToLoad := []StringTuple{
		StringTuple{"services", "fixtures/services.json"},
	}
	dbTestHarness(t, fixturesToLoad, func(t *testing.T, s *CAS) {

		// Create the service we expect to find in the fixture
		expectedService := &CASService{
			Name:       "test_service",
			Url:        "localhost:9090/validateCASLogin",
			AdminEmail: "noone@nowhere.com",
		}

		// Attempt to get a service by name
		returnedService, casErr := s.dbAdapter.FindServiceByUrl(expectedService.Url)
		if casErr != nil {
			if casErr.err != nil {
				t.Logf("Internal Error: %v", *casErr.err)
			}
			t.Errorf("Failed to find service that was expected to be present. err: %v", casErr)
			return
		}

		// Ensure received data matches expected
		if returnedService != nil && *returnedService != *expectedService {
			t.Errorf("Returned service %v is not equal to expected service %v", returnedService, expectedService)
			return
		}

	})
}

// Test getting a user by email
func TestFindUserByEmail(t *testing.T) {
	fixturesToLoad := []StringTuple{
		StringTuple{"users", "fixtures/users.json"},
	}
	dbTestHarness(t, fixturesToLoad, func(t *testing.T, s *CAS) {
		// Create the user we're expecting to get back
		expectedUser := &User{
			Email:    "test@test.com",
			Password: "thisisnotarealpassword",
		}

		// Attempt to get a user by name
		returnedUser, casErr := s.dbAdapter.FindUserByEmail(expectedUser.Email)
		if casErr != nil {
			if casErr.err != nil {
				t.Logf("Internal Error: %v", *casErr.err)
			}
			t.Errorf("Failed to find user that was expected to be present. (returnedUser: %v) err: %v", returnedUser, casErr)
			return
		}

		// Ensure received data matches expected
		if returnedUser != nil && !compareUsers(*returnedUser, *expectedUser) {
			t.Errorf("Returned user %v is not equal to expected user %v", returnedUser, expectedUser)
			return
		}

	})
}

// Test adding new user
func TestAddNewUser(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping DB-involved test (in short mode).")
	}

	s := setupTestCASServer(t)
	setupDb(s, t)

	// Setup users table
	casErr := s.dbAdapter.SetupUsersTable()
	if casErr != nil {
		if casErr.err != nil {
			t.Logf("DB Error: %v", *casErr.err)
		}
		t.Errorf("Failed to setup users table, err: %v", casErr)
		return
	}

	// Add the user
	newUser, casErr := s.dbAdapter.AddNewUser("test_user@test.com", "randompassword")
	if casErr != nil {
		t.Errorf("Failed to add new user, err: %v", casErr)
		return
	}

	// Find the user by email
	returnedUser, casErr := s.dbAdapter.FindUserByEmail(newUser.Email)
	if casErr != nil {
		t.Errorf("Failed to find created temporary user, err: %v", casErr)
		return
	}

	if returnedUser != nil && newUser.Email != returnedUser.Email {
		t.Errorf("Newly created user and returned user's emails don't match")
		return
	}

	teardownDb(s, t)
}

// Utility function for creating a service and adding a ticket to it
func addTicketForService(s *CAS, t *testing.T) (*CASTicket, *CASService, *CASServerError) {

	// Setup tickets table
	casErr := s.dbAdapter.SetupTicketsTable()
	if casErr != nil {
		t.Logf("Failed to setup tickets table")
		return nil, nil, casErr
	}

	// Create a new CASTicket to store
	ticket := &CASTicket{
		UserEmail:      "test@test.com",
		UserAttributes: map[string]string{},
		WasSSO:         false,
	}

	mockService := &CASService{
		Url:        "localhost:8080",
		Name:       "mock_service",
		AdminEmail: "noone@nowhere.com",
	}

	ticket, casErr = s.dbAdapter.AddTicketForService(ticket, mockService)
	if casErr != nil {
		t.Logf("Failed to add ticket to database for service [%s]", mockService.Name)
		return nil, nil, casErr
	}

	// Ensure that the ticket has been updated with the right ID
	if ticket != nil && len(ticket.Id) == 0 {
		t.Logf("Received ticket does not have a proper Id attribute set: %v", ticket)
		return nil, nil, &FailedToCreateTicketError
	}

	return ticket, mockService, nil
}

// Test adding ticket for service
func TestAddTicketForService(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping DB-involved test (in short mode).")
	}

	// Setup server and Db
	s := setupTestCASServer(t)
	setupDb(s, t)

	// Add ticket for the service
	addTicketForService(s, t)

	teardownDb(s, t)
}

// Test finding tickets by Id given a service
func TestFindTicketByIdForService(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping DB-involved test (in short mode).")
	}

	// Setup server and Db
	s := setupTestCASServer(t)
	setupDb(s, t)

	// Add ticket for the service
	originalTicket, service, casErr := addTicketForService(s, t)
	if casErr != nil {
		if casErr.err != nil {
			t.Logf("DB err: %v", *casErr.err)
		}
		t.Errorf("Utility function to add ticket for service failed, err: %v", casErr)
		return
	}

	// Find the ticket that was just added
	ticket, err := s.dbAdapter.FindTicketByIdForService(originalTicket.Id, service)
	if err != nil {
		t.Errorf("Failed to find ticket that should have been added: %v", originalTicket)
		return
	}

	// Ensure the tickets are the same
	if ticket != nil && originalTicket != nil && !compareTickets(*ticket, *originalTicket) {
		t.Errorf("Found ticket ( %v ) != original ticket ( %v )", ticket, originalTicket)
		return
	}

	teardownDb(s, t)
}

// Test removing an added tickets for a given user
func TestRemoveTicketsForUser(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping DB-involved test (in short mode).")
	}

	// Setup server and DB
	s := setupTestCASServer(t)
	casErr := setupDb(s, t)
	if casErr != nil {
		if casErr.err != nil {
			t.Logf("DB err: %v", *casErr.err)
		}
		t.Errorf("Utility function to add ticket for service failed, err: %v", casErr)
		return
	}

	// Add a ticket for the user
	ticket, service, casErr := addTicketForService(s, t)
	if casErr != nil {
		if casErr.err != nil {
			t.Logf("DB err: %v", *casErr.err)
		}
		t.Errorf("Utility function to add ticket for service failed, err: %v", casErr)
		return
	}

	// Remove ticket for the user
	err := s.dbAdapter.RemoveTicketsForUserWithService(ticket.UserEmail, service)
	if err != nil {
		t.Errorf("Failed to remove tickets for user with service, err: %v", casErr)
		return
	}

	// Attempt to find ticket (that should have been removed
	ticket, err = s.dbAdapter.FindTicketByIdForService(ticket.Id, service)
	if ticket != nil || err == nil {
		t.Errorf("Found ticket (or did not recieve expected error) that should have been deleted: %v", ticket)
		return
	}

	teardownDb(s, t)
}
