package cas

import (
	"net/http"
)

func (err *CASServerError) Error() string { return err.msg }

// Error declarations
var (
	// Input errors (error codes 100-199)
	InvalidEmailAddressError = CASServerError{
		msg:        "An error occurred finding a user with that email address.. Please wait a while and try again",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 100,
	}

	InvalidCredentialsError = CASServerError{
		msg:        "Invalid email/password combination",
		httpCode:   http.StatusUnauthorized,
		casErrCode: 101,
	}
	FailedToFindServiceError = CASServerError{
		msg:        "Failed to find matching service",
		httpCode:   http.StatusNotImplemented,
		casErrCode: 102,
	}
	FailedToFindTicketError = CASServerError{
		msg:        "Failed to find matching ticket",
		httpCode:   http.StatusNotImplemented,
		casErrCode: 103,
	}
	SSOAuthenticatedUserRenewError = CASServerError{
		msg:        "Failed to validate ticket, renew option specified and user was SSO authenticated",
		httpCode:   http.StatusNotImplemented,
		casErrCode: 103,
	}
	EmailAlreadyTakenError = CASServerError{
		msg:        "Looks like that email address is already taken. If you've forgotten your password, please contact the administrator",
		httpCode:   http.StatusBadRequest,
		casErrCode: 104,
	}

	// Internal Server errors (error codes 200 - 299)
	FailedToSaveSessionError = CASServerError{
		msg:        "Failed to save session",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 200,
	}
	FailedToDeleteSessionError = CASServerError{
		msg:        "Failed to delete session",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 201,
	}
	FailedToCreateNewAuthTicketError = CASServerError{
		msg:        "Failed to create new authentication ticket",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 202,
	}
	AuthMethodNotSupportedError = CASServerError{
		msg:        "Failed to create new authentication ticket",
		httpCode:   http.StatusMethodNotAllowed,
		casErrCode: 203,
	}
	FailedToCreateUserError = CASServerError{
		msg:        "An error occurred while creating your account.. Please verify fields and try again",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 204,
	}
	FailedToTeardownDatabase = CASServerError{
		msg:        "Failed to tear down database",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 205,
	}
	FailedToSetupDatabase = CASServerError{
		msg:        "Failed to setup database",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 206,
	}
	FailedToImportTableDataFromFile = CASServerError{
		msg:        "Failed to import database information from file",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 207,
	}

	// Other (error codes 300 - 399)
	UnsupportedFeatureError = CASServerError{
		msg:        "Feature not supported by CASGO",
		httpCode:   http.StatusNotImplemented,
		casErrCode: 300,
	}
)
