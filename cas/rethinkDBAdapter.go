package cas

import (
	"encoding/json"
	r "github.com/dancannon/gorethink"
	"log"
	"os"
	"path/filepath"
)

type RethinkDBAdapter struct {
	session          *r.Session
	dbName           string
	ticketTableName  string
	serviceTableName string
	userTableName    string
}

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
		session:          dbSession,
		dbName:           c.Config["dbName"],
		ticketTableName:  "tickets",
		serviceTableName: "services",
		userTableName:    "users",
	}

	return adapter, nil
}

// Create/Setup all relevant tables in the database
func (db *RethinkDBAdapter) SetupDB() *CASServerError {

	// Setup the Database
	_, err := r.
		DbCreate(db.dbName).
		Run(db.session)
	if err != nil {
		casError := &FailedToSetupDatabase
		casError.err = &err
		return casError
	}

	// Setup required tables
	_, ticketTableErr := r.Db(db.dbName).TableCreate(db.ticketTableName).Run(db.session)
	_, userTableErr := r.Db(db.dbName).TableCreate(db.userTableName).Run(db.session)
	_, serviceTableErr := r.Db(db.dbName).TableCreate(db.serviceTableName).Run(db.session)
	if ticketTableErr != nil || userTableErr != nil || serviceTableErr != nil {
		casError := &FailedToSetupDatabase
		return casError
	}

	return nil
}

// Import JSON data into the database
func (db *RethinkDBAdapter) ImportTableDataFromFile(tableName, path string) *CASServerError {
	// Get the absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		casError := &FailedToImportTableDataFromFile
		casError.err = &err
		return casError
	}

	// Open the file
	file, err := os.Open(absPath)
	if err != nil {
		casError := &FailedToImportTableDataFromFile
		casError.err = &err
		return casError
	}

	// Read the JSON data and upload all of it into the specified table
	// Build inert query to use to insert all the objects into the database
	insertQuery := r.Db(db.dbName).Table(tableName)
	decoder := json.NewDecoder(file)
	for {
		var v map[string]interface{}
		if err := decoder.Decode(&v); err != nil {
			break
		}
		log.Print("Want to insert:", v)
		insertQuery.Insert(v)
	}

	// Run the insert query
	_, err = insertQuery.Run(db.session)
	if err != nil {
		casError := &FailedToImportTableDataFromFile
		casError.err = &err
		return casError
	}

	return nil
}

// Clear all relevant databases and/or tables
func (db *RethinkDBAdapter) TeardownDB() *CASServerError {
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
