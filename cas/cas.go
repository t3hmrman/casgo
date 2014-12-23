package cas

import (
	"code.google.com/p/go.crypto/bcrypt"
	"encoding/gob"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/*
 * CAS server implementation
 */

func NewCASServer(userConfigOverrides map[string]string) (*CAS, error) {
	// Create and initialize the CAS server
	cas := &CAS{
		Config:      nil,
		render:      nil,
		cookieStore: nil,
		ServeMux:    nil,
	}

	// Create configuration with user overrides provided
	config, err := NewCASServerConfig(userConfigOverrides)
	if err != nil {
		log.Fatalf("Failed to create new CAS server configuration, err: %v", err)
	}
	cas.Config = config

	// Setup rendering function
	render := render.New(render.Options{
		Directory: cas.Config["templatesDirectory"],
		Layout:    "layout",
	})
	cas.render = render

	// Cookie store setup
	cookieStore := sessions.NewCookieStore([]byte(cas.Config["cookieSecret"]))
	cookieStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	cas.cookieStore = cookieStore

	// Register types for encoding/decoding
	gob.Register([]CASService{})

	cas.init()
	cas.setLogLevel(cas.Config["logLevel"])
	return cas, nil
}

func (c *CAS) setLogLevel(lvl string) {
	switch lvl {
	case "WARN":
		c.LogLevel = WARN
	case "DEBUG":
		c.LogLevel = DEBUG
	case "INFO":
		c.LogLevel = INFO
	default:
		c.LogLevel = WARN
	}
}

func (c *CAS) init() {
	// Override config with ENV variables
	c.Config = overrideConfigWithEnv(c.Config)

	// Setup database adapter
	Db, err := NewRethinkDBAdapter(c)
	if err != nil {
		log.Fatal("Failed to setup database adapter", err)
	}
	c.Db = Db

	// Setup the internal HTTP Server
	c.server = &http.Server{
		Addr: c.GetAddr(),
	}

	// Setup handlers
	serveMux := mux.NewRouter()

	// Front end endpoints
	serveMux.HandleFunc("/login", c.HandleLogin)
	serveMux.HandleFunc("/logout", c.HandleLogout)
	serveMux.HandleFunc("/register", c.HandleRegister)

	// User-accessible API endpoints
	serveMux.HandleFunc("/api/sessions/{userEmail}/services", c.listSessionUserServices).Methods("GET")
	serveMux.HandleFunc("/api/sessions", c.SessionHandler).Methods("GET")

	// Admin-only endpoints
	serveMux.HandleFunc("/api/users", c.UsersHandler)
	serveMux.HandleFunc("/api/services", c.ServicesHandler)

	// CAS-specific endpoints
	serveMux.HandleFunc("/validate", c.HandleValidate)
	serveMux.HandleFunc("/serviceValidate", c.HandleServiceValidate)
	serveMux.HandleFunc("/proxyValidate", c.HandleProxyValidate)
	serveMux.HandleFunc("/proxy", c.HandleProxy)

	// Static file serving
	serveMux.PathPrefix("/public/").Handler(http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))
	serveMux.HandleFunc("/", c.HandleIndex)

	c.ServeMux = serveMux
	c.server.Handler = c.ServeMux
}

// Set up the underlying database
func (c *CAS) SetupDb() {
	c.Db.Setup()
}

// Teardown the underlying database
func (c *CAS) TeardownDb() {
	c.Db.Teardown()
}

// Start the CAS server
func (c *CAS) Start() {
	// Start server
	log.Fatal(c.server.ListenAndServe())
}

// Get the address of the server based on server configuration
func (c *CAS) GetAddr() string {
	return c.Config["host"] + ":" + c.Config["port"]
}

// (Optional) Handles Index route
func (c *CAS) HandleIndex(w http.ResponseWriter, req *http.Request) {

	// Attempt to retrieve user session and populate template context
	session, _ := c.cookieStore.Get(req, "casgo-session")
	templateContext := c.augmentTemplateContext(map[string]interface{}{}, session)

	// Exit early (and show landing page) if not user not logged in (in session)
	if _, ok := templateContext["userEmail"]; !ok {
		c.render.HTML(w, http.StatusOK, "landing", templateContext)
		return
	}

	c.render.HTML(w, http.StatusOK, "index", templateContext)
}

// Augment information in given context with information from given session
// Will overwrite any fields that are already filled
func (c *CAS) augmentTemplateContext(context map[string]interface{}, session *sessions.Session) map[string]interface{} {
	context["CompanyName"] = c.Config["companyName"]

	// Add information from session
	if session != nil {
		if userEmail, ok := session.Values["userEmail"]; ok {
			context["userEmail"] = userEmail.(string)
		}

		if isAdmin, ok := session.Values["userIsAdmin"]; ok {
			context["userIsAdmin"] = isAdmin
		}
	}

	return context
}

