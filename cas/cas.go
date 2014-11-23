package cas

/**
 * CAS protocol API
 */

import (
	"code.google.com/p/go.crypto/bcrypt"
	r "github.com/dancannon/gorethink"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
	"log"
	"net/http"
	"strings"
	"strconv"
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
	Config      map[string]string
	RDBSession  *r.Session
	render      *render.Render
	cookieStore *sessions.CookieStore
}

func NewCASServer(config map[string]string) *CAS {
	// Setup rendering function
	render := render.New(render.Options{Directory: config["TemplatesDirectory"]})

	// Cookie store setup
	cookieStore := sessions.NewCookieStore([]byte(config["CookieSecret"]))
	cookieStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	// Database setup
	dbSession, err := r.Connect(r.ConnectOpts{
		Address:  config["DBHost"],
		Database: config["DBName"],
	})
	if err != nil {
		log.Fatalln(err.Error())
	} 

	return &CAS{config, dbSession, render, cookieStore}
}

func (c *CAS) init() {
	// override config with ENV variables
	c.overrideConfigWithEnv()
}

func (c *CAS) GetAddr() string {
	return c.Config["Host"] + ":" + c.Config["Port"]
}

// (Optional) Handles Index route
func (c *CAS) HandleIndex(w http.ResponseWriter, req *http.Request) {
	c.render.HTML(w, http.StatusOK, "index", map[string]string{"CompanyName": c.Config["CompanyName"]})
}

// Credential acceptor endpoint (requestor is Handled in main)
func (c *CAS) HandleLogin(w http.ResponseWriter, req *http.Request) {
	// Generate context
	context := map[string]string{"CompanyName": c.Config["CompanyName"]}

	// Trim and lightly pre-process/validate service
	service := strings.TrimSpace(strings.ToLower(req.FormValue("service")))
	gateway := strings.TrimSpace(strings.ToLower(req.FormValue("gateway")))
	renew := strings.TrimSpace(strings.ToLower(req.FormValue("renew")))
	method := strings.TrimSpace(strings.ToLower(req.FormValue("method")))

	// Handle service being not set early
	casService, err := c.getService(service)
	if err != nil {
		context["Error"] = "Failed to find matching service with URL [" + service + "]."
		c.render.HTML(w, http.StatusNotFound, "login", context)
	}

	// Pass method along in context if specified & valid
	if method == "post" || method == "get" {
		context["Method"] = method
	}

	// Both gateway and  cannot be set -- Croak here? Maybe also take renew over gateway (as docs suggest?)
	if gateway == "true" && renew == "true" {
		context["Error"] = "Invalid Request: Both gateway and renew options specified"
		c.render.HTML(w, http.StatusBadRequest, "login", context)
	}

	if renew == "true" {

		// If renew is set, automatic sign on is disabled, user must present credentials regardless of whether a sign on session exists
		// Renew takes priority over gateway
		c.render.HTML(w, http.StatusOK, "login", context)
		return

	} else if gateway == "true" {

		// If gateway is set, CAS will try to use previous session or authenticate with non-interactive means (ex. LDAP)
		// If no CAS session and no non-interactive means, then redirect with no ticket parameter to service URL

		// Finish early if the user is already logged in (has session)
		session, _ := c.cookieStore.Get(req, "casgo-session")
		if _, ok := session.Values["currentUserEmail"]; ok {

			// If session is not set and gateway is set, behavior is undefined, act as if nothing was given, let user know they are logged in
			// Otherwiser make new ticket and properly redirect to service
			if casService == nil {
				context["Success"] = "User already logged in..."
				c.render.HTML(w, http.StatusOK, "login", context)
			} else {
				_, err := c.makeNewTicketAndRedirect(w, req, casService)
				if err != nil {
					http.Error(w, "Failed to create new authentication ticket. Please contact administrator if problem persists.", 500)
				}
			}

			return
		}

		// Attempt non-interactive authentication
		returnedUser, err := c.validateUserCredentials("", "")
		if err != nil {
			// In the case of an error, redirect to the service with no ticket
			if casService == nil {
				context["Error"] = err.msg
				c.render.HTML(w, err.http_code, "login", context)
			} else {
				http.Redirect(w, req, casService.Url, 401)
			}
			return
		}

		// Save session since non-interactive auth succeeded
		_, err = c.saveUserEmailInSession(w, req, "casgo-session", returnedUser.Email)
		if err != nil {
			log.Fatal("Failed to save session!")
		}

		if casService == nil {
			// If service is not set, render login with context
			c.render.HTML(w, err.http_code, "login", context)
		} else {
			// If service is set, redirect
			ticket, err := c.makeNewTicketForService(casService)
			if err != nil {
				http.Error(w, "Failed to create new authentication ticket. Please contact administrator if problem persists.", 500)
				return
			}
			http.Redirect(w, req, service+"?ticket="+ticket, 302)
			return
		}

	} // /if gateway == true

	// Trim and lightly pre-process/validate email/password
	email := strings.TrimSpace(strings.ToLower(req.FormValue("email")))
	password := strings.TrimSpace(strings.ToLower(req.FormValue("password")))
	if email == "" || password == "" {
		c.render.HTML(w, http.StatusOK, "login", context)
		return
	}

	// Find user, and attempt to validate provided credentials
	returnedUser, err := c.validateUserCredentials(email, password)
	if err != nil {
		context["Error"] = err.msg
		c.render.HTML(w, err.http_code, "login", context)
		return
	}

	// Save session in cookies
	c.saveUserEmailInSession(w, req, "casgo-session", returnedUser.Email)
	if err != nil {
		log.Fatal("Failed to save session, err:", err)
	}

	// If the user has logged in and service was provided, redirect
	// Otherwise render login page
	if service != "" {
		// Get ticket for the service
		ticket, err := c.makeNewTicketForService(casService)
		if err != nil {
			http.Error(w, "Failed to create new authentication ticket. Please contact administrator if problem persists.", 500)
			return
		}
		http.Redirect(w, req, service+"?ticket="+ticket, 302)
		return
	} else {
		context["Success"] = "Successful log in! Redirecting to services page..."
		context["currentUserEmail"] = returnedUser.Email
		c.render.HTML(w, http.StatusOK, "login", context)
	}
}

