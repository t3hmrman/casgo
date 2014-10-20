package cas

/**
 * CAS protocol V3 API
 */

import (
	"net/http"
	"io"
)


func HandleLogin(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Login!")
}