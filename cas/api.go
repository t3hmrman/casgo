package cas

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

/*
 * CAS FrontendAPI implementation
 */

func NewCasgoFrontendAPI(c *CAS) (*FrontendAPI, error) {
	return &FrontendAPI{casServer: c}, nil
}

// Utility function to authenticate an API user, whether user is using a web-session or passed an API key
func authenticateAPIUser(api *FrontendAPI, req *http.Request) (*User, *CASServerError) {

	// Attempt to authenticate with HTTP session
	user, casErr := authenticateWithSession(api, req)
	if user != nil && casErr == nil {
		return user, nil
	}

	// Attempt to authenticate with API key and secret if present
	user, casErr = api.authenticateWithAPIKey(req)
	if user != nil && casErr == nil {
		return user, nil
	}

	// If all authentication methods fail, return error
	return nil, &FailedToAuthenticateUserError
}

// Authenticate user with session
func authenticateWithSession(api *FrontendAPI, req *http.Request) (*User, *CASServerError) {
	// Get the current session
	session, err := api.casServer.cookieStore.Get(req, "casgo-session")
	if err != nil {
		casErr := &FailedToRetrieveServicesError
		casErr.err = &err
		return nil, casErr
	}

	// Retrive current user from session
	user, ok := session.Values["currentUser"].(User)
	if !ok {
		casErr := &FailedToRetrieveInformationFromSessionError
		casErr.err = &err
		return nil, casErr
	}

	return &user, nil
}

func (api *FrontendAPI) authenticateWithAPIKey(req *http.Request) (*User, *CASServerError) {
	// Get the api key and secret
	apiKey := req.Header.Get("X-Api-Key")
	apiSecret := req.Header.Get("X-Api-Secret")
	if len(apiKey) == 0 || len(apiSecret) == 0 {
		return nil, &FailedToAuthenticateUserError
	}

	user, casErr := api.casServer.Db.FindUserByApiKeyAndSecret(apiKey, apiSecret)
	if casErr != nil {
		return nil, casErr
	}

	return user, nil
}

// Middleware for routes that require admin access
func (api *FrontendAPI) WrapAdminOnlyEndpoint(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// Get session and user
		requestingUser, casErr := authenticateAPIUser(api, req)
		if casErr != nil {
			api.casServer.render.JSON(w, casErr.HttpCode, map[string]string{
				"status":  "error",
				"message": casErr.Msg,
			})
			return
		}

		// Ensure user is admin
		if !requestingUser.IsAdmin {
			api.casServer.render.JSON(w, InsufficientPermissionsError.HttpCode, map[string]string{
				"status":  "error",
				"message": InsufficientPermissionsError.Msg,
			})
			return
		}

		// Run the actual handler
		handler(w, req)
	}
}

// Hook up API endpoints to given mux
func (api *FrontendAPI) HookupAPIEndpoints(m *mux.Router) {
	// Session information endpoints
	m.HandleFunc("/api/sessions/{userEmail}/services", api.listSessionUserServices).Methods("GET")
	m.HandleFunc("/api/sessions", api.SessionsHandler).Methods("GET")

	// Service endpoints
	m.HandleFunc("/api/users", api.GetUsers).Methods("GET")
	m.HandleFunc("/api/users", api.CreateUser).Methods("POST")
	m.HandleFunc("/api/users/{userEmail}", api.UpdateUser).Methods("PUT")
	m.HandleFunc("/api/users/{userEmail}", api.RemoveUser).Methods("DELETE")
	m.HandleFunc("/api/services", api.GetServices).Methods("GET")
	m.HandleFunc("/api/services", api.WrapAdminOnlyEndpoint(api.CreateService)).Methods("POST")
	m.HandleFunc("/api/services/{serviceName}", api.WrapAdminOnlyEndpoint(api.UpdateService)).Methods("PUT")
	m.HandleFunc("/api/services/{serviceName}", api.RemoveService).Methods("DELETE")
}

