package cas

import (
	"github.com/gorilla/mux"
	"net/http"
	"strings"
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

// Hook up API endpoints to given mux
func (api *FrontendAPI) HookupAPIEndpoints(m *mux.Router) {
	// Session information endpoints
	m.HandleFunc("/api/sessions/{userEmail}/services", api.listSessionUserServices).Methods("GET")
	m.HandleFunc("/api/sessions", api.SessionsHandler).Methods("GET")

	// Service endpoints
	m.HandleFunc("/api/services", api.GetServices).Methods("GET")
	m.HandleFunc("/api/services", api.CreateService).Methods("POST")
	m.HandleFunc("/api/services/{serviceName}", api.UpdateService).Methods("PUT")
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
		api.casServer.render.JSON(w, http.StatusUnauthorized, map[string]string{
			"status":  "error",
			"message": "Insufficient permissions",
		})
		return
	}

	// Return the user's services
	api.casServer.render.JSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   user.Services,
	})
}

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
		api.casServer.render.JSON(w, http.StatusUnauthorized, map[string]string{
			"status":  "error",
			"message": "Insufficient permissions.",
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
	// Get session and user
	user, casErr := authenticateAPIUser(api, req)
	if casErr != nil {
		api.casServer.render.JSON(w, casErr.HttpCode, map[string]string{
			"status":  "error",
			"message": casErr.Msg,
		})
		return
	}

	// Build service from passed in data
	service := CASService{
		Name:       strings.TrimSpace(req.FormValue("name")),
		Url:        strings.TrimSpace(strings.ToLower(req.FormValue("url"))),
		AdminEmail: strings.TrimSpace(strings.ToLower(req.FormValue("adminEmail"))),
	}

	// Ensure user is admin
	if !user.IsAdmin {
		api.casServer.render.JSON(w, http.StatusUnauthorized, map[string]string{
			"status":  "error",
			"message": "Insufficient permissions.",
		})
		return
	}

	// Attempt to add service
	casErr = api.casServer.Db.AddNewService(&service)
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
		api.casServer.render.JSON(w, http.StatusUnauthorized, map[string]string{
			"status":  "error",
			"message": "Insufficient permissions.",
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
		api.casServer.render.JSON(w, http.StatusUnauthorized, map[string]string{
			"status":  "error",
			"message": "Insufficient permissions.",
		})
		return
	}

	service := CASService{
		Name:       strings.TrimSpace(req.FormValue("name")),
		Url:        strings.TrimSpace(strings.ToLower(req.FormValue("url"))),
		AdminEmail: strings.TrimSpace(strings.ToLower(req.FormValue("adminEmail"))),
	}

	// Attempt to update the service
	casErr = api.casServer.Db.UpdateService(&service)
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
