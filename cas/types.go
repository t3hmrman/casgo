package cas

import (
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
	"log"
	"net/http"
)

type User struct {
	Email    string `gorethink:"email" json:"email"`
	Password string `gorethink:"password" json:"password"`
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
	Name        string `gorethink:"name" json:"name"`
	AdminEmail string `gorethink:"adminEmail" json:"adminEmail"`
}

// Function to create a CASService object from a parsed generic map[string]interface{}
func createCASServiceFromGenericObject(generic map[string]interface{}) CASService {
	return CASService{
		Url:        generic["url"].(string),
		AdminEmail: generic["adminEmail"].(string),
	}
}

type CASTicket struct {
	serviceUrl        string `gorethink:"serviceUrl" json:"serviceUrl"`
	wasFromSSOSession bool   `gorethink:"wasFromSSOSession" json:"wasFromSSOSession"`
}

// Function to create a CASTicket object from a parsed generic map[string]interface{}
func createCASTicketFromGenericObject(generic map[string]interface{}) CASTicket {
	return CASTicket{
		serviceUrl:        generic["serviceUrl"].(string),
		wasFromSSOSession: generic["wasFromSSOSession"].(bool),
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
	SetupServicesTable() *CASServerError
	TeardownServicesTable() *CASServerError
	SetupUsersTable() *CASServerError
	TeardownUsersTable() *CASServerError
	SetupTicketsTable() *CASServerError
	TeardownTicketsTable() *CASServerError

	// Fixture loading utility function
	LoadJSONFixture(string, string, string) *CASServerError

	// App functions
	GetServiceByUrl(string) (*CASService, *CASServerError)
	FindUserByEmail(string) (*User, *CASServerError)
	MakeNewTicketForService(service *CASService) (*CASTicket, *CASServerError)
	RemoveTicketsForUser(string, *CASService) *CASServerError
	FindTicketForService(string, *CASService) (*CASTicket, *CASServerError)
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
}
