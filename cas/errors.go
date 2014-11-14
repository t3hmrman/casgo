package cas

import(
	"net/http"
)

// Custom type for CAS server errors
type CASServerError struct {
	msg string
	http_code int
}

func (err *CASServerError) Error() string { return err.msg }


// Error declarations
var(
	InvalidEmailAddressError = CASServerError{"An error occurred finding a user with that email address.. Please wait a while and try again", http.StatusInternalServerError}
	InvalidCredentialsError = CASServerError{"Invalid email/password combination", http.StatusUnauthorized}
)