// Get the service that belongs to the cas
func (c *CAS) getService(serviceName string) (*CASService, *CASServerError) {
	return &CASService{serviceName, "nobody@nowhere.net"}, nil
}

// Make a new ticket for a service
func (c *CAS) makeNewTicketForService(service *CASService) (string, *CASServerError) {
	return "123456", nil
}

func (c *CAS) makeNewTicketAndRedirect(w http.ResponseWriter, req *http.Request, service *CASService) (bool, *CASServerError) {
	// If service is set, redirect
	ticket, err := c.makeNewTicketForService(service)
	if err != nil {
		http.Error(w, "Failed to create new authentication ticket. Please contact administrator if problem persists.", 500)
		return false, &FailedToCreateNewAuthTicketError
	}
	redirectUrl := service.Url + "?ticket=" + ticket
	http.Redirect(w, req, redirectUrl, 302)
	return true, nil
}

// Save session in cookiestore
func (c *CAS) saveUserEmailInSession(w http.ResponseWriter, req *http.Request, sessionName string, email string) (bool, *CASServerError) {
	// Save session in cookies
	session, _ := c.cookieStore.Get(req, sessionName)
	session.Values["currentUserEmail"] = email
	sessionSaveErr := session.Save(req, w)
	if sessionSaveErr != nil {
		log.Fatal("Failed to save logged in user to session:", sessionSaveErr)
		return false, &FailedToSaveSessionError
	}
	return true, nil
}

// Validate user credentials
// Returns a valid user object if validation succeeds
func (c *CAS) validateUserCredentials(email string, password string) (*User, *CASServerError) {

	// Find the user
	cursor, err := r.Db(c.Config["DBName"]).Table("users").Get(email).Run(c.RDBSession)
	if err != nil {
		return nil, &InvalidEmailAddressError
	}

	// Get the user from the returned cursor
	var returnedUser *User
	err = cursor.One(&returnedUser)
	if err != nil {
		return nil, &InvalidEmailAddressError
	}

	// Use default authentication typeDepending on the authentication type
	switch c.Config["DefaultAuthMethod"] {
	case "password":
		// Check hash
		err = bcrypt.CompareHashAndPassword([]byte(returnedUser.Password), []byte(password))
		if err != nil {
			return nil, &InvalidCredentialsError
		}
		break
	default:
		return nil, &AuthMethodNotSupportedError
		break
	}

	// Successful validation
	return returnedUser, nil
}

