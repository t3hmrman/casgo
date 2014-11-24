package cas

import (
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
	"net/http"
)

type User struct {
	Email    string `gorethink:"email"`
	Password string `gorethink:"password"`
}

type CASService struct {
	Url               string `gorethink:"url"`
	AdminstratorEmail string `gorethink:"admin_email"`
}

type CASTicket struct {
	serviceUrl          string `gorethink:"serviceUrl"`
	wasFromSSOSession   bool `gorethink:"wasFromSSOSession"`
}

type CASServerError struct {
	msg string
	http_code int
	err_code int
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
	GetServiceByName(string) (*CASService, *CASServerError)
	FindUserByEmail(string) (*User, *CASServerError)
	MakeNewTicketForService(service *CASService) (*CASTicket, *CASServerError)
	RemoveTicketsForUser(string, *CASService) *CASServerError
	FindTicketForService(string, *CASService) (*CASTicket, *CASServerError)
	AddNewUser(string, string) (*User, *CASServerError)
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

