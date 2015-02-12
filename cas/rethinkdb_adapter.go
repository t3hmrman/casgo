package cas

import (
	"errors"
	"fmt"
	r "github.com/dancannon/gorethink"
	"os/exec"
	"path/filepath"
)

func (db *RethinkDBAdapter) GetDbName() string            { return db.dbName }
func (db *RethinkDBAdapter) GetTicketsTableName() string  { return db.ticketsTableName }
func (db *RethinkDBAdapter) GetServicesTableName() string { return db.servicesTableName }
func (db *RethinkDBAdapter) GetUsersTableName() string    { return db.usersTableName }
func (db *RethinkDBAdapter) GetApiKeysTableName() string  { return db.apiKeysTableName }

func NewRethinkDBAdapter(c *CAS) (*RethinkDBAdapter, error) {
	// Database setup
	dbSession, err := r.Connect(r.ConnectOpts{
		Address:  c.Config["dbHost"],
		Database: c.Config["dbName"],
	})
	if err != nil {
		return nil, err
	}

	// Create the adapter
	adapter := &RethinkDBAdapter{
		session:              dbSession,
		dbName:               c.Config["dbName"],
		ticketsTableName:     "tickets",
		ticketsTableOptions:  nil,
		servicesTableName:    "services",
		servicesTableOptions: &r.TableCreateOpts{PrimaryKey: "name"},
		usersTableName:       "users",
		usersTableOptions:    &r.TableCreateOpts{PrimaryKey: "email"},
		apiKeysTableName:     "api_keys",
		apiKeysTableOptions:  &r.TableCreateOpts{PrimaryKey: "key"},
		LogLevel:             c.Config["logLevel"],
	}

	return adapter, nil
}

// Check if the database has been setup
func (db *RethinkDBAdapter) DbExists() (bool, *CASServerError) {
	cursor, err := r.
		DbList().
		Run(db.session)
	if err != nil {
		casErr := &DbExistsCheckFailedError
		casErr.err = &err
		return false, casErr
	}

	var response []interface{}
	err = cursor.All(&response)

	// Check that the list contains the database name for the adapter
	for _, listedDb := range response {
		if listedDb == db.dbName {
			return true, nil
		}
	}
	return false, nil
}

// Create/Setup all relevant tables in the database
func (db *RethinkDBAdapter) Setup() *CASServerError {

	// Setup the Database
	_, err := r.
		DbCreate(db.dbName).
		Run(db.session)
	if err != nil {
		casError := &FailedToSetupDatabaseError
		casError.err = &err
		return casError
	}

	// Setup tables
	db.SetupServicesTable()
	db.SetupTicketsTable()
	db.SetupUsersTable()
	db.SetupApiKeysTable()

	return nil
}

func (db *RethinkDBAdapter) teardownTable(tableName string) *CASServerError {
	_, err := r.Db(db.dbName).TableDrop(tableName).Run(db.session)
	if err != nil {
		casError := &FailedToTeardownDatabaseError
		casError.err = &err
		return casError
	}
	return nil
}

func (db *RethinkDBAdapter) setupTable(tableName string, dbOptions interface{}) *CASServerError {
	// Ensure rdbOptions is of correct type, if non-nil
	switch t := dbOptions.(type) {
	case *r.TableCreateOpts:
		return db.createTableWithOptions(tableName, dbOptions.(*r.TableCreateOpts))
	default:
		casError := &FailedToSetupTableError
		err := fmt.Errorf("Unexpected type of dbOptions: %T", t)
		casError.err = &err
		return casError
	}
	return nil
}

func (db *RethinkDBAdapter) createTableWithOptions(tableName string, rdbOptions *r.TableCreateOpts) *CASServerError {
	logMessagef(db.LogLevel, "INFO", "Creating table [%s], options: %v", tableName, rdbOptions)

	// Check again that rdbOptions is not nil, optionally leave out argument
	var err error
	if rdbOptions == nil {

		// Create table with no options
		_, err = r.Db(db.dbName).TableCreate(tableName).Run(db.session)

	} else {

		// Set and get the table options (so that they can be retrieved later
		db.setTableSetupOptions(tableName, rdbOptions)
		options, err := db.getTableSetupOptions(tableName)
		if err != nil {
			casError := &FailedToCreateTableError
			casError.err = &err
			return casError
		}

		// Create table
		_, err = r.Db(db.dbName).TableCreate(tableName, *options).Run(db.session)
	}

	if err != nil {
		casError := &FailedToCreateTableError
		casError.err = &err
		return casError
	}

	return nil
}

