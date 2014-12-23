package cas

import (
	"github.com/gorilla/mux"
	"net/http"
)

/*
 * CAS FrontendAPI implementation
 */

func NewCasgoFrontendAPI(c *CAS) (*FrontendAPI, error) {
	return &FrontendAPI{casServer: c}, nil
}

// Hook up API endpoints to given mux
func (api *FrontendAPI) HookupAPIEndpoints(m *mux.Router) {
	m.HandleFunc("/api/sessions/{userEmail}/services", api.listSessionUserServices).Methods("GET")
	m.HandleFunc("/api/sessions", api.SessionsHandler).Methods("GET")
	m.HandleFunc("/api/services", api.GetService).Methods("GET")
	m.HandleFunc("/api/services", api.CreateService).Methods("POST")
	m.HandleFunc("/api/services/{serviceName}", api.RemoveService).Methods("DELETE")
}

// Handle sessions endpoint
func (api *FrontendAPI) SessionsHandler(w http.ResponseWriter, req *http.Request) {
	// Get the current session
	session, err := api.casServer.cookieStore.Get(req, "casgo-session")
	if err != nil {
		api.casServer.render.JSON(w, http.StatusInternalServerError, map[string]string{
			"status":  "error",
			"message": "Failed to retrieve services for logged in user. Please ensure you are logged in.",
		})
		return
	}

	// Retrieve information from Check whether the user is an admin
	userIsAdmin, isAdminOk := session.Values["userIsAdmin"].(bool)
	userEmail, emailOk := session.Values["userEmail"].(string)
	userServices, servicesOk := session.Values["userServices"].([]CASService)
	if !isAdminOk || !emailOk || !servicesOk {
		api.casServer.render.JSON(w, http.StatusInternalServerError, map[string]string{
			"status":  "error",
			"message": "Internal server error, Failed to retrieve user information from session.",
		})
		return
	}

	api.casServer.render.JSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"email":    userEmail,
			"isAdmin":  userIsAdmin,
			"services": userServices,
		},
	})
}

// Get the services for a logged in user
func (api *FrontendAPI) listSessionUserServices(w http.ResponseWriter, req *http.Request) {
	// Get the current session
	session, err := api.casServer.cookieStore.Get(req, "casgo-session")
	if err != nil {
		api.casServer.render.JSON(w, http.StatusInternalServerError, map[string]string{
			"status":  "error",
			"message": "Failed to retrieve services for logged in user. Please ensure you are logged in.",
		})
		return
	}

	// Retrieve information from Check whether the user is an admin
	userIsAdmin, isAdminOk := session.Values["userIsAdmin"].(bool)
	userEmail, emailOk := session.Values["userEmail"].(string)
	userServices, servicesOK := session.Values["userServices"].([]CASService)
	if !isAdminOk || !emailOk || !servicesOK {
		api.casServer.render.JSON(w, http.StatusInternalServerError, map[string]string{
			"status":  "error",
			"message": "Internal server error, Failed to retrieve user information from session.",
		})
		return
	}

	// Quit early if the user is not an admin and is not the requested user
	routeVars := mux.Vars(req)
	routeUserEmail := routeVars["userEmail"]

	// Ensure non-admin user is not trying to lookup another users session information
	if !userIsAdmin && userEmail != routeUserEmail {
		api.casServer.render.JSON(w, http.StatusUnauthorized, map[string]string{
			"status":  "error",
			"message": "Insufficient permissions",
		})
		return
	}

	// Return the user's services
	api.casServer.render.JSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   userServices,
	})
}

// Get a service
func (api *FrontendAPI) GetService(w http.ResponseWriter, req *http.Request) {
}

// Create a new service
func (api *FrontendAPI) CreateService(w http.ResponseWriter, req *http.Request) {
}

// Remove a service
func (api *FrontendAPI) RemoveService(w http.ResponseWriter, req *http.Request) {

}
