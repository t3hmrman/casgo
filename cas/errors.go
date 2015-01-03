package cas

import (
	"net/http"
)

func (err *CASServerError) Error() string { return err.Msg }

// Error declarations
var (
	// Input errors (error codes 100-199)
	InvalidEmailAddressError = CASServerError{
		Msg:        "An error occurred finding a user with that email address.. Please wait a while and try again",
		HttpCode:   http.StatusInternalServerError,
		CasgoErrCode: 100,
	}

	InvalidCredentialsError = CASServerError{
		Msg:        "Invalid email/password combination",
		HttpCode:   http.StatusUnauthorized,
		CasgoErrCode: 101,
	}
	FailedToFindServiceError = CASServerError{
		Msg:        "Failed to find matching service",
		HttpCode:   http.StatusNotImplemented,
		CasgoErrCode: 102,
	}
	FailedToFindTicketError = CASServerError{
		Msg:        "Failed to find matching ticket",
		HttpCode:   http.StatusBadRequest,
		CasgoErrCode: 103,
	}
	SSOAuthenticatedUserRenewError = CASServerError{
		Msg:        "Failed to validate ticket, renew option specified and user was SSO authenticated",
		HttpCode:   http.StatusNotImplemented,
		CasgoErrCode: 103,
	}
	EmailAlreadyTakenError = CASServerError{
		Msg:        "Looks like that email address is already taken. If you've forgotten your password, please contact the administrator",
		HttpCode:   http.StatusBadRequest,
		CasgoErrCode: 104,
	}
	FailedToFindUserError = CASServerError{
		Msg:        "Failed to find matching email/password combination",
		HttpCode:   http.StatusBadRequest,
		CasgoErrCode: 105,
	}
	FailedToRetrieveServicesError = CASServerError{
		Msg:        "Failed to retrieve services for logged in user. Please ensure you are logged in.",
		HttpCode:   http.StatusBadRequest,
		CasgoErrCode: 106,
	}
	ServiceNameAlreadyTakenError = CASServerError{
		Msg:        "Looks like that service name is already taken. Please use a different service name.",
		HttpCode:   http.StatusBadRequest,
		CasgoErrCode: 107,
	}
	InvalidServiceNameError = CASServerError{
		Msg:        "Invalid service name provided.",
		HttpCode:   http.StatusBadRequest,
		CasgoErrCode: 108,
	}
	FailedToAuthenticateUserError = CASServerError{
		Msg:        "Failed to authenticate API user. Please ensure that you have provided sufficient credentials (whether through relevant headers or session information)..",
		HttpCode:   http.StatusUnauthorized,
		CasgoErrCode: 109,
	}

	// Internal Server errors (error codes 200 - 299)
	FailedToSaveSessionError = CASServerError{
		Msg:        "Failed to save session",
		HttpCode:   http.StatusInternalServerError,
		CasgoErrCode: 200,
	}
	FailedToDeleteSessionError = CASServerError{
		Msg:        "Failed to delete session",
		HttpCode:   http.StatusInternalServerError,
		CasgoErrCode: 201,
	}
	FailedToCreateNewAuthTicketError = CASServerError{
		Msg:        "Failed to create new authentication ticket",
		HttpCode:   http.StatusInternalServerError,
		CasgoErrCode: 202,
	}
	AuthMethodNotSupportedError = CASServerError{
		Msg:        "Failed to create new authentication ticket",
		HttpCode:   http.StatusMethodNotAllowed,
		CasgoErrCode: 203,
	}
	FailedToCreateUserError = CASServerError{
		Msg:        "An error occurred while creating your account.. Please verify fields and try again",
		HttpCode:   http.StatusInternalServerError,
		CasgoErrCode: 204,
	}
	FailedToTeardownDatabaseError = CASServerError{
		Msg:        "Failed to tear down database",
		HttpCode:   http.StatusInternalServerError,
		CasgoErrCode: 205,
	}
	FailedToSetupDatabaseError = CASServerError{
		Msg:        "Failed to setup database",
		HttpCode:   http.StatusInternalServerError,
		CasgoErrCode: 206,
	}
	FailedToLoadJSONFixtureError = CASServerError{
		Msg:        "Failed to import database information from file",
		HttpCode:   http.StatusInternalServerError,
		CasgoErrCode: 207,
	}
	FailedToLookupServiceByUrlError = CASServerError{
		Msg:        "An error occurred while searching for service with given URL",
		HttpCode:   http.StatusInternalServerError,
		CasgoErrCode: 208,
	}
	FailedToCreateTicketError = CASServerError{
		Msg:        "Failed to create ticket",
		HttpCode:   http.StatusInternalServerError,
		CasgoErrCode: 209,
	}
	FailedToDeleteTicketsForUserError = CASServerError{
		Msg:        "Failed to delete tickets for user",
		HttpCode:   http.StatusInternalServerError,
		CasgoErrCode: 210,
	}
	FailedToSetupTableError = CASServerError{
		Msg:        "Failed to setup table",
		HttpCode:   http.StatusInternalServerError,
		CasgoErrCode: 211,
	}
	FailedToCreateTableError = CASServerError{
		Msg:        "Failed to setup database",
		HttpCode:   http.StatusInternalServerError,
		CasgoErrCode: 212,
	}
	DbExistsCheckFailedError = CASServerError{
		Msg:        "Failed to check whether database existed",
		HttpCode:   http.StatusInternalServerError,
		CasgoErrCode: 213,
	}
	FailedToFindServiceByUrlError = CASServerError{
		Msg:        "Failed to find service with given URL",
		HttpCode:   http.StatusInternalServerError,
		CasgoErrCode: 214,
	}
	FailedToFindUserByEmailError = CASServerError{
		Msg:        "Failed to find user with given email address",
		HttpCode:   http.StatusInternalServerError,
		CasgoErrCode: 215,
	}
	FailedToRetrieveInformationFromSessionError = CASServerError{
		Msg:        "Failed to retrieve information from session",
		HttpCode:   http.StatusInternalServerError,
		CasgoErrCode: 216,
	}
	FailedToCreateServiceError = CASServerError{
		Msg:        "An error occurred while creating the service... Please verify fields and try again",
		HttpCode:   http.StatusInternalServerError,
		CasgoErrCode: 217,
	}
	FailedToDeleteServiceError = CASServerError{
		Msg:        "Failed to delete service.",
		HttpCode:   http.StatusInternalServerError,
		CasgoErrCode: 218,
	}
	FailedToListServicesError = CASServerError{
		Msg:        "Failed to list services.",
		HttpCode:   http.StatusInternalServerError,
		CasgoErrCode: 219,
	}
	FailedToUpdateServiceError = CASServerError{
		Msg:        "Failed to update service.",
		HttpCode:   http.StatusInternalServerError,
		CasgoErrCode: 220,
	}

	// Other (error codes 300 - 399)
	UnsupportedFeatureError = CASServerError{
		Msg:        "Feature not supported by CASGO",
		HttpCode:   http.StatusNotImplemented,
		CasgoErrCode: 300,
	}
)