// Set up the table that holds services
func (db *RethinkDBAdapter) SetupServicesTable() *CASServerError {
	return db.setupTable(db.servicesTableName, db.servicesTableOptions)
}

// Tear down the table that holds services
func (db *RethinkDBAdapter) TeardownServicesTable() *CASServerError {
	return db.teardownTable(db.servicesTableName)
}

// Set up the table that holds tickets
func (db *RethinkDBAdapter) SetupTicketsTable() *CASServerError {
	return db.setupTable(db.ticketsTableName, db.ticketsTableOptions)
}

// Tear down the table that holds tickets
func (db *RethinkDBAdapter) TeardownTicketsTable() *CASServerError {
	return db.teardownTable(db.ticketsTableName)
}

// Set up the table that holds users
func (db *RethinkDBAdapter) SetupUsersTable() *CASServerError {
	return db.setupTable(db.usersTableName, db.usersTableOptions)
}

// Tear down the table that holds users
func (db *RethinkDBAdapter) TeardownUsersTable() *CASServerError {
	return db.teardownTable(db.usersTableName)
}

// Set up the table that holds apikeys
func (db *RethinkDBAdapter) SetupApiKeysTable() *CASServerError {
	return db.setupTable(db.apiKeysTableName, db.apiKeysTableOptions)
}

// Tear down the table that holds apikeys
func (db *RethinkDBAdapter) TeardownApiKeysTable() *CASServerError {
	return db.teardownTable(db.apiKeysTableName)
}

// Dynamically setup tables - dispatch because each table might have special implementations
func (db *RethinkDBAdapter) SetupTable(tableName string) *CASServerError {
	switch tableName {
	case db.ticketsTableName:
		return db.SetupTicketsTable()
	case db.servicesTableName:
		return db.SetupServicesTable()
	case db.usersTableName:
		return db.SetupUsersTable()
	case db.apiKeysTableName:
		return db.SetupApiKeysTable()
	default:
		casError := &FailedToSetupDatabaseError
		return casError
	}
}

// Dynamically teardown tables - dispatch because each table might have special implementations
func (db *RethinkDBAdapter) TeardownTable(tableName string) *CASServerError {
	switch tableName {
	case db.ticketsTableName:
		return db.TeardownTicketsTable()
	case db.servicesTableName:
		return db.TeardownServicesTable()
	case db.usersTableName:
		return db.TeardownUsersTable()
	default:
		casError := &FailedToTeardownDatabaseError
		return casError
	}
}

// Dynamically get options that were used when setting up a table
func (db *RethinkDBAdapter) getTableSetupOptions(tableName string) (*r.TableCreateOpts, error) {
	switch tableName {
	case db.ticketsTableName:
		return db.ticketsTableOptions, nil
	case db.servicesTableName:
		return db.servicesTableOptions, nil
	case db.usersTableName:
		return db.usersTableOptions, nil
	case db.apiKeysTableName:
		return db.apiKeysTableOptions, nil
	default:
		return nil, errors.New(fmt.Sprintf("Invalid tableName, can't find setup options for table [%s]", tableName))
	}
}

// Dynamically get options that were used when setting up a table
func (db *RethinkDBAdapter) setTableSetupOptions(tableName string, opts *r.TableCreateOpts) error {
	switch tableName {
	case db.ticketsTableName:
		db.ticketsTableOptions = opts
	case db.servicesTableName:
		db.servicesTableOptions = opts
	case db.usersTableName:
		db.usersTableOptions = opts
	case db.apiKeysTableName:
		db.apiKeysTableOptions = opts
	default:
		return errors.New(fmt.Sprintf("Failed to set table setup options for table [%s]", tableName))
	}
	return nil
}

// Load database fixture, given intended database name, table and path to fixture file (JSON)
func (db *RethinkDBAdapter) LoadJSONFixture(dbName, tableName, path string) *CASServerError {
	// Get the absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		casError := &FailedToLoadJSONFixtureError
		casError.err = &err
		return casError
	}

	// Start import command
	importCmd := exec.Command("rethinkdb", "import",
		"--table", dbName+"."+tableName,
		"--format", "json",
		"--force",
		"-f", absPath)

	// Add special options based on table information from setup
	options, err := db.getTableSetupOptions(tableName)
	if err != nil {
		return &CASServerError{Msg: "Failed to find table setup options for table!", err: &err}
	}

	// Check for and apply special options
	if options != nil && options.PrimaryKey != nil {
		importCmd.Args = append(importCmd.Args, "--pkey", options.PrimaryKey.(string))
	}
	logMessagef(db.LogLevel, "INFO", "Args: %v", importCmd.Args)

	// Run the import command
	output, err := importCmd.CombinedOutput()
	if err != nil {
		casError := &FailedToLoadJSONFixtureError
		casError.Msg = string(output)
		casError.err = &err
		return casError
	}

	return nil
}

