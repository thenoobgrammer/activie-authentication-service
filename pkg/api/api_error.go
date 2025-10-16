package api

import (
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int `json:"statusCode"`
	Msg        any `json:"msg"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: %d", e.StatusCode)
}

func NewAPIError(statusCode int, err error) APIError {
	return APIError{
		StatusCode: statusCode,
		Msg:        err.Error(),
	}
}

func InvalidRequestData(errors map[string]string) APIError {
	return APIError{
		StatusCode: http.StatusUnprocessableEntity,
		Msg:        errors,
	}
}

func EmailExists() APIError {
	return NewAPIError(http.StatusConflict, fmt.Errorf("email already in use"))
}

func FailedToCreateSession() APIError {
	return NewAPIError(http.StatusInternalServerError, fmt.Errorf("failed to create session"))
}

func FailedToDeleteSession() APIError {
	return NewAPIError(http.StatusInternalServerError, fmt.Errorf("failed.to.delete.session"))
}

func InvalidJSON() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid JSON request data"))
}

func InvalidParamRequest() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid param request"))
}

func InvalidHeaderRequest() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid header request"))
}

func InvalidBearer() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid bearer"))
}

func SessionNotFound() APIError {
	return NewAPIError(http.StatusNotFound, fmt.Errorf("session not found - please login"))
}

func InvalidClaims() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid claims"))
}

func InvalidCredentials() APIError {
	return NewAPIError(http.StatusUnauthorized, fmt.Errorf("invalid credentials"))
}

func BadRequest() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("bad request"))
}

func NotFound() APIError {
	return NewAPIError(http.StatusNotFound, fmt.Errorf("not found"))
}

func Unauthorized() APIError {
	return NewAPIError(http.StatusUnauthorized, fmt.Errorf("unauthorized"))
}

func Forbidden() APIError {
	return NewAPIError(http.StatusForbidden, fmt.Errorf("forbidden access"))
}

func FailedToGenerateBearerToken() APIError {
	return NewAPIError(http.StatusInternalServerError, fmt.Errorf("failed to generate bearer token"))
}

func OperationFailed() APIError {
	return NewAPIError(http.StatusInternalServerError, fmt.Errorf("operation failed"))
}

type InternalServerError struct {
	Err error
}

func (e InternalServerError) Error() string {
	return fmt.Sprintf("internal server error: %v", e.Err)
}

func NewInternalServerError(err error) InternalServerError {
	return InternalServerError{
		Err: err,
	}
}