// Endpoint for registering new users
func (c *CAS) HandleRegister(w http.ResponseWriter, req *http.Request) {
	context := map[string]string{"CompanyName": c.Config["CompanyName"]}

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

	res, err := r.Db(c.Config["DBName"]).Table("users").Insert(newUser, r.InsertOpts{Conflict: "error"}).RunWrite(c.RDBSession)
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
	context := map[string]string{"CompanyName": c.Config["CompanyName"]}

	// Get the user's session
	session, _ := c.cookieStore.Get(req, "casgo-session")

	service := strings.TrimSpace(strings.ToLower(req.FormValue("service")))
	// Get the CASService for this service URL
	casService, err := c.getService(service)
	if err != nil {
		context["Error"] = "Failed to find matching service with URL [" + service + "]."
		c.render.HTML(w, http.StatusNotFound, "login", context)
	}

	// Attempt to find the session for the user, exit early if there is no session
	userEmail := session.Values["currentUserEmail"]
	if userEmail == nil {
		http.Redirect(w, req, "/login", 301)
		return
	}

	// If service was specified, Delete any ticket granting tickets that belong to the user
	c.removeTicketsForUser(userEmail.(string), casService)
	if err != nil {
		log.Printf("Failed to remove ticket for user %s", userEmail.(string))
	}

	// Exit early if the user is not already logged in (in session)
	if _, ok := session.Values["currentUserEmail"]; !ok {
		// Redirect if the person was never logged in
		http.Redirect(w, req, "/login", 301)
		return
	}

	// Remove current user information from session
	c.removeCurrentUserFromSession(w, req, session)
	if err != nil {
		context["Error"] = "Failed to log out... Please contact your IT administrator"
		c.render.HTML(w, err.http_code, "login", context)
		return
	}

	context["Success"] = "Successfully logged out"
	c.render.HTML(w, http.StatusOK, "login", context)
}

// Remove the ticket granting tickets for a given user on a given service
// If service is nil, for all services.
func (c *CAS) removeTicketsForUser(userEmail string, service *CASService) {
	// Something
}

// Remove all current user information from the session object
func (c *CAS) removeCurrentUserFromSession(w http.ResponseWriter, req *http.Request, session *sessions.Session) *CASServerError {
	// Delete current user from session
	delete(session.Values, "currentUserEmail")

	// Save the modified session
	err := session.Save(req, w)
	if err != nil {
		return &FailedToDeleteSessionError
	}

	return nil
}

// Endpoint for validating service tickets
func (c *CAS) HandleValidate(w http.ResponseWriter, req *http.Request) {

	// Grab important request parameters
	serviceUrl := strings.TrimSpace(strings.ToLower(req.FormValue("service")))
	ticket := strings.TrimSpace(strings.ToLower(req.FormValue("ticket")))
	renew := strings.TrimSpace(strings.ToLower(req.FormValue("renew")))

	// Get the CASService for the given service URL
	casService, err := c.getService(serviceUrl)
	if err != nil {
		log.Printf("Failed to find matching service with URL [%s]", serviceUrl)
		c.render.JSON(w, http.StatusOK, map[string]string{
			"status": "error",
			"code": strconv.Itoa(*&FailedToFindServiceError.err_code),
			"message": *&FailedToFindServiceError.msg,
		})
		return
	}

	// Look up ticket
	casTicket, err := c.getTicketForService(casService, ticket)
	if err != nil {
		log.Printf("Failed to find matching ticket", casService.Url)
		c.render.JSON(w, http.StatusOK, map[string]string{
			"status": "error",
			"code": strconv.Itoa(*&FailedToFindTicketError.err_code),
			"message": *&FailedToFindTicketError.msg,
		})
		return
	}

	// If renew is specified, validation only works if the login is fresh (not from a single sign on session)
	if renew == "true" && casTicket.wasFromSSOSession {
		c.render.JSON(w, http.StatusOK, map[string]string{
			"status": "error",
			"code": strconv.Itoa(*&SSOAuthenticatedUserRenewError.err_code),
			"message": *&SSOAuthenticatedUserRenewError.msg,
		})
		return
	}

	// Successfully validated user send user information along
	c.render.JSON(w, http.StatusOK, map[string]string{
		"status": "success",
		"message": "Successfully authenticated user",
		"username": "",
	})
}

// Get the ticket for a given service
func (c *CAS) getTicketForService(service *CASService, ticket string) (*CASTicket, *CASServerError) {
	return &CASTicket{"ABC", false}, nil
}

// Endpoint for validating service tickets for possible proxies (CAS 2.0)
func (c *CAS) HandleServiceValidate(w http.ResponseWriter, req *http.Request) {
	log.Print("Attempt to use /serviceValidate, feature not supported yet")
	c.render.JSON(w, http.StatusOK, map[string]string{"error": *&UnsupportedFeatureError.msg})
}

// Endpoint for validating proxy tickets (CAS 2.0)
func (c *CAS) HandleProxyValidate(w http.ResponseWriter, req *http.Request) {
	log.Print("Attempt to use /proxyValidate, feature not supported yet")
	c.render.JSON(w, http.StatusOK, map[string]string{"error": *&UnsupportedFeatureError.msg})
}

// Endpoint for handling proxy tickets (CAS 2.0)
func (c *CAS) HandleProxy(w http.ResponseWriter, req *http.Request) {
	log.Print("Attempt to use /proxy, feature not supported yet")
	c.render.JSON(w, http.StatusOK, map[string]string{"error": *&UnsupportedFeatureError.msg})
}
