package cas

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	"strings"
)

/*
 * CAS FrontendAPI implementation
 */

func NewCasgoFrontendAPI(c *CAS) (*FrontendAPI, error) {
	return &FrontendAPI{casServer: c}, nil
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
	_, user, casErr := getSessionAndUser(api, req)
	if casErr != nil {
		api.casServer.render.JSON(w, casErr.httpCode, map[string]string{
			"status":  "error",
			"message": casErr.msg,
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
	_, user, casErr := getSessionAndUser(api, req)
	if casErr != nil {
		api.casServer.render.JSON(w, casErr.httpCode, map[string]string{
			"status":  "error",
			"message": casErr.msg,
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

// Utility function to retrieve session and user information
func getSessionAndUser(api *FrontendAPI, req *http.Request) (*sessions.Session, *User, *CASServerError) {
	// Get the current session
	session, err := api.casServer.cookieStore.Get(req, "casgo-session")
	if err != nil {
		casErr := &FailedToRetrieveServicesError
		casErr.err = &err
		return nil, nil, casErr
	}

	// Retrieve information from Check whether the user is an admin
	user, ok := session.Values["currentUser"].(User)
	if !ok {
		casErr := &FailedToRetrieveInformationFromSessionError
		casErr.err = &err
		return nil, nil, casErr
	}

	return session, &user, nil
}

// Get list of services (admin only)
func (api *FrontendAPI) GetServices(w http.ResponseWriter, req *http.Request) {
	// Get the current session and user
	_, user, casErr := getSessionAndUser(api, req)
	if casErr != nil {
		api.casServer.render.JSON(w, casErr.httpCode, map[string]string{
			"status":  "error",
			"message": casErr.msg,
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
		api.casServer.render.JSON(w, casErr.httpCode, map[string]string{
			"status":  "error",
			"message": casErr.msg,
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
	_, user, casErr := getSessionAndUser(api, req)
	if casErr != nil {
		api.casServer.render.JSON(w, casErr.httpCode, map[string]string{
			"status":  "error",
			"message": casErr.msg,
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
		api.casServer.render.JSON(w, casErr.httpCode, map[string]string{
			"status":  "error",
			"message": casErr.msg,
		})
		return
	}

	api.casServer.render.JSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   service,
	})
}

// Remove a service
func (api *FrontendAPI) RemoveService(w http.ResponseWriter, req *http.Request) {
	// Get session and user
	_, user, casErr := getSessionAndUser(api, req)
	if casErr != nil {
		api.casServer.render.JSON(w, casErr.httpCode, map[string]string{
			"status":  "error",
			"message": casErr.msg,
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
		api.casServer.render.JSON(w, casErr.httpCode, map[string]string{
			"status":  "error",
			"message": casErr.msg,
		})
	}

	api.casServer.render.JSON(w, http.StatusOK, map[string]string{
		"status": "success",
		"data":   serviceName,
	})
}

// Update an existing service
func (api *FrontendAPI) UpdateService(w http.ResponseWriter, req *http.Request) {
	// Build service object from passed in data
	// api.casServer.Db.UpdateService(service)
}
