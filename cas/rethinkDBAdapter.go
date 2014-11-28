package cas

import (
	r "github.com/dancannon/gorethink"
	"log"
	"os/exec"
	"path/filepath"
)

type RethinkDBAdapter struct {
	session           *r.Session
	dbName            string
	ticketsTableName  string
	servicesTableName string
	usersTableName    string
}

func (db *RethinkDBAdapter) getDbName() string { return db.dbName }
func (db *RethinkDBAdapter) getTicketsTableName() string { return db.dbName }
func (db *RethinkDBAdapter) getServicesTableName() string { return db.dbName }
func (db *RethinkDBAdapter) getUsersTableName() string { return db.dbName }

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
		session:           dbSession,
		dbName:            c.Config["dbName"],
		ticketsTableName:  "tickets",
		servicesTableName: "services",
		usersTableName:    "users",
	}

	return adapter, nil
}

// Create/Setup all relevant tables in the database
func (db *RethinkDBAdapter) Setup() *CASServerError {

	// Setup the Database
	_, err := r.
		DbCreate(db.dbName).
		Run(db.session)
	if err != nil {
		casError := &FailedToSetupDatabase
		casError.err = &err
		return casError
	}

	return nil
}

// Set up the table that holds services
func (db *RethinkDBAdapter) SetupServicesTable() *CASServerError {
	_, err := r.Db(db.dbName).TableCreate(db.servicesTableName).Run(db.session)
	if err != nil {
		casError := &FailedToSetupDatabase
		casError.err = &err
		return casError
	}
	return nil
}

// Tear down the table that holds services
func (db *RethinkDBAdapter) TeardownServicesTable() *CASServerError {
	_, err := r.Db(db.dbName).TableDrop(db.servicesTableName).Run(db.session)
	if err != nil {
		casError := &FailedToSetupDatabase
		casError.err = &err
		return casError
	}
	return nil
}

// Set up the table that holds tickets
func (db *RethinkDBAdapter) SetupTicketsTable() *CASServerError {
	_, err := r.Db(db.dbName).TableCreate(db.ticketsTableName).Run(db.session)
	if err != nil {
		casError := &FailedToSetupDatabase
		casError.err = &err
		return casError
	}
	return nil
}

// Tear down the table that holds tickets
func (db *RethinkDBAdapter) TeardownTicketsTable() *CASServerError {
	_, err := r.Db(db.dbName).TableDrop(db.servicesTableName).Run(db.session)
	if err != nil {
		casError := &FailedToSetupDatabase
		casError.err = &err
		return casError
	}
	return nil
}

// Set up the table that holds users
func (db *RethinkDBAdapter) SetupUsersTable() *CASServerError {
	_, err := r.Db(db.dbName).TableCreate(db.usersTableName).Run(db.session)
	if err != nil {
		casError := &FailedToSetupDatabase
		casError.err = &err
		return casError
	}
	return nil
}

// Tear down the table that holds users
func (db *RethinkDBAdapter) TeardownUsersTable() *CASServerError {
	_, err := r.Db(db.dbName).TableDrop(db.servicesTableName).Run(db.session)
	if err != nil {
		casError := &FailedToSetupDatabase
		casError.err = &err
		return casError
	}
	return nil
}

// Load database fixture, given intended database name, table and path to fixture file (JSON)
func (db *RethinkDBAdapter) LoadJSONFixture(dbName, tableName, path string) *CASServerError {
	// Get the absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		casError := &FailedToLoadJSONFixture
		casError.err = &err
		return casError
	}

	// Start import command
	importCmd := exec.Command("rethinkdb", "import",
		"--table", dbName+"."+tableName,
		"--format", "json",
		"--force",
		"-f", absPath)
	output, err := importCmd.CombinedOutput()
	if err != nil {
		casError := &FailedToLoadJSONFixture
		casError.msg = string(output)
		casError.err = &err
		return casError
	}

	log.Println("[DB IMPORT]:", string(output))

	return nil
}

// Clear all relevant databases and/or tables
func (db *RethinkDBAdapter) Teardown() *CASServerError {
	_, err := r.
		DbDrop(db.dbName).
		Run(db.session)
	if err != nil {
		casError := &FailedToTeardownDatabase
		casError.err = &err
		return casError
	}

	return nil
}

func (db *RethinkDBAdapter) GetServiceByName(serviceName string) (*CASService, *CASServerError) {
	return &CASService{serviceName, "nobody@nowhere.net"}, nil
}

func (db *RethinkDBAdapter) FindUserByEmail(username string) (*User, *CASServerError) {
	// Find the user
	cursor, err := r.
		Db(db.dbName).
		Table("users").
		Get(username).
		Run(db.session)
	if err != nil {
		return nil, &InvalidEmailAddressError
	}

	// Get the user from the returned cursor
	var returnedUser *User
	err = cursor.One(&returnedUser)
	if err != nil {
		return nil, &InvalidEmailAddressError
	}

	return returnedUser, nil
}

func (db *RethinkDBAdapter) MakeNewTicketForService(service *CASService) (*CASTicket, *CASServerError) {
	// TODO
	return &CASTicket{}, nil
}

func (db *RethinkDBAdapter) RemoveTicketsForUser(username string, service *CASService) *CASServerError {
	return nil
}

func (db *RethinkDBAdapter) FindTicketForService(ticket string, service *CASService) (*CASTicket, *CASServerError) {
	return &CASTicket{}, nil
}

func (db *RethinkDBAdapter) AddNewUser(username, password string) (*User, *CASServerError) {
	user := &User{username, password}

	// Insert user into database
	res, err := r.
		Db(db.dbName).
		Table("users").
		Insert(user, r.InsertOpts{Conflict: "error"}).
		RunWrite(db.session)
	if err != nil {
		return nil, &FailedToCreateUserError
	} else if res.Errors > 0 {
		return nil, &EmailAlreadyTakenError
	}

	return user, nil
}