// Clear all relevant databases and/or tables
func (db *RethinkDBAdapter) Teardown() *CASServerError {
	_, err := r.
		DbDrop(db.dbName).
		Run(db.session)
	if err != nil {
		casError := &FailedToTeardownDatabaseError
		casError.err = &err
		return casError
	}

	return nil
}

// Find a service by given URL (callback URL)
func (db *RethinkDBAdapter) FindServiceByUrl(serviceUrl string) (*CASService, *CASServerError) {
	// Get the first service with the given name
	cursor, err := r.Db(db.dbName).
		Table(db.servicesTableName).
		Filter(map[string]string{"url": serviceUrl}).
		Run(db.session)
	if err != nil {
		casErr := &FailedToLookupServiceByUrlError
		casErr.err = &err
		return nil, casErr
	}

	// Create user object from the returned data cursor
	var returnedService *CASService
	err = cursor.One(&returnedService)
	if err != nil {
		casErr := &FailedToLookupServiceByUrlError
		casErr.err = &err
		return nil, casErr
	}

	return returnedService, nil
}

// Find a user by email address ("username")
func (db *RethinkDBAdapter) FindUserByEmail(email string) (*User, *CASServerError) {
	// Find the user
	cursor, err := r.
		Db(db.dbName).
		Table(db.usersTableName).
		Get(email).
		Run(db.session)
	if err != nil {
		casErr := &FailedToFindUserByEmailError
		casErr.err = &err
		return nil, casErr
	}

	// Get the user from the returned cursor
	var returnedUser *User
	err = cursor.One(&returnedUser)
	if err != nil {
		casErr := &FailedToFindUserByEmailError
		casErr.err = &err
		return nil, casErr
	}

	return returnedUser, nil
}

// Find a user by API secret and key
func (db *RethinkDBAdapter) FindUserByApiKeyAndSecret(key, secret string) (*User, *CASServerError) {
	// Find the user
	cursor, err := r.
		Db(db.dbName).
		Table(db.apiKeysTableName).
		Get(key).
		Run(db.session)
	if err != nil {
		casErr := &FailedToFindUserByApiKeyAndSecretError
		casErr.err = &err
		return nil, casErr
	}

	// Get the user from the returned cursor
	var apiKeyPair *CasgoAPIKeyPair
	err = cursor.One(&apiKeyPair)
	if err != nil {
		casErr := &FailedToFindUserByApiKeyAndSecretError
		casErr.err = &err
		return nil, casErr
	}

	// Return error of the secret is invalid
	if apiKeyPair.Secret != secret {
		casErr := &FailedToFindUserByApiKeyAndSecretError
		return nil, casErr
	}

	return apiKeyPair.User, nil
}

// Add a new user to the database
func (db *RethinkDBAdapter) AddNewUser(username, password string) (*User, *CASServerError) {
	user := &User{
		Email:    username,
		Password: password,
	}

	// Insert user into database
	res, err := r.
		Db(db.dbName).
		Table(db.usersTableName).
		Insert(user, r.InsertOpts{Conflict: "error"}).
		RunWrite(db.session)
	if err != nil || res.Inserted == 0 {
		return nil, &FailedToCreateUserError
	} else if res.Errors > 0 {
		return nil, &EmailAlreadyTakenError
	}

	return user, nil
}

func (db *RethinkDBAdapter) AddNewService(service *CASService) *CASServerError {
	res, err := r.
		Db(db.dbName).
		Table(db.servicesTableName).
		Insert(service, r.InsertOpts{Conflict: "error"}).
		RunWrite(db.session)
	if err != nil || res.Inserted == 0 {
		return &FailedToCreateServiceError
	} else if res.Errors > 0 {
		return &ServiceNameAlreadyTakenError
	}

	// Update the passed in ticket with the ID that was given by the database
	if len(res.GeneratedKeys) > 0 {
		service.Name = res.GeneratedKeys[0]
	}

	return nil
}

// Add new CASTicket to the database for the given service
func (db *RethinkDBAdapter) AddTicketForService(ticket *CASTicket, service *CASService) (*CASTicket, *CASServerError) {
	res, err := r.
		Db(db.dbName).
		Table(db.ticketsTableName).
		Insert(ticket).
		RunWrite(db.session)
	if err != nil || res.Errors > 0 || len(res.GeneratedKeys) == 0 {
		casErr := &FailedToCreateTicketError
		casErr.err = &err
		return nil, casErr
	}

	// Update the passed in ticket with the ID that was given by the database
	if len(res.GeneratedKeys) > 0 {
		ticket.Id = res.GeneratedKeys[0]
	}

	return ticket, nil
}

