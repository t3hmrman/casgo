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

	// If the database already

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
func TestFindServiceByUrl(t *testing.T) {
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
	if returnedService != nil && *returnedService != *expectedService {
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
	if returnedUser != nil && !compareUsers(*returnedUser, *expectedUser) {
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

	// Setup users table
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

	teardownDb(s, t)
}

// Utility function for creating a service
func addTicketForService(s *CAS, t *testing.T) (*CASTicket, *CASService) {
	// Setup tickets table
	err := s.dbAdapter.SetupTicketsTable()
	if err != nil {
		t.Errorf("Failed to setup tickets table", err)
	}

	// Create a new CASTicket to store
	ticket := &CASTicket{
		UserEmail: "test@test.com",
		UserAttributes: map[string]string{},
		WasSSO: false,
	}

	mockService := &CASService{
		Url: "localhost:8080",
		Name: "mock_service",
		AdminEmail: "noone@nowhere.com",
	}

	ticket, err = s.dbAdapter.AddTicketForService(ticket, mockService)
	if err != nil {
		t.Errorf("Failed to add ticket to database for service [%s]", mockService.Name)
	}

	// Ensure that the ticket has been updated with the right ID
	if len(ticket.Id) == 0 {
		t.Errorf("Received ticket does not have a proper Id attribute set: %v", ticket)
	}

	return ticket, mockService
}

// Test adding ticket for service
func TestAddTicketForService(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping DB-involved test (in short mode).")
	}

	// Setup server and Db
	s := setupCASServer(t)
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
	s := setupCASServer(t)
	setupDb(s, t)

	// Add ticket for the service
	originalTicket, service := addTicketForService(s, t)

	// Find the ticket that was just added
	ticket, err := s.dbAdapter.FindTicketByIdForService(originalTicket.Id, service)
	if err != nil {
		t.Errorf("Failed to find ticket that should have been added: %v", originalTicket)
	}

	// Ensure the tickets are the same
	if ticket != nil && originalTicket != nil && !compareTickets(*ticket, *originalTicket) {
		t.Errorf("Found ticket ( %v ) != original ticket ( %v )", ticket, originalTicket)
	}

	teardownDb(s, t)
}


// Test removing an added tickets for a given user
func TestRemoveTicketsForUser(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping DB-involved test (in short mode).")
	}

	// Setup server and DB
	s := setupCASServer(t)
	setupDb(s, t)

	// Add a ticket for the user
	ticket, service := addTicketForService(s, t)

	// Remove ticket for the user
	s.dbAdapter.RemoveTicketsForUser(ticket.UserEmail, service)

	// Attempt to find ticket (that should have been removed
	ticket, err := s.dbAdapter.FindTicketByIdForService(ticket.Id, service)
	if ticket != nil || err == nil {
		t.Errorf("Found ticket (or did not recieve expected error) that should have been deleted: %v", ticket)
	}

	teardownDb(s, t)
}