// Handle sessions endpoint
func (api *FrontendAPI) SessionsHandler(w http.ResponseWriter, req *http.Request) {
	user, casErr := authenticateAPIUser(api, req)
	if casErr != nil {
		api.casServer.render.JSON(w, casErr.HttpCode, map[string]string{
			"status":  "error",
			"message": casErr.Msg,
		})
		return
	}

	api.casServer.render.JSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   user,
	})
}

// Get the services for a logged in user
func (api *FrontendAPI) listSessionUserServices(w http.ResponseWriter, req *http.Request) {
	user, casErr := authenticateAPIUser(api, req)
	if casErr != nil {
		api.casServer.render.JSON(w, casErr.HttpCode, map[string]string{
			"status":  "error",
			"message": casErr.Msg,
		})
		return
	}

	// Quit early if the user is not an admin and is not the requested user
	routeVars := mux.Vars(req)
	routeUserEmail := routeVars["userEmail"]

	// Ensure non-admin user is not trying to lookup another users session information
	if !user.IsAdmin && user.Email != routeUserEmail {
		api.casServer.render.JSON(w, InsufficientPermissionsError.HttpCode, map[string]string{
			"status":  "error",
			"message": InsufficientPermissionsError.Msg,
		})
		return
	}

	// Return the user's services
	api.casServer.render.JSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   user.Services,
	})
}

///////////
// Users //
///////////

// Get list of users (admin only)
func (api *FrontendAPI) GetUsers(w http.ResponseWriter, req *http.Request) {
	// Get the current session and user
	user, casErr := authenticateAPIUser(api, req)
	if casErr != nil {
		api.casServer.render.JSON(w, casErr.HttpCode, map[string]string{
			"status":  "error",
			"message": casErr.Msg,
		})
		return
	}

	// Ensure user is admin
	if !user.IsAdmin {
		api.casServer.render.JSON(w, InsufficientPermissionsError.HttpCode, map[string]string{
			"status":  "error",
			"message": InsufficientPermissionsError.Msg,
		})
		return
	}

	// Grab list of all users
	users, casErr := api.casServer.Db.GetAllUsers()
	if casErr != nil {
		api.casServer.render.JSON(w, casErr.HttpCode, map[string]string{
			"status":  "error",
			"message": casErr.Msg,
		})
		return
	}

	api.casServer.render.JSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   users,
	})
}

// Create a new user
func (api *FrontendAPI) CreateUser(w http.ResponseWriter, req *http.Request) {
	// Get session and user
	requestingUser, casErr := authenticateAPIUser(api, req)
	if casErr != nil {
		api.casServer.render.JSON(w, casErr.HttpCode, map[string]string{
			"status":  "error",
			"message": casErr.Msg,
		})
		return
	}

	// Read JSON from request body
	var user User
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		api.casServer.render.JSON(w, InvalidUserError.HttpCode, map[string]string{
			"status":  "error",
			"message": InvalidUserError.Msg,
		})
		return
	}

	// Unmarshal JSON & build user from passed in data
	err = json.Unmarshal(reqBody, &user)
	if err != nil {
		api.casServer.render.JSON(w, FailedToParseJSONError.HttpCode, map[string]string{
			"status":  "error",
			"message": FailedToParseJSONError.Msg,
		})
		return
	}

	// Ensure user is valid
	if !user.IsValid() {
		api.casServer.render.JSON(w, InvalidUserError.HttpCode, map[string]string{
			"status":  "error",
			"message": InvalidUserError.Msg,
		})
		return
	}

	// Ensure user is admin before adding user
	if !requestingUser.IsAdmin {
		api.casServer.render.JSON(w, InsufficientPermissionsError.HttpCode, map[string]string{
			"status":  "error",
			"message": InsufficientPermissionsError.Msg,
		})
		return
	}

	// Attempt to add user
	newUser, casErr := api.casServer.Db.AddNewUser(user.Email, user.Password)
	if casErr != nil {
		api.casServer.render.JSON(w, casErr.HttpCode, map[string]string{
			"status":  "error",
			"message": casErr.Msg,
		})
		return
	}

	api.casServer.render.JSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   newUser,
	})
}

