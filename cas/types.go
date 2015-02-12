package cas

import (
	r "github.com/dancannon/gorethink"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
	"net/http"
)

// Small string tuple class implementation (see util.go)
type StringTuple [2]string

// CasGo user
type User struct {
	Email      string            `gorethink:"email" json:"email"`
	Attributes map[string]string `gorethink:"attributes" json:"attributes"`
	Password   string            `gorethink:"password" json:"password"`
	Services   []CASService      `gorethink:"services" json:"services"`
	IsAdmin    bool              `gorethink:"isAdmin" json:"isAdmin"`
}

// Enforce schema for Users
func (u *User) IsValid() bool {
	return len(u.Email) > 0 && len(u.Password) > 0
}

// Enforce lax schema for user updates (as they may not include Password field)
// at least the email must be present (used when getting the user, as it is the PK)
func (u *User) IsValidUpdate() bool {
	return len(u.Email) > 0
}

// Comparison function for Users
func compareUsers(a, b User) bool {
	if &a == &b || (a.Email == b.Email && a.Password == b.Password) {
		return true
	}
	return false
}

// CasGo registered service
type CASService struct {
	Url        string `gorethink:"url" json:"url"`
	Name       string `gorethink:"name" json:"name"`
	AdminEmail string `gorethink:"adminEmail" json:"adminEmail"`
}

// Enforce schema for CASService
func (s *CASService) IsValid() bool {
	return len(s.Url) > 0 && len(s.Name) > 0 && len(s.AdminEmail) > 0
}

// Enforce lax schema for CASService updates (as they may not include some otherwise required fields)
// at least the name must be present (used when getting the service, as it is the PK)
func (s *CASService) IsValidUpdate() bool {
	return len(s.Name) > 0
}

// CasGo ticket
type CASTicket struct {
	Id             string            `gorethink:"id,omitempty" json:"id"`
	UserEmail      string            `gorethink:"userEmail" json:"userEmail"`
	UserAttributes map[string]string `gorethink:"userAttributes" json:"userAttributes"`
	WasSSO         bool              `gorethink:"wasSSO" json:"wasSSO"`
}

// CasGo API keypair
type CasgoAPIKeyPair struct {
	Key    string `gorethink:"key" json:"key"`
	Secret string `gorethink:"secret" json:"secret"`
	User   *User  `gorethink:"user" json:"user"`
}

// Compairson function for CASTickets
func CompareTickets(a, b CASTicket) bool {
	if &a == &b || (a.Id == b.Id && a.UserEmail == b.UserEmail && a.WasSSO == b.WasSSO) {
		return true
	}
	return false
}

type CASServerError struct {
	Msg          string // Message string
	HttpCode     int    // HTTP error code, if applicable
	CasgoErrCode int    // CASGO specific error code
	err          *error // Actual error that was thrown (if any)
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
	FindUserByApiKeyAndSecret(string, string) (*User, *CASServerError)
	AddTicketForService(ticket *CASTicket, service *CASService) (*CASTicket, *CASServerError)
	RemoveTicketsForUserWithService(string, *CASService) *CASServerError
	FindTicketByIdForService(string, *CASService) (*CASTicket, *CASServerError)
	AddNewUser(string, string) (*User, *CASServerError)

	// REST API functions (CRUD)
	GetAllUsers() ([]User, *CASServerError)
	UpdateUser(*User) *CASServerError
	RemoveUserByEmail(string) *CASServerError

	GetAllServices() ([]CASService, *CASServerError)
	AddNewService(*CASService) *CASServerError
	RemoveServiceByName(string) *CASServerError
	UpdateService(*CASService) *CASServerError

	// Property getter utility functions
	GetDbName() string
	GetTicketsTableName() string
	GetServicesTableName() string
	GetUsersTableName() string
	GetApiKeysTableName() string
}

type CasgoFrontendAPI interface {
	HookupAPIEndpoints(*mux.Router)

	// Services Endpoint
	GetServices(http.ResponseWriter, *http.Request)
	RemoveService(http.ResponseWriter, *http.Request)
	CreateService(http.ResponseWriter, *http.Request)

	listSessionUserServices(http.ResponseWriter, *http.Request)
	SessionsHandler(http.ResponseWriter, *http.Request)
}

// CAS Server
type CAS struct {
	server      *http.Server
	ServeMux    *mux.Router
	Config      map[string]string
	Db          CASDBAdapter
	Api         CasgoFrontendAPI
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
	apiKeysTableName     string
	apiKeysTableOptions  *r.TableCreateOpts
	LogLevel             string
}

// CasGo frontend RESTful API
type FrontendAPI struct {
	casServer *CAS
}
