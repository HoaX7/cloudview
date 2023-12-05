package custom_errors

import (
	"errors"
)

var (
	UnknownError       = errors.New("Unknown error")
	MissingAuthToken   = errors.New("No authorization token found in header")
	InvalidJWTToken    = errors.New("Invalid JWT token")
	Forbidden          = errors.New("Forbidden")
	InvalidCredentials = errors.New("Invalid credentials provided")
)
