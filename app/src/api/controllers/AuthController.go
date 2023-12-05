package controllers

import (
	"cloudview/app/src/api/authentication"
	"cloudview/app/src/api/middleware"
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/database"
	"cloudview/app/src/helpers"
	"cloudview/app/src/helpers/constants"
	"cloudview/app/src/providers/oauth"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type OauthBody struct {
	Code string `json:"code"`
}

/*
This method allows github & google oauth logins
*/
func (c *AuthController) OauthLogin(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		logger.Logger.Log(r.Body)
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Logger.Error(err)
			rw.Error("Unable to parse json body", http.StatusBadRequest)
			return
		}
		var body OauthBody
		if err := json.Unmarshal(reqBody, &body); err != nil {
			rw.Error("Internal server error", http.StatusInternalServerError)
			return
		}
		params := mux.Vars(r)
		provider := params["provider"]
		if body.Code == "" {
			rw.Error("Invalid code provided", http.StatusBadRequest)
			return
		}
		logger.Logger.Log("AuthController.OauthLogin: authenticating " + provider + " provider with code: " + body.Code)
		switch strings.ToLower(provider) {
		case constants.OAuth.GOOGLE:
			// TODO - Unimplemented
			oauth.Login(oauth.Google{Code: body.Code}, db)
			rw.Success("Google login", http.StatusOK)
		case constants.OAuth.GITHUB:
			sessionUser, err := oauth.Login(oauth.Github{Code: body.Code}, db)
			if err != nil {
				rw.Error(err.Error(), http.StatusForbidden)
				break
			}
			authentication.SetSession(w, sessionUser)
			rw.Success(sessionUser, http.StatusOK)
		default:
			rw.NotFound()
		}
		return
	}
}

func (c *AuthController) OauthCallback(w http.ResponseWriter, r *http.Request) {
	rw := middleware.RegisterResponses(w)
	rw.Success("Login successful", http.StatusOK)
	return
}

/*
The way to remove cookie is to set cookie with empty values
and max-age of negative or zero value.
https://thinkingeek.com/2018/05/31/setting-and-deleting-cookies-in-go/
*/
func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	rw := middleware.RegisterResponses(w)
	ENV := helpers.GoDotEnvVariable("ENV")
	cookie := http.Cookie{
		Name:     constants.COOKIE_NAME,
		Value:    "",
		HttpOnly: true,
		SameSite: 0,
		Secure:   false,
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   -1,
	}
	if strings.ToLower(ENV) == "production" {
		cookie.Domain = "" // Add domain
		cookie.Secure = true
	}
	http.SetCookie(w, &cookie)
	rw.Success("Logout successful", http.StatusOK)
	return
}
