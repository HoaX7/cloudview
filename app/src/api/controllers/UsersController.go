package controllers

import (
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/api/middleware"
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/database"
	users_model "cloudview/app/src/models/users"
	"errors"
	"net/http"
)

// To pass Metadata see the following example:
/*
metadataData := struct {
    Key1 string `json:"key1"`
    Key2 string `json:"key2"`
}{
    Key1: "value1",
    Key2: "value2",
}

metadataJSON, err := json.Marshal(metadataData)
if err != nil {
    // Handle the error
}

user := model.Users{
    Username: username,
    Email: email,
    Metadata: json.RawMessage(metadataJSON),
}
*/

var logu = logger.NewLogger()

// For testing
func (c *UsersController) CreateUser(db *database.DB) http.HandlerFunc {
	logu.SetName(c.Name() + ".CreateUser")
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		var (
			username = "hoax"
			email    = "test@gmail.com"
		)
		user := users_model.Users{
			Username: username,
			Email:    email,
		}
		logu.Log("Creating new user: ", user)
		result, err := users_model.Create(db, user)
		if err != nil {
			if errors.Is(err, custom_errors.UniqueConstraintViolation) {
				rw.Error("Email already in use", http.StatusConflict)
				return
			}
			logu.Error("ERROR", err)
			rw.Error(custom_errors.UnknownError.Error(), http.StatusBadRequest)
			return
		}
		rw.Success(result, http.StatusOK)
		return
	}
}

func (c *UsersController) GetUserByEmail(db *database.DB) http.HandlerFunc {
	logu.SetName(c.Name() + ".GetUserByEmail")
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		email := "test@gmail.com"
		logu.Log("Fetching user with email: ", email)
		result, err := users_model.GetByEmail(db, email)
		if err != nil {
			if errors.Is(err, custom_errors.NoDataFound) {
				rw.Error(custom_errors.NoDataFound.Error(), http.StatusNotFound)
				return
			}
			logu.Error("ERROR", err)
			rw.Error(custom_errors.UnknownError.Error(), http.StatusBadRequest)
			return
		}
		rw.Success(result, http.StatusOK)
		return
	}
}
