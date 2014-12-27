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
		httpCode:   http.StatusBadRequest,
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
	FailedToFindUserError = CASServerError{
		msg:        "Failed to find matching email/password combination",
		httpCode:   http.StatusBadRequest,
		casErrCode: 105,
	}
	FailedToRetrieveServicesError = CASServerError{
		msg: "Failed to retrieve services for logged in user. Please ensure you are logged in.",
		httpCode: http.StatusBadRequest,
		casErrCode: 106,
	}
	ServiceNameAlreadyTakenError = CASServerError{
		msg:        "Looks like that service name is already taken. Please use a different service name.",
		httpCode:   http.StatusBadRequest,
		casErrCode: 107,
	}
	InvalidServiceNameError = CASServerError{
		msg:        "Invalid service name provided.",
		httpCode:   http.StatusBadRequest,
		casErrCode: 107,
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
	FailedToTeardownDatabaseError = CASServerError{
		msg:        "Failed to tear down database",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 205,
	}
	FailedToSetupDatabaseError = CASServerError{
		msg:        "Failed to setup database",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 206,
	}
	FailedToLoadJSONFixtureError = CASServerError{
		msg:        "Failed to import database information from file",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 207,
	}
	FailedToLookupServiceByUrlError = CASServerError{
		msg:        "An error occurred while searching for service with given URL",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 208,
	}
	FailedToCreateTicketError = CASServerError{
		msg:        "Failed to create ticket",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 209,
	}
	FailedToDeleteTicketsForUserError = CASServerError{
		msg:        "Failed to delete tickets for user",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 210,
	}
	FailedToSetupTableError = CASServerError{
		msg:        "Failed to setup table",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 211,
	}
	FailedToCreateTableError = CASServerError{
		msg:        "Failed to setup database",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 212,
	}
	DbExistsCheckFailedError = CASServerError{
		msg:        "Failed to check whether database existed",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 213,
	}
	FailedToFindServiceByUrlError = CASServerError{
		msg:        "Failed to find service with given URL",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 214,
	}
	FailedToFindUserByEmailError = CASServerError{
		msg:        "Failed to find user with given email address",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 215,
	}
	FailedToRetrieveInformationFromSessionError = CASServerError{
		msg:        "Failed to retrieve information from session",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 216,
	}
	FailedToCreateServiceError = CASServerError{
		msg:        "An error occurred while creating the service... Please verify fields and try again",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 217,
	}
	FailedToDeleteServiceError = CASServerError{
		msg:        "Failed to delete service.",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 218,
	}
	FailedToListServicesError = CASServerError{
		msg:        "Failed to list services.",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 219,
	}
	FailedToUpdateServiceError = CASServerError{
		msg:        "Failed to update service.",
		httpCode:   http.StatusInternalServerError,
		casErrCode: 220,
	}


	// Other (error codes 300 - 399)
	UnsupportedFeatureError = CASServerError{
		msg:        "Feature not supported by CASGO",
		httpCode:   http.StatusNotImplemented,
		casErrCode: 300,
	}
)
