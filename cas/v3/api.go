package cas

/**
 * CAS protocol V3 API
 */

import (
	"net/http"
	"io"
)


// Credential acceptor endpoint (requestor is handled in main)
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Show login page if credentials are not provided, attempt login otherwise
	if (true) {
		io.WriteString(w, "Login!")
	} else {
		io.WriteString(w, "Login!")		
	}
}

// Endpoint for destroying CAS sessions (logging out) 
func HandleLogout(w http.ResponseWriter, r *http.Request) {

}

// Endpoint for validating service tickets
func HandleValidate(w http.ResponseWriter, r *http.Request) {

}

// Endpoint for validating service tickets (CAS 2.0)
func HandleServiceValidate(w http.ResponseWriter, r *http.Request) {

}

// Endpoint for validating proxy tickets (CAS 2.0)
func HandleProxyValidate(w http.ResponseWriter, r *http.Request) {

}

// Endpoint for handling proxy tickets (CAS 2.0)
func HandleProxy(w http.ResponseWriter, r *http.Request) {

}