// Remove a user
// Returns the removed user's email
func (api *FrontendAPI) RemoveUser(w http.ResponseWriter, req *http.Request) {
	// Get session and user
	requestingUser, casErr := authenticateAPIUser(api, req)
	if casErr != nil {
		api.casServer.render.JSON(w, casErr.HttpCode, map[string]string{
			"status":  "error",
			"message": casErr.Msg,
		})
		return
	}

	// Ensure user is admin
	if !requestingUser.IsAdmin {
		api.casServer.render.JSON(w, InsufficientPermissionsError.HttpCode, map[string]string{
			"status":  "error",
			"message": InsufficientPermissionsError.Msg,
		})
		return
	}

	// Get passed in user name
	routeVars := mux.Vars(req)
	userEmail := routeVars["userEmail"]

	casErr = api.casServer.Db.RemoveUserByEmail(userEmail)
	if casErr != nil {
		api.casServer.render.JSON(w, casErr.HttpCode, map[string]string{
			"status":  "error",
			"message": casErr.Msg,
		})
	}

	api.casServer.render.JSON(w, http.StatusOK, map[string]string{
		"status": "success",
		"data":   userEmail,
	})
}

// Update an existing user
// Returns the modified user
func (api *FrontendAPI) UpdateUser(w http.ResponseWriter, req *http.Request) {
	// Get session and user
	requestingUser, casErr := authenticateAPIUser(api, req)
	if casErr != nil {
		api.casServer.render.JSON(w, casErr.HttpCode, map[string]string{
			"status":  "error",
			"message": casErr.Msg,
		})
		return
	}

	// Ensure user is admin
	if !requestingUser.IsAdmin {
		api.casServer.render.JSON(w, InsufficientPermissionsError.HttpCode, map[string]string{
			"status":  "error",
			"message": InsufficientPermissionsError.Msg,
		})
		return
	}

	// Get passed in user name
	routeVars := mux.Vars(req)
	userEmail := routeVars["userEmail"]

	// Read JSON from request body
	var user User
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		api.casServer.render.JSON(w, InvalidUserError.HttpCode, map[string]string{
			"status":  "error",
			"message": InvalidUserError.Msg,
		})
		return
	}

	// Unmarshal JSON & build user from passed in data
	err = json.Unmarshal(reqBody, &user)
	if err != nil {
		api.casServer.render.JSON(w, FailedToParseJSONError.HttpCode, map[string]string{
			"status":  "error",
			"message": FailedToParseJSONError.Msg,
		})
		return
	}

	// Ensure user is valid
	if !user.IsValidUpdate() || userEmail != user.Email {
		api.casServer.render.JSON(w, InvalidUserError.HttpCode, map[string]string{
			"status":  "error",
			"message": InvalidUserError.Msg,
		})
		return
	}

	// Attempt to update the user
	casErr = api.casServer.Db.UpdateUser(&user)
	if casErr != nil {
		api.casServer.render.JSON(w, casErr.HttpCode, map[string]string{
			"status":  "error",
			"message": casErr.Msg,
		})
		return
	}

	api.casServer.render.JSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   user,
	})
}

//////////////
// Services //
//////////////

// Get list of services (admin only)
func (api *FrontendAPI) GetServices(w http.ResponseWriter, req *http.Request) {
	// Get the current session and user
	user, casErr := authenticateAPIUser(api, req)
	if casErr != nil {
		api.casServer.render.JSON(w, casErr.HttpCode, map[string]string{
			"status":  "error",
			"message": casErr.Msg,
		})
		return
	}

	// Ensure user is admin
	if !user.IsAdmin {
		api.casServer.render.JSON(w, InsufficientPermissionsError.HttpCode, map[string]string{
			"status":  "error",
			"message": InsufficientPermissionsError.Msg,
		})
		return
	}

	// Grab list of all services
	services, casErr := api.casServer.Db.GetAllServices()
	if casErr != nil {
		api.casServer.render.JSON(w, casErr.HttpCode, map[string]string{
			"status":  "error",
			"message": casErr.Msg,
		})
		return
	}

	api.casServer.render.JSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   services,
	})
}

