package cas

import(
	"net/http"
)

// Custom type for CAS server errors
type CASServerError struct {
	msg string
	http_code int
	err_code int
}

func (err *CASServerError) Error() string { return err.msg }

// Error declarations
var(
	// Input errors (error codes 100-199)
	InvalidEmailAddressError = CASServerError{"An error occurred finding a user with that email address.. Please wait a while and try again", http.StatusInternalServerError, 100}
	InvalidCredentialsError = CASServerError{"Invalid email/password combination", http.StatusUnauthorized, 101}
	FailedToFindServiceError = CASServerError{"Failed to find matching service", http.StatusNotImplemented, 102}
	FailedToFindTicketError = CASServerError{"Failed to find matching ticket", http.StatusNotImplemented, 103}
	SSOAuthenticatedUserRenewError = CASServerError{"Failed to validate ticket, renew option specified and user was SSO authenticated", http.StatusNotImplemented, 103}

	// Internal Server errors (error codes 200 - 299)
	FailedToSaveSessionError = CASServerError{"Failed to save session", http.StatusInternalServerError, 200}
	FailedToDeleteSessionError = CASServerError{"Failed to delete session", http.StatusInternalServerError, 201}
	FailedToCreateNewAuthTicketError = CASServerError{"Failed to create new authentication ticket", http.StatusInternalServerError, 202}
	AuthMethodNotSupportedError = CASServerError{"Failed to create new authentication ticket", http.StatusMethodNotAllowed, 203}

	// Other (error codes 300 - 399)
	UnsupportedFeatureError = CASServerError{"Feature not supported by CASGO", http.StatusNotImplemented, 300}
)