// Find ticket by Id for a given service
func (db *RethinkDBAdapter) FindTicketByIdForService(ticketId string, service *CASService) (*CASTicket, *CASServerError) {
	cursor, err := r.
		Db(db.dbName).
		Table(db.ticketsTableName).
		Get(ticketId).
		Run(db.session)
	if err != nil || cursor.IsNil() {
		casErr := &FailedToFindTicketError
		casErr.err = &err
		return nil, casErr
	}

	// Create CASTicket from result
	var returnedTicket *CASTicket
	err = cursor.One(&returnedTicket)
	if err != nil {
		casErr := &FailedToFindTicketError
		casErr.err = &err
		return nil, casErr
	}

	return returnedTicket, nil
}

// Remove tickets for a given user under a given service
func (db *RethinkDBAdapter) RemoveTicketsForUserWithService(email string, service *CASService) *CASServerError {
	_, err := r.
		Db(db.dbName).
		Table(db.ticketsTableName).
		Filter(map[string]string{"userEmail": email}).
		Delete().
		Run(db.session)
	if err != nil {
		casErr := &FailedToDeleteTicketsForUserError
		casErr.err = &err
		return casErr
	}

	return nil
}

// Remove a service by name (pkey)
func (db *RethinkDBAdapter) RemoveServiceByName(name string) *CASServerError {
	if len(name) == 0 {
		return &InvalidServiceNameError
	}

	_, err := r.
		Db(db.dbName).
		Table(db.servicesTableName).
		Get(name).
		Delete().
		Run(db.session)
	if err != nil {
		casErr := &FailedToDeleteServiceError
		casErr.err = &err
		return casErr
	}

	return nil
}

// Remove a user by email (pkey)
func (db *RethinkDBAdapter) RemoveUserByEmail(email string) *CASServerError {
	if len(email) == 0 {
		return &InvalidUserEmailError
	}

	_, err := r.
		Db(db.dbName).
		Table(db.usersTableName).
		Get(email).
		Delete().
		Run(db.session)
	if err != nil {
		casErr := &FailedToDeleteUserError
		casErr.err = &err
		return casErr
	}

	return nil
}

// Update service with a similar name to the passed in service (key)
func (db *RethinkDBAdapter) UpdateService(service *CASService) *CASServerError {
	if len(service.Name) == 0 {
		return &InvalidServiceNameError
	}

	res, err := r.
		Db(db.dbName).
		Table(db.servicesTableName).
		Get(service.Name).
		Update(service, r.UpdateOpts{ReturnChanges: true}).
		RunWrite(db.session)
	if err != nil || res.Replaced == 0 || len(res.Changes) == 0 {
		casErr := &FailedToUpdateServiceError
		casErr.err = &err
		return casErr
	}

	return nil
}

// Update user with a similar name to the passed in user (key)
func (db *RethinkDBAdapter) UpdateUser(user *User) *CASServerError {
	if len(user.Email) == 0 {
		return &InvalidUserEmailError
	}

	res, err := r.
		Db(db.dbName).
		Table(db.usersTableName).
		Get(user.Email).
		Update(user, r.UpdateOpts{ReturnChanges: true}).
		RunWrite(db.session)
	if err != nil || res.Replaced == 0 || len(res.Changes) == 0 {
		casErr := &FailedToUpdateUserError
		casErr.err = &err
		return casErr
	}

	return nil
}

// Get all services
func (db *RethinkDBAdapter) GetAllServices() ([]CASService, *CASServerError) {
	cursor, err := r.
		Db(db.dbName).
		Table(db.servicesTableName).
		Run(db.session)
	if err != nil {
		casErr := &FailedToListServicesError
		casErr.err = &err
		return nil, casErr
	}

	var services []CASService
	err = cursor.All(&services)
	if err != nil {
		casErr := &FailedToListServicesError
		casErr.err = &err
		return nil, casErr
	}

	return services, nil
}

// Get all users
func (db *RethinkDBAdapter) GetAllUsers() ([]User, *CASServerError) {
	cursor, err := r.
		Db(db.dbName).
		Table(db.usersTableName).
		Without("password").
		Run(db.session)
	if err != nil {
		casErr := &FailedToListUsersError
		casErr.err = &err
		return nil, casErr
	}

	var users []User
	err = cursor.All(&users)
	if err != nil {
		casErr := &FailedToListUsersError
		casErr.err = &err
		return nil, casErr
	}

	return users, nil
}