// Create a new service
func (api *FrontendAPI) CreateService(w http.ResponseWriter, req *http.Request) {
	// Read JSON from request body
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		api.casServer.render.JSON(w, FailedToParseJSONError.HttpCode, map[string]string{
			"status":  "error",
			"message": FailedToParseJSONError.Msg,
		})
		return
	}

	// Unmarshal JSON & build service from passed in data
	var service CASService
	err = json.Unmarshal(reqBody, &service)
	if err != nil {
		api.casServer.render.JSON(w, InvalidServiceError.HttpCode, map[string]string{
			"status":  "error",
			"message": InvalidServiceError.Msg,
		})
		return
	}

	// Ensure service is valid
	if !service.IsValid() {
		api.casServer.render.JSON(w, InvalidServiceError.HttpCode, map[string]string{
			"status":  "error",
			"message": InvalidServiceError.Msg,
		})
		return
	}

	// Attempt to add service
	casErr := api.casServer.Db.AddNewService(&service)
	if casErr != nil {
		api.casServer.render.JSON(w, casErr.HttpCode, map[string]string{
			"status":  "error",
			"message": casErr.Msg,
		})
		return
	}

	api.casServer.render.JSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   service,
	})
}

// Remove a service
// Returns the removed service's name
func (api *FrontendAPI) RemoveService(w http.ResponseWriter, req *http.Request) {
	// Get session and user
	user, casErr := authenticateAPIUser(api, req)
	if casErr != nil {
		api.casServer.render.JSON(w, casErr.HttpCode, map[string]string{
			"status":  "error",
			"message": casErr.Msg,
		})
		return
	}

	// Ensure user is admin
	if !user.IsAdmin {
		api.casServer.render.JSON(w, InsufficientPermissionsError.HttpCode, map[string]string{
			"status":  "error",
			"message": InsufficientPermissionsError.Msg,
		})
		return
	}

	// Get passed in service name
	routeVars := mux.Vars(req)
	serviceName := routeVars["serviceName"]

	casErr = api.casServer.Db.RemoveServiceByName(serviceName)
	if casErr != nil {
		api.casServer.render.JSON(w, casErr.HttpCode, map[string]string{
			"status":  "error",
			"message": casErr.Msg,
		})
	}

	api.casServer.render.JSON(w, http.StatusOK, map[string]string{
		"status": "success",
		"data":   serviceName,
	})
}

// Update an existing service
// Returns the modified service
func (api *FrontendAPI) UpdateService(w http.ResponseWriter, req *http.Request) {
	// Get passed in service name
	routeVars := mux.Vars(req)
	serviceName := routeVars["serviceName"]

	// Read JSON from request body
	var service CASService
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		api.casServer.render.JSON(w, InvalidServiceError.HttpCode, map[string]string{
			"status":  "error",
			"message": InvalidServiceError.Msg,
		})
		return
	}

	// Unmarshal JSON & build service from passed in data
	err = json.Unmarshal(reqBody, &service)
	if err != nil {
		api.casServer.render.JSON(w, FailedToParseJSONError.HttpCode, map[string]string{
			"status":  "error",
			"message": FailedToParseJSONError.Msg,
		})
		return
	}

	// Set the service name if it wasn't on the incoming request's object
	if len(service.Name) == 0 {
		service.Name = serviceName
	}

	// Ensure service is valid
	if !service.IsValidUpdate() || serviceName != service.Name {
		api.casServer.render.JSON(w, InvalidServiceError.HttpCode, map[string]string{
			"status":  "error",
			"message": InvalidServiceError.Msg,
		})
		return
	}

	// Attempt to update the service
	casErr := api.casServer.Db.UpdateService(&service)
	if casErr != nil {
		api.casServer.render.JSON(w, casErr.HttpCode, map[string]string{
			"status":  "error",
			"message": casErr.Msg,
		})
		return
	}

	api.casServer.render.JSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   service,
	})
}
