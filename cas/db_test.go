package cas

import (
	"testing"
)

// Utility functions for tearing down/setting up the database
func setupDb(server *CAS, t *testing.T) {
	setupErr := server.dbAdapter.Setup()
	if setupErr != nil {
		if setupErr != nil {
			t.Errorf("Failed to set up database: %s", *setupErr.err)
			return
		}
	}
}

func teardownDb(server *CAS, t *testing.T) {
	teardownErr := server.dbAdapter.Teardown()
	if teardownErr != nil {
		if teardownErr.err != nil {
			t.Errorf("Failed to tear down database: %s", *teardownErr.err)
			return
		}
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

func dbTestHarness(t *testing.T, tablesToSetup []string, fixturesToLoad []string, dbTestFunc func(*testing.T, *CAS)) {
	if testing.Short() {
		t.Skip("Skipping DB-involved test (in short mode).")
	}

	// Do setup
	s := setupCASServer(t)
	setupDb(s, t)
	defer teardownDb(s, t)

	// Possibly execute setup functions (mostly to setup tables)
	for _, tableName := range tablesToSetup {
		err := s.dbAdapter.SetupTable(tableName)
		if err != nil {
			t.Errorf("Failed to setup %s table: %v", tableName, err)
		}
	}

	// Possibly load fixtures
	for _, fixturePath := range fixturesToLoad {
		casErr := s.dbAdapter.LoadJSONFixture(s.dbAdapter.getDbName(), s.dbAdapter.getServicesTableName(), fixturePath)
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
}

// Test getting service by URL
func TestFindServiceByUrl(t *testing.T) {
	tablesToSetup := []string{"services"}
	fixturesToLoad := []string{"fixtures/services.json"}
	dbTestHarness(t, tablesToSetup, fixturesToLoad, func(t *testing.T, s *CAS) {

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
	if testing.Short() {
		t.Skip("Skipping DB-involved test (in short mode).")
	}

	s := setupCASServer(t)
	setupDb(s, t)

	// Setup & load users table
	err := s.dbAdapter.SetupUsersTable()
	if err != nil {
		t.Errorf("Failed to setup users table, err: %v", err)
		return
	}

	casErr := s.dbAdapter.LoadJSONFixture(s.dbAdapter.getDbName(), s.dbAdapter.getUsersTableName(), "fixtures/users.json")
	if casErr != nil {
		if casErr.err != nil {
			t.Logf("Failed to import data into database: %s", *casErr.err)
		}
		t.Errorf("Failed to load JSON fixture: %s", casErr)
		return
	}

	// Inspect the error, it should have properties we expect
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
		if casErr.err != nil {
			t.Logf("DB Error: %v", *casErr.err)
		}
		t.Errorf("Failed to setup tickets table, %v", casErr)
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
		t.Errorf("Failed to add ticket to database for service [%s]", mockService.Name)
		return nil, nil, casErr
	}

	// Ensure that the ticket has been updated with the right ID
	if ticket != nil && len(ticket.Id) == 0 {
		t.Errorf("Received ticket does not have a proper Id attribute set: %v", ticket)
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
	s := setupCASServer(t)
	setupDb(s, t)

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
	s.dbAdapter.RemoveTicketsForUserWithService(ticket.UserEmail, service)

	// Attempt to find ticket (that should have been removed
	ticket, err := s.dbAdapter.FindTicketByIdForService(ticket.Id, service)
	if ticket != nil || err == nil {
		t.Errorf("Found ticket (or did not recieve expected error) that should have been deleted: %v", ticket)
		return
	}

	teardownDb(s, t)
}
