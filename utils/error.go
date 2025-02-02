package utils

import (
	"errors"
	"net/http"
)

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	if e.Message == "" {
		return "err: not found"
	}
	return e.Message
}

type InsertionError struct {
	Message string
}

func (e *InsertionError) Error() string {
	if e.Message == "" {
		return "insertion error"
	}
	return e.Message
}

type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	if e.Message == "" {
		return "error: bad request"
	}
	return e.Message
}

type AuthenticationError struct {
	Message string
}

func (e *AuthenticationError) Error() string {
	if e.Message == "" {
		return "full authentication is required to access this resource"
	}
	return e.Message
}

type AccessDeniedError struct {
	Message string
}

func (e *AccessDeniedError) Error() string {
	if e.Message == "" {
		return "access denied"
	}
	return e.Message
}

func errorStatus(err error) int {
	var notFoundError *NotFoundError
	var insertionError *InsertionError
	var badRequestError *BadRequestError
	var authenticationError *AuthenticationError
	var accessDeniedError *AccessDeniedError
	switch {
	case errors.As(err, &notFoundError):
		return http.StatusNotFound
	case errors.As(err, &insertionError):
		return http.StatusConflict
	case errors.As(err, &badRequestError):
		return http.StatusBadRequest
	case errors.As(err, &authenticationError):
		return http.StatusUnauthorized
	case errors.As(err, &accessDeniedError):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
