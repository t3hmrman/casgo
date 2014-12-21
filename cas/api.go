package cas

import (
	"net/http"
)

/*
 * CAS API implementation
 */

// Get the services for a logged in user
func (c *CAS) getServicesForUser(w http.ResponseWriter, req *http.Request) {
	// Get the current session
	session, err := c.cookieStore.Get(req, "casgo-session")
	if err != nil {
		c.render.JSON(w, http.StatusInternalServerError, map[string]string{
			"status": "error",
			"message": "Failed to retrieve services for logged in user. Please ensure you are logged in.",
		})
		return
	}

	// Get the user's services out of the session
	services, ok := session.Values["userServices"]
	if !ok {
		c.render.JSON(w, http.StatusInternalServerError, map[string]string{
			"status": "error",
			"message": "Failed to retrieve services for logged in user. Please ensure you are logged in.",
		})
		return
	}
	
	// Return the user's services
	c.render.JSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data": services,
	})
}
