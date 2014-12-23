package cas

import (
	"github.com/gorilla/mux"
	"net/http"
)

/*
 * CAS API implementation
 */

// Get the services for a logged in user
func (c *CAS) listSessionUserServices(w http.ResponseWriter, req *http.Request) {
	// Get the current session
	session, err := c.cookieStore.Get(req, "casgo-session")
	if err != nil {
		c.render.JSON(w, http.StatusInternalServerError, map[string]string{
			"status":  "error",
			"message": "Failed to retrieve services for logged in user. Please ensure you are logged in.",
		})
		return
	}

	// Retrieve information from Check whether the user is an admin
	userIsAdmin, isAdminOk := session.Values["userIsAdmin"].(bool)
	userEmail, emailOk := session.Values["userIsAdmin"].(string)
	userServices, servicesOK := session.Values["userServices"].([]CASService)
	if !isAdminOk || !emailOk || !servicesOK {
		c.render.JSON(w, http.StatusInternalServerError, map[string]string{
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
		c.render.JSON(w, http.StatusUnauthorized, map[string]string{
			"status":  "error",
			"message": "Insufficient permissions",
		})
	}

	// Return the user's services
	c.render.JSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   userServices,
	})
}

// Services Handler
func (c *CAS) ServicesHandler(w http.ResponseWriter, req *http.Request) {

}

// Users Handler
func (c *CAS) UsersHandler(w http.ResponseWriter, req *http.Request) {

}
