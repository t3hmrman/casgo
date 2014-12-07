package cas

import (
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
	"log"
	"net/http"
	r "github.com/dancannon/gorethink"
)

// Small string tuple class implementation (see util.go)
type StringTuple [2]string

// CasGo user
type User struct {
	Email      string            `gorethink:"email" json:"email"`
	Attributes map[string]string `gorethink:"attributes" json:"attributes"`
	Password   string            `gorethink:"password" json:"password"`
}

// Comparison function for Users
func compareUsers(a, b User) bool {
	if &a == &b || (a.Email == b.Email && a.Password == b.Password) {
		return true
	}
	return false
}

// Function to create a user object from a parsed generic map[string]interface{}
func createUserFromGenericObject(generic map[string]interface{}) User {
	return User{
		Email:    generic["email"].(string),
		Password: generic["password"].(string),
	}
}

type CASService struct {
	Url        string `gorethink:"url" json:"url"`
	Name       string `gorethink:"name" json:"name"`
	AdminEmail string `gorethink:"adminEmail" json:"adminEmail"`
}

// Function to create a CASService object from a parsed generic map[string]interface{}
func createCASServiceFromGenericObject(generic map[string]interface{}) CASService {
	return CASService{
		Url:        generic["url"].(string),
		AdminEmail: generic["adminEmail"].(string),
	}
}

// CasGo ticket
type CASTicket struct {
	Id             string            `gorethink:"id,omitempty" json:"id"`
	UserEmail      string            `gorethink:"userEmail" json:"userEmail"`
	UserAttributes map[string]string `gorethink:"userAttributes" json:"userAttributes"`
	WasSSO         bool              `gorethink:"wasSSO" json:"wasSSO"`
}

// Compairson function for CASTickets
func compareTickets(a, b CASTicket) bool {
	if &a == &b || (a.Id == b.Id && a.UserEmail == b.UserEmail && a.WasSSO == b.WasSSO) {
		return true
	}
	return false
}

// Function to create a CASTicket object from a parsed generic map[string]interface{}
func createCASTicketFromGenericObject(generic map[string]interface{}) CASTicket {
	return CASTicket{
		WasSSO: generic["wasFromSSOSession"].(bool),
	}
}

// Utility function to translate a generic map to a proper databse type
func translateGenericObjectToDBStruct(tableName string, obj map[string]interface{}) interface{} {
	// Determine the function that should be used to translate the object
	switch tableName {
	case "services":
		return createCASServiceFromGenericObject(obj)
		break
	case "tickets":
		return createCASTicketFromGenericObject(obj)
		break
	case "users":
		return createUserFromGenericObject(obj)
		break
	default:
		log.Fatal("Invalid table name, could not find generic-to-object conversion function")
	}
	return nil
}

type CASServerError struct {
	msg        string // Message string
	httpCode   int    // HTTP error code, if applicable
	casErrCode int    // CASGO specific error code
	err        *error // Actual error that was thrown (if any)
}

// CAS server interface
type CASServer interface {
	HandleLogin(w http.ResponseWriter, r *http.Request)
	HandleLogout(w http.ResponseWriter, r *http.Request)
	HandleRegister(w http.ResponseWriter, r *http.Request)
	HandleValidate(w http.ResponseWriter, r *http.Request)
	HandleServiceValidate(w http.ResponseWriter, r *http.Request)
	HandleProxyValidate(w http.ResponseWriter, r *http.Request)
	HandleProxy(w http.ResponseWriter, r *http.Request)
}

// CAS DB interface
type CASDBAdapter interface {
	// Database setup & teardown logic
	Setup() *CASServerError
	Teardown() *CASServerError

	// Table setup & teardown logic
	DbExists() (bool, *CASServerError)
	SetupTable(string) *CASServerError
	TeardownTable(string) *CASServerError
	SetupServicesTable() *CASServerError
	TeardownServicesTable() *CASServerError
	SetupUsersTable() *CASServerError
	TeardownUsersTable() *CASServerError
	SetupTicketsTable() *CASServerError
	TeardownTicketsTable() *CASServerError

	// Fixture loading utility function
	LoadJSONFixture(string, string, string) *CASServerError

	// App functions
	FindServiceByUrl(string) (*CASService, *CASServerError)
	FindUserByEmail(string) (*User, *CASServerError)
	AddTicketForService(ticket *CASTicket, service *CASService) (*CASTicket, *CASServerError)
	RemoveTicketsForUserWithService(string, *CASService) *CASServerError
	FindTicketByIdForService(string, *CASService) (*CASTicket, *CASServerError)
	AddNewUser(string, string) (*User, *CASServerError)

	// Property getter utility functions
	getDbName() string
	getTicketsTableName() string
	getServicesTableName() string
	getUsersTableName() string
}

// CAS Server
type CAS struct {
	server      *http.Server
	serveMux    *http.ServeMux
	Config      map[string]string
	dbAdapter   CASDBAdapter
	render      *render.Render
	cookieStore *sessions.CookieStore
	LogLevel    int
}

// RethinkDB Adapter
type RethinkDBAdapter struct {
	session              *r.Session
	dbName               string
	ticketsTableName     string
	ticketsTableOptions  *r.TableCreateOpts
	servicesTableName    string
	servicesTableOptions *r.TableCreateOpts
	usersTableName       string
	usersTableOptions    *r.TableCreateOpts
	LogLevel             string
}
