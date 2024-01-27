package middleware

import (
	"cloudview/app/src/api/authentication"
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/database"
	users_model "cloudview/app/src/models/users"
	jwtAuth "cloudview/app/src/providers/oauth/jwt"
	"cloudview/app/src/types"
	"encoding/json"
	"errors"
	"net/http"
)

type HttpResponseWriter interface {
	http.ResponseWriter
	success(data interface{}, status int)
	error(message string, status int)
	forbidden()
	unauthorized()
	notFound()
	user() (interface{}, error)
}

type HttpResponseWriterImpl struct {
	http.ResponseWriter
	ErrorMessage string
}

func (rw HttpResponseWriterImpl) Success(data interface{}, status int) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	response := map[string]interface{}{
		"success": true,
		"data":    data,
	}
	json.NewEncoder(rw).Encode(response)
}

func (rw HttpResponseWriterImpl) Error(message string, status int) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	response := map[string]interface{}{
		"error":   true,
		"message": message,
	}
	json.NewEncoder(rw).Encode(response)
}

func (rw HttpResponseWriterImpl) Forbidden() {
	msg := rw.ErrorMessage
	if msg == "" {
		msg = "You are not allowed to perform this action"
	}
	rw.ErrorMessage = ""
	rw.Error(msg, http.StatusForbidden)
}

func (rw HttpResponseWriterImpl) Unauthorized() {
	rw.Error("Unauthorized", http.StatusUnauthorized)
}

func (rw HttpResponseWriterImpl) NotFound() {
	rw.Error("Not Found", http.StatusNotFound)
}

/*
Fetch user data from cookie
*/
func (rw HttpResponseWriterImpl) User(db *database.DB, r *http.Request) (*types.SessionUser, error) {
	logger.Logger.Log("fetching user from jwt auth token")
	token, err := authentication.GetAuthToken(r)
	if err != nil {
		logger.Logger.Error("Auth error", err.Error())
		return nil, custom_errors.NoDataFound
	}
	result, err := jwtAuth.DecodeToken(token)
	if err != nil {
		logger.Logger.Error("middleware.responses.User: ERROR", err)
		return nil, custom_errors.NoDataFound
	}
	logger.Logger.Log(result.ID, result.Email)
	/*
		Verify if the user is in our system.
	*/
	if _, err := users_model.GetByEmail(db, result.Email); err != nil {
		if errors.Is(err, custom_errors.NoDataFound) {
			logger.Logger.Error("middleware.response.User: ERROR User not found in DB", err)
			return nil, custom_errors.NoDataFound
		}
		return nil, err
	}
	return result, nil
}

func RegisterResponses(w http.ResponseWriter) *HttpResponseWriterImpl {
	resp := HttpResponseWriterImpl{w, ""}
	return &resp
}

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
