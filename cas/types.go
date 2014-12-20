package cas

import (
  r "github.com/dancannon/gorethink"
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
  IsAdmin    bool              `gorethink:"isAdmin" json:"isAdmin"`
}

// Comparison function for Users
func compareUsers(a, b User) bool {
  if &a == &b || (a.Email == b.Email && a.Password == b.Password) {
    return true
  }
  return false
}

type CASService struct {
  Url        string `gorethink:"url" json:"url"`
  Name       string `gorethink:"name" json:"name"`
  AdminEmail string `gorethink:"adminEmail" json:"adminEmail"`
}

// CasGo ticket
type CASTicket struct {
  Id             string            `gorethink:"id,omitempty" json:"id"`
  UserEmail      string            `gorethink:"userEmail" json:"userEmail"`
  UserAttributes map[string]string `gorethink:"userAttributes" json:"userAttributes"`
  WasSSO         bool              `gorethink:"wasSSO" json:"wasSSO"`
}

// Compairson function for CASTickets
func CompareTickets(a, b CASTicket) bool {
  if &a == &b || (a.Id == b.Id && a.UserEmail == b.UserEmail && a.WasSSO == b.WasSSO) {
    return true
  }
  return false
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
  GetDbName() string
  GetTicketsTableName() string
  GetServicesTableName() string
  GetUsersTableName() string
}

// CAS Server
type CAS struct {
  server      *http.Server
  ServeMux    *http.ServeMux
  Config      map[string]string
  Db          CASDBAdapter
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
