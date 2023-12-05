package authentication

import (
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/database"
	"cloudview/app/src/helpers"
	"cloudview/app/src/helpers/constants"
	"cloudview/app/src/models"
	project_members_model "cloudview/app/src/models/project_members"
	jwtAuth "cloudview/app/src/providers/oauth/jwt"
	"cloudview/app/src/types"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func SetSession(w http.ResponseWriter, s types.SessionUser) error {
	ENV := helpers.GoDotEnvVariable("ENV")
	token, err := jwtAuth.GenerateToken(&s)
	if err != nil {
		return errors.New("Error generating JWT token")
	}
	logger.Logger.Log("JWT token generated successfully")
	// Cookie valid for 1 year
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:     constants.COOKIE_NAME,
		Value:    token,
		HttpOnly: true,
		SameSite: 0,
		Secure:   false,
		Path:     "/",
		Domain:   "localhost",
		Expires:  expiration,
	}
	if strings.ToLower(ENV) == "production" {
		cookie.Domain = "" // Add domain
		cookie.Secure = true
	}
	http.SetCookie(w, &cookie)

	logger.Logger.Log("authentication.SetSession: session data successfully set")

	return nil
}

func GetAuthToken(r *http.Request) (string, error) {
	/*
		`Authorization` header is not passed
		since cookie is not available on the client-side.

		`Authorization` is only available when server side api calls are made.
		Therefore, we check for both cookie or auth header.
	*/
	auth_header := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth_header, "Bearer") {
		logger.Logger.Log("Authorization header not found, checking for cookie..")
		tok, err := r.Cookie(constants.COOKIE_NAME)
		if err != nil {
			logger.Logger.Error("Cookie not found", err)
			return "", custom_errors.MissingAuthToken
		}
		logger.Logger.Log("Cookie header found")
		return tok.Value, nil
	}
	logger.Logger.Log("Authorization header found")
	return strings.TrimPrefix(auth_header, "Bearer "), nil
}

/*
Verify if the auth user has access to project and services
*/
func VerifyProjectAccess(db *database.DB, projectId interface{}, userId uuid.UUID) (*models.Projects, error) {
	var projectUUID uuid.UUID
	switch v := projectId.(type) {
	case string:
		isValid := helpers.IsValidUUID(projectId.(string))
		if !isValid {
			logger.Logger.Error("authentication.VerifyProjectAccess: ERROR invalid projectId type, got:", v)
			return nil, errors.New("Invalid `projectId` type. Expected uuid string, got " + v)
		}
		uid, err := uuid.Parse(projectId.(string))
		if err != nil {
			logger.Logger.Error("authentication.VerifyProjectAccess: ERROR Unable to parse uuid", err)
			return nil, custom_errors.UnknownError
		}
		projectUUID = uid
		break
	case uuid.UUID:
		projectUUID = projectId.(uuid.UUID)
		break
	default:
		logger.Logger.Error("authentication.VerifyProjectAccess: ERROR invalid projectId type 'unknown'.")
		return nil, errors.New("Invalid `projectId` type. Expected uuid string, got `unknown`.")
	}

	result, err := project_members_model.GetProjectByIdAndUserId(db, projectUUID, userId)
	if err != nil {
		logger.Logger.Error("authentication.VerifyProjectAccess: ERROR Unable to fetch project", err)
		return nil, errors.New("You do not have access to this project.")
	}
	return &result, nil
}