// Handle logins (functions as both a credential acceptor and requestor)
func (c *CAS) HandleLogin(w http.ResponseWriter, req *http.Request) {
	// Generate context
	context := map[string]interface{}{"CompanyName": c.Config["companyName"]}

	// Trim and lightly pre-process/validate service
	serviceUrl := strings.TrimSpace(strings.ToLower(req.FormValue("service")))
	gateway := strings.TrimSpace(strings.ToLower(req.FormValue("gateway")))
	renew := strings.TrimSpace(strings.ToLower(req.FormValue("renew")))
	method := strings.TrimSpace(strings.ToLower(req.FormValue("method")))

	// In the case login is being used as an acceptor
	email := strings.TrimSpace(strings.ToLower(req.FormValue("email")))
	password := strings.TrimSpace(strings.ToLower(req.FormValue("password")))

	// Handle service being not set early
	var casService *CASService
	if len(serviceUrl) > 0 {
		foundService, err := c.Db.FindServiceByUrl(serviceUrl)
		if err != nil {
			context["Error"] = "Failed to find matching service with URL [" + serviceUrl + "]."
			c.render.HTML(w, http.StatusNotFound, "login", context)
			return
		}
		casService = foundService
	}

	// Pass method along in context if specified & valid
	if method == "post" || method == "get" {
		context["Method"] = method
	}

	// Both gateway and  cannot be set -- Croak here? Maybe also take renew over gateway (as docs suggest?)
	if gateway == "true" && renew == "true" {
		context["Error"] = "Invalid Request: Both gateway and renew options specified"
		c.render.HTML(w, http.StatusBadRequest, "login", context)
		return
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
		if _, ok := session.Values["userEmail"]; ok {

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
		returnedUser, err := c.validateUserCredentials(email, password)
		if err != nil {
			// In the case of an error, redirect to the service with no ticket
			if casService == nil {
				context["Error"] = err.msg
				c.render.HTML(w, err.httpCode, "login", context)
			} else {
				http.Redirect(w, req, casService.Url, 401)
			}
			return
		}

		// Save session since non-interactive auth succeeded
		_, err = c.saveUserInfoInSession(w, req, "casgo-session", returnedUser)
		if err != nil {
			log.Fatal("Failed to save session!")
		}

		if casService == nil {
			// If service is not set, render login with context
			c.render.HTML(w, err.httpCode, "login", context)
		} else {
			// Create a new ticket
			ticket := &CASTicket{
				UserEmail:      returnedUser.Email,
				UserAttributes: returnedUser.Attributes,
				WasSSO:         false,
			}

			// If service is set, redirect
			ticket, err := c.Db.AddTicketForService(ticket, casService)
			if err != nil {
				http.Error(w, "Failed to create new authentication ticket. Please contact administrator if problem persists.", 500)
				return
			}
			http.Redirect(w, req, serviceUrl+"?ticket="+ticket.Id, 302)
			return
		}

	} // /if gateway == true

	// Trim and lightly pre-process/validate email/password
	if email == "" || password == "" {
		c.render.HTML(w, http.StatusOK, "login", context)
		return
	}

	// Find user, and attempt to validate provided credentials
	returnedUser, err := c.validateUserCredentials(email, password)
	if err != nil {
		context["Error"] = err.msg
		c.render.HTML(w, err.httpCode, "login", context)
		return
	}

	// Save session in cookies
	session, err := c.saveUserInfoInSession(w, req, "casgo-session", returnedUser)
	if err != nil {
		log.Fatal("Failed to save session, err:", err)
	}

	// Update context with session
	if err == nil {
		c.augmentTemplateContext(context, session)
	}

	// If the user was already logged in service was provided, create a new ticket (with SSO set to true) and redirect
	// Otherwise render login page
	if serviceUrl != "" {

		ssoTicket := &CASTicket{
			UserEmail:      returnedUser.Email,
			UserAttributes: returnedUser.Attributes,
			WasSSO:         true,
		}

		// Get ticket for the service
		ticket, err := c.Db.AddTicketForService(ssoTicket, casService)
		if err != nil {
			http.Error(w, "Failed to create new authentication ticket. Please contact administrator if problem persists.", 500)
			return
		}

		http.Redirect(w, req, serviceUrl+"?ticket="+ticket.Id, 302)
		return
	} else {

		context["Success"] = "Successful log in! Redirecting to services page..."
		c.render.HTML(w, http.StatusOK, "login", context)
	}
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
func (c *CAS) saveUserInfoInSession(w http.ResponseWriter, req *http.Request, sessionName string, user *User) (*sessions.Session, *CASServerError) {
	// Save session in cookies
	session, _ := c.cookieStore.Get(req, sessionName)

	// Save user information onto session
	session.Values["userEmail"] = user.Email
	session.Values["userServices"] = user.Services
	session.Values["userIsAdmin"] = user.IsAdmin

	// Save the session
	sessionSaveErr := session.Save(req, w)
	if sessionSaveErr != nil {
		log.Fatal("Failed to save logged in user to session:", sessionSaveErr)
		return nil, &FailedToSaveSessionError
	}

	return session, nil
}

// Validate user credentials
// Returns a valid user object if validation succeeds
func (c *CAS) validateUserCredentials(email string, password string) (*User, *CASServerError) {

	// TODO get the user from the current database adapter
	returnedUser, err := c.Db.FindUserByEmail(email)
	if err != nil {
		return nil, &FailedToFindUserError
	}

	// Use default authentication typeDepending on the authentication type
	switch c.Config["authMethod"] {
	case "password":
		// Check hash
		err := bcrypt.CompareHashAndPassword([]byte(returnedUser.Password), []byte(password))
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
	context := map[string]interface{}{"CompanyName": c.Config["companyName"]}

	// Show login page if credentials are not provided, attempt login otherwise
	email := strings.TrimSpace(strings.ToLower(req.FormValue("email")))
	password := strings.TrimSpace(strings.ToLower(req.FormValue("password")))

	// Exit early if email/password are empty
	if email == "" || password == "" {
		c.render.HTML(w, http.StatusOK, "register", context)
		return
	}

	// Generate hashed password
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10) // Default cost
	if err != nil {
		context["Error"] = "Registration failed... Please contact server administrator"
		c.render.HTML(w, http.StatusInternalServerError, "register", context)
		return
	}

	// Create new user object
	_, casErr := c.Db.AddNewUser(email, string(encryptedPassword))
	if casErr != nil {
		context["Error"] = casErr.msg
		c.render.HTML(w, http.StatusBadRequest, "register", context)
		return
	}

	context["Success"] = "Registration successful!"
	c.render.HTML(w, http.StatusOK, "register", context)
}

// Endpoint for destroying CAS sessions (logging out)
func (c *CAS) HandleLogout(w http.ResponseWriter, req *http.Request) {
	context := map[string]interface{}{"CompanyName": c.Config["companyName"]}

	// Get the user's session
	session, _ := c.cookieStore.Get(req, "casgo-session")

	serviceUrl := strings.TrimSpace(strings.ToLower(req.FormValue("service")))

	// Get the CASService for this service URL
	var casService *CASService
	if len(serviceUrl) > 0 {
		returnedService, err := c.Db.FindServiceByUrl(serviceUrl)
		if err != nil {
			context["Error"] = "Failed to find matching service with URL [" + serviceUrl + "]."
			c.render.HTML(w, http.StatusNotFound, "login", context)
			return
		}
		casService = returnedService
	}

	// Exit early if the user is not already logged in (in session), otherwise get their email
	if _, ok := session.Values["userEmail"]; !ok {
		// Redirect if the person was never logged in
		http.Redirect(w, req, "/login", 401)
		return
	}
	userEmail := session.Values["userEmail"]

	// If service was specified, Delete any ticket granting tickets that belong to the user
	err := c.Db.RemoveTicketsForUserWithService(userEmail.(string), casService)
	if err != nil {
		log.Printf("Failed to remove ticket for user %s", userEmail.(string))
		http.Redirect(w, req, "/login", 500)
		return
	}

	// Remove current user information from session
	err = c.removeCurrentUserFromSession(w, req, session)
	if err != nil {
		context["Error"] = "Failed to log out... Please contact your IT administrator"
		c.render.HTML(w, err.httpCode, "login", context)
		return
	}

	context["Success"] = "Successfully logged out"
	c.render.HTML(w, http.StatusOK, "login", context)
}

// Remove all current user information from the session object
func (c *CAS) removeCurrentUserFromSession(w http.ResponseWriter, req *http.Request, session *sessions.Session) *CASServerError {
	// Delete current user from session
	delete(session.Values, "userEmail")

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
	casService, err := c.Db.FindServiceByUrl(serviceUrl)
	if err != nil {
		log.Printf("Failed to find matching service with URL [%s]", serviceUrl)
		c.render.JSON(w, http.StatusOK, map[string]string{
			"status":  "error",
			"code":    strconv.Itoa(*&FailedToFindServiceError.casErrCode),
			"message": *&FailedToFindServiceError.msg,
		})
		return
	}

	// Look up ticket
	casTicket, err := c.Db.FindTicketByIdForService(ticket, casService)
	if err != nil {
		log.Printf("Failed to find matching ticket", casService.Url)
		c.render.JSON(w, http.StatusOK, map[string]string{
			"status":  "error",
			"code":    strconv.Itoa(*&FailedToFindTicketError.casErrCode),
			"message": *&FailedToFindTicketError.msg,
		})
		return
	}

	// If renew is specified, validation only works if the login is fresh (not from a single sign on session)
	if renew == "true" && casTicket.WasSSO {
		c.render.JSON(w, http.StatusOK, map[string]string{
			"status":  "error",
			"code":    strconv.Itoa(*&SSOAuthenticatedUserRenewError.casErrCode),
			"message": *&SSOAuthenticatedUserRenewError.msg,
		})
		return
	}

	// Successfully validated user send user information along
	c.render.JSON(w, http.StatusOK, map[string]interface{}{
		"status":         "success",
		"message":        "Successfully authenticated user",
		"userEmail":      casTicket.UserEmail,
		"userAttributes": casTicket.UserAttributes,
	})
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
