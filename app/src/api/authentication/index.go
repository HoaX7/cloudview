package authentication

import (
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/database"
	"cloudview/app/src/helpers"
	"cloudview/app/src/helpers/constants"
	"cloudview/app/src/models"
	project_members_model "cloudview/app/src/models/project_members"
	provider_accounts_model "cloudview/app/src/models/provider_accounts"
	jwtAuth "cloudview/app/src/providers/oauth/jwt"
	"cloudview/app/src/types"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func SetSession(w http.ResponseWriter, s types.SessionUser) error {
	ENV := helpers.GoDotEnvVariable("GO_ENV")
	DOMAIN := helpers.GoDotEnvVariable("EXTERNAL_DOMAIN")
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
		cookie.Domain = "." + DOMAIN // Add domain
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
Verify if the auth user has access to project and services.

This function returns member permissions and minimal project details.
If you pass providerAccId in the input - the data returned also contains
minimal account details and acc permissions.
*/
func VerifyProjectAccess(db *database.DB, userId uuid.UUID, input types.VerifyProjectAccessInput) (*types.VerifyProjectAccessOutput, error) {
	var projectUUID uuid.UUID
	var providerAccount models.ProviderAccounts
	if input.ProjectID == nil && input.ProviderAccountID == nil {
		return nil, errors.New("Invalid 'VerifyProjectAccessInput'. Expected 'ProjectID' or 'ProviderAccountID'")
	}
	projectId := input.ProjectID
	providerAccountId := input.ProviderAccountID
	if providerAccountId != nil {
		uuid, err := checkUUID(providerAccountId, "providerAccountId")
		if err != nil {
			return nil, err
		}
		providerAcc, err := provider_accounts_model.GetById(db, *uuid)
		if err != nil {
			logger.Logger.Error("VerifyProjectAccess: Error fetching provider account", err)
			return nil, custom_errors.UnknownError
		}
		projectUUID = providerAcc.Project.ID
		providerAccount = providerAcc.ProviderAccounts
	} else {
		uuid, err := checkUUID(projectId, "projectId")
		if err != nil {
			return nil, err
		}
		projectUUID = *uuid
	}

	result, err := project_members_model.GetProjectByIdAndUserId(db, projectUUID, userId)
	if err != nil {
		logger.Logger.Error("authentication.VerifyProjectAccess: ERROR Unable to fetch project", err)
		return nil, errors.New("You do not have access to this project.")
	}
	res := &types.VerifyProjectAccessOutput{
		ProjectAccessDetails: result,
		ProviderAccount:      &providerAccount,
	}
	return res, nil
}

func checkUUID(id interface{}, checkString string) (*uuid.UUID, error) {
	var result uuid.UUID
	errorStr := fmt.Sprintf("Invalid `%s` type. Expected uuid string, got ", checkString)
	switch v := id.(type) {
	case string:
		isValid, uid := helpers.IsValidUUID(id.(string))
		if !isValid {
			logger.Logger.Error("authentication.checkUUID: ERROR", errorStr, v)
			return nil, errors.New(errorStr + v)
		}
		result = *uid
		break
	case uuid.UUID:
		result = id.(uuid.UUID)
		if helpers.IsDummyUUID(result) {
			return nil, errors.New(errorStr + "'Unknown")
		}
		break
	default:
		logger.Logger.Error("authentication.checkUUID: ERROR ", errorStr, "'unknown'")
		return nil, errors.New(errorStr + "'Unkown'")
	}
	logger.Logger.Log("authenticate.checkUUID: success", &result)
	return &result, nil
}
