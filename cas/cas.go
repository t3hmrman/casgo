package cas

/**
 * CAS protocol V3 API
 */

import (
	"code.google.com/p/go.crypto/bcrypt"
	r "github.com/dancannon/gorethink"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
	"net/http"
	"strings"
	"log"
)

type User struct {
	Email    string `gorethink:"email"`
	Password string `gorethink:"password"`
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

// CAS Server
type CAS struct {
	config  *CASServerConfig
	render  *render.Render
	cookieStore *sessions.CookieStore
}

func New(config *CASServerConfig) *CAS {
	r := render.New(render.Options{Directory: config.TemplatesDirectory})
	cookieStore := sessions.NewCookieStore([]byte(config.CookieSecret))
	cookieStore.Options = &sessions.Options{
		Path: "/",
		MaxAge: 86400 * 7,
		HttpOnly: true,
	}
	c := &CAS{config, r, cookieStore}
	return c
}

// (Optional) Handles Index route
func (c *CAS) HandleIndex(w http.ResponseWriter, req *http.Request) {
	c.render.HTML(w, http.StatusOK, "index", map[string]string{"CompanyName": c.config.CompanyName})
}

// Credential acceptor endpoint (requestor is Handled in main)
func (c *CAS) HandleLogin(w http.ResponseWriter, req *http.Request) {
	// Generate context
	context := map[string]string{	"CompanyName": c.config.CompanyName }

	// Exit early if the user is already logged in (in session)
	session, _ := c.cookieStore.Get(req, "casgo-session")
	if currentUserEmail,ok := session.Values["currentUserEmail"]; ok {
		context["currentUserEmail"] = currentUserEmail.(string)
		c.render.HTML(w, http.StatusOK, "login", context)
		return
	}

	// Trim and lightly pre-process/validate email/password
	email := strings.TrimSpace(strings.ToLower(req.FormValue("email")))
	password := strings.TrimSpace(strings.ToLower(req.FormValue("password")))
	if email == "" || password == "" {
		c.render.HTML(w, http.StatusOK, "login", context)
		return
	}

	// Find, and validate user credentials
	// validateUserCredentials returns CASServerError if there is an error
	returnedUser, err := c.validateUserCredentials(email,password)
	if err != nil {
		context["Error"] = err.msg
		c.render.HTML(w, err.http_code, "login", context)
		return
	}

	// Save session in cookies
	session, _ = c.cookieStore.Get(req, "casgo-session")
	session.Values["currentUserEmail"] = returnedUser.Email
	sessionSaveErr := session.Save(req, w)
	if sessionSaveErr != nil {
		log.Fatal("Failed to save logged in user to session:", err)
	}

	// if the user has logged in and there was a service, do stuff
	// if req.FormValue("service"), ok; ok {
	//	handleServiceLogin(w, req)
	//	return
	// }

	context["Success"] = "Successful log in! Redirecting to services page..."
	context["currentUserEmail"] = returnedUser.Email
	c.render.HTML(w, http.StatusOK, "login", context)
}

// Validate user credentials
// Returns a valid user object if validation succeeds
func (c *CAS) validateUserCredentials(email string, password string) (*User, *CASServerError) {

	// Find the user
	cursor, err := r.Db(c.config.DBName).Table("users").Get(email).Run(c.config.RDBSession)
	if err != nil { return nil, &InvalidEmailAddressError }

	// Get the user from the returned cursor
	var returnedUser *User
	err = cursor.One(&returnedUser)
	if err != nil {	return nil, &InvalidEmailAddressError }

	// Check hash
	err = bcrypt.CompareHashAndPassword([]byte(returnedUser.Password), []byte(password))
	if err != nil { return nil, &InvalidCredentialsError }

	// Successful validation
	return returnedUser, nil
}

// Handle login accesses (credential acceptance -- verifying attempt to log in)
func (c *CAS) handleLoginCredentialAccept(w http.ResponseWriter, req *http.Request) {

}

// Handle login accesses (credential requests -- user logging in)
func (c *CAS) handleLoginCredentialRequest(w http.ResponseWriter, req *http.Request) {
	context := map[string]string{}

	// Trim and lightly pre-process/validate service
	service := strings.TrimSpace(strings.ToLower(req.FormValue("service")))
	gateway := strings.TrimSpace(strings.ToLower(req.FormValue("gateway")))
	renew := strings.TrimSpace(strings.ToLower(req.FormValue("renew")))
	method := strings.TrimSpace(strings.ToLower(req.FormValue("method")))

	if gateway == "true" && renew == "true" {
		// Both gateway and  cannot be set -- Croak here? Maybe also take renew over gateway (as docs suggest?)
		context["Error"] = "Invalid Request: Both gateway and renew options specified"
			c.render.HTML(w, http.StatusBadRequest, "login", context)
	} 

	if gateway == "true" {
		// If gateway is set, CAS will try to authenticate with non-interactive means (ex. LDAP)
		// If no CAS session and no non-interactive means, then redirect with no ticket parameter to service URL

	} else if renew == "true" {
		// If renew is set, automatic sign on is disabled, user must present credentials regardless of whether a sign on session exists
		// Ignore gateway if renew is set.

	}

	// Look up the service by URL -- URL encoded
	// Verify that it exists for possible redirect

	// If renew is set, gateway cannot be

	//

		context["Success"] = "Handling service login, service:" + service + ", gateway:" + gateway + ", method: " + method
		c.render.HTML(w, http.StatusOK, "login", context)
}

// Endpoint for destroying CAS sessions (logging out)
func (c *CAS) HandleRegister(w http.ResponseWriter, req *http.Request) {
	context := map[string]string{	"CompanyName": c.config.CompanyName }

	// Show login page if credentials are not provided, attempt login otherwise
	email := strings.TrimSpace(strings.ToLower(req.FormValue("email")))
	password := strings.TrimSpace(strings.ToLower(req.FormValue("password")))

	// Exit early if email/password are empty
	if email == "" || password == "" {
		c.render.HTML(w, http.StatusOK, "register", context)
		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10) // Default cost
	if err != nil {
		context["Error"] = "Registration failed... Please contact server administrator"
		c.render.HTML(w, http.StatusInternalServerError, "register", context)
		return
	}

	// Create new user object
	newUser := map[string]string{
		"email":    email,
		"password": string(encryptedPassword),
	}

	res, err := r.Db(c.config.DBName).Table("users").Insert(newUser, r.InsertOpts{Conflict: "error"}).RunWrite(c.config.RDBSession)
	if err != nil || res.Errors > 0 {
		if err != nil {
			context["Error"] = "An error occurred while creating your account.. Please verify fields and try again"
		} else if res.Errors > 0 {
			context["Error"] = "Looks like that email address is already taken. If you've forgotten your password, please contact the administrator"
		}
		c.render.HTML(w, http.StatusOK, "register", context)
		return
	}

	context["Success"] = "Successfully registered email and password"
	c.render.HTML(w, http.StatusOK, "register", context)
}

// Endpoint for destroying CAS sessions (logging out)
func (c *CAS) HandleLogout(w http.ResponseWriter, req *http.Request) {
	context := map[string]string{ "CompanyName": c.config.CompanyName }

	// Exit early if the user is already logged in (in session)
	session, _ := c.cookieStore.Get(req, "casgo-session")
	if _,ok := session.Values["currentUserEmail"]; !ok {
		// Redirect if the person was never logged in
		http.Redirect(w, req, "/login", 301)
	}

	// Delete current user email (logging out user)
	delete(session.Values, "currentUserEmail")

	// Save the modified session
	err := session.Save(req, w)
	if err != nil {
		context["Error"] = "Failed to log out... Please contact your IT administrator"
		log.Fatal("Failed to remove logged in user from session:", err)
	} else {
		context["Success"] = "Successfully logged out"
	}

	c.render.HTML(w, http.StatusOK, "login", context)
}

// Endpoint for validating service tickets
func (c *CAS) HandleValidate(w http.ResponseWriter, req *http.Request) {

}

// Endpoint for validating service tickets (CAS 2.0)
func (c *CAS) HandleServiceValidate(w http.ResponseWriter, req *http.Request) {

}

// Endpoint for validating proxy tickets (CAS 2.0)
func (c *CAS) HandleProxyValidate(w http.ResponseWriter, req *http.Request) {

}

// Endpoint for handling proxy tickets (CAS 2.0)
func (c *CAS) HandleProxy(w http.ResponseWriter, req *http.Request) {

}
