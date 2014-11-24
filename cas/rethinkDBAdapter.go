package cas

import (
	"bytes"
	r "github.com/dancannon/gorethink"
	"os/exec"
	"path/filepath"
)

type RethinkDBAdapter struct {
	session *r.Session
	dbName  string
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

	return &RethinkDBAdapter{dbSession, c.Config["dbName"]}, nil
}

// Create/Setup all relevant tables in the database
func (db *RethinkDBAdapter) SetupDB() *CASServerError {
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

// Import JSON data into the database
func (db *RethinkDBAdapter) ImportTableDataFromFile(tableName, tablePK, path string) *CASServerError {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return &FailedToImportTableDataFromFile
	}

	cmd := exec.Command("rethinkdb", "import", "--table", tableName, "--pkey", tablePK, "--format", "json", "-f", absPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
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
		return &FailedToTeardownDatabase
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
