package controllers

import (
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/api/middleware"
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/database"
	"cloudview/app/src/helpers"
	"cloudview/app/src/models"
	project_members_model "cloudview/app/src/models/project_members"
	projects_model "cloudview/app/src/models/projects"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

/*
TODO - Decide the `member_limit` prop value based on the user
subscription plan. Default value is 1
*/
func (c *ProjectsController) CreateProject(db *database.DB) http.HandlerFunc {
	/*
		TODO - Add middleware to check if the user
		has the correct subscription plan to be able to
		create more than 1 project.

		Pricing Model also includes
		cost per contributor.

	*/
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		authenticatedUser, err := rw.User(db, r)
		if err != nil {
			logger.Logger.Error("ProjectsController.CreateProject: ERROR", err)
			if errors.Is(err, custom_errors.NoDataFound) {
				rw.Unauthorized()
				return
			}
			rw.Forbidden()
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Logger.Error("ProjectsController.CreateProject: Error reading request body:", err)
			rw.Error("Bad request", http.StatusUnprocessableEntity)
			return
		}
		var request models.Projects
		if err := json.Unmarshal(body, &request); err != nil {
			logger.Logger.Error("ProjectsController.CreateProject: Error parsing request body:", err)
			rw.Error(custom_errors.UnknownError.Error(), http.StatusUnprocessableEntity)
			return
		}
		defer r.Body.Close()
		if request.Name == "" {
			rw.Error("Required field `name`", http.StatusBadRequest)
			return
		}
		if request.Type == "" || (request.Type != "PUBLIC" && request.Type != "PRIVATE") {
			rw.Error("Required enum field `type` must be one of value ('PUBLIC', 'PRIVATE')", http.StatusBadRequest)
			return
		}
		request.OwnerID = &authenticatedUser.ID
		request.Email = authenticatedUser.Email
		if request.MemberLimit <= 0 {
			request.MemberLimit = 1
		}
		logger.Logger.Log("ProjectsController.CreateProject: creating new project for owner: ", request.OwnerID)
		result, err := projects_model.Create(db, request)
		if err != nil {
			logger.Logger.Error("ProjectsController.CreateProject: Unable to create project", err.Error())
			rw.Error("Something went wrong, Please try again", http.StatusUnprocessableEntity)
			return
		}
		// create project members
		project_members_model.Create(db, models.ProjectMembers{
			ProjectID: result.ID,
			UserID:    authenticatedUser.ID,
			IsOwner:   true,
		})
		rw.Success(result, http.StatusOK)
		return
	}
}

func (c *ProjectsController) GetProject(db *database.DB) http.HandlerFunc {
	/*
		TODO - Add pagination
	*/
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		authenticatedUser, err := rw.User(db, r)
		if err != nil {
			logger.Logger.Error("ProjectsController.GetProject: ERROR", err)
			rw.Forbidden()
			return
		}

		// Fetch projects from `project_members` table to also
		// fetch the projects that the user is a member of.
		// result, err := models.GetByOwnerId(db, authenticatedUser.ID)

		result, err := project_members_model.GetProjectsByUserId(db, authenticatedUser.ID)
		if err != nil {
			logger.Logger.Error("ProjectsController.GetProject: ERROR", err)
			rw.Error("Unable to fetch data", http.StatusInternalServerError)
			return
		}
		rw.Success(result, http.StatusOK)
		return
	}
}

// @deprecated - in favor of using 'cross-account-access' over 'access_keys'
// Allow users to create more projects based on subscription plan.
func (c *ProjectsController) _createWithService(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		authenticatedUser, err := rw.User(db, r)
		if err != nil {
			logger.Logger.Error("ProjectsController.createWithService: ERROR", err)
			if errors.Is(err, custom_errors.NoDataFound) {
				rw.Unauthorized()
				return
			}
			rw.Forbidden()
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Logger.Error("ProjectsController.createWithService: Error reading request body:", err)
			rw.Error("Bad request", http.StatusUnprocessableEntity)
			return
		}
		var request projects_model.CreateWithServiceProps
		if err := json.Unmarshal(body, &request); err != nil {
			logger.Logger.Error("ProjectsController.createWithService: Error parsing request body:", err)
			rw.Error(custom_errors.UnknownError.Error(), http.StatusUnprocessableEntity)
			return
		}
		defer r.Body.Close()
		if request.Type == "" || (request.Type != "PUBLIC" && request.Type != "PRIVATE") {
			rw.Error("Required enum field `type` must be one of value ('PUBLIC', 'PRIVATE')", http.StatusBadRequest)
			return
		}
		if err := helpers.CheckEmptyFields(request); err != nil {
			rw.Error(err.Error(), http.StatusBadRequest)
			return
		}
		request.OwnerID = authenticatedUser.ID
		request.Email = authenticatedUser.Email
		result, err := projects_model.CreateWithService(db, request)
		if err != nil {
			rw.Error("Unexpected Error Occured, Please try again later", http.StatusInternalServerError)
			return
		}
		rw.Success(result, http.StatusOK)
		return
	}
}

func (c *ProjectsController) GetProjectById(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		authenticatedUser, err := rw.User(db, r)
		if err != nil {
			logger.Logger.Error("ProjectsController.GetProjectById: ERROR", err)
			rw.Forbidden()
			return
		}
		id := mux.Vars(r)["id"]
		isUUIDValid := helpers.IsValidUUID(id)
		if !isUUIDValid {
			rw.Error("Invalid project ID provided", http.StatusBadRequest)
			return
		}
		uuid, err := uuid.Parse(id)
		if err != nil {
			rw.Error("Invalid project ID of type uuid provided", http.StatusUnprocessableEntity)
			return
		}
		// Fetching projects from `project_members` table mapping.
		// Oriignal projects table only has `owner_id`.
		// result, err := models.GetByIdAndUserId(db, uuid, authenticatedUser.ID)
		result, err := project_members_model.GetProjectByIdAndUserId(db, uuid, authenticatedUser.ID)
		if err != nil {
			logger.Logger.Error("ProjectsController.GetProjectById: ERROR", err)
			rw.Error("Unable to fetch data", http.StatusInternalServerError)
			return
		}
		rw.Success(result, http.StatusOK)
		return
	}
}

type UpdateDTO struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Email       string  `json:"email"`
	IsDeleted   *bool   `json:"isDeleted,omitempty"`
}

func (c *ProjectsController) Update(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		authenticatedUser, err := rw.User(db, r)
		if err != nil {
			logger.Logger.Error("ProjectsController.Update: ERROR", err)
			rw.Forbidden()
			return
		}
		id := mux.Vars(r)["id"]
		isUUIDValid := helpers.IsValidUUID(id)
		if !isUUIDValid {
			rw.Error("Invalid project ID provided", http.StatusBadRequest)
			return
		}
		uuid, err := uuid.Parse(id)
		if err != nil {
			rw.Error("Invalid project ID of type uuid provided", http.StatusUnprocessableEntity)
			return
		}
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			logger.Logger.Error("ProjectsController.Update: Error reading request body:", err)
			rw.Error("Bad request", http.StatusUnprocessableEntity)
			return
		}
		var request UpdateDTO
		if err := json.Unmarshal(body, &request); err != nil {
			logger.Logger.Error("ProjectsController.Update: Error parsing request body:", err)
			rw.Error(custom_errors.UnknownError.Error(), http.StatusUnprocessableEntity)
			return
		}
		if err := projects_model.Update(db, uuid, authenticatedUser.ID, models.Projects{
			Name:        request.Name,
			Description: request.Description,
			Email:       request.Email,
			IsDeleted:   request.IsDeleted,
		}); err != nil {
			logger.Logger.Error("ProjectsController.Update: ERROR", err)
			rw.Error("Unable to save data", http.StatusInternalServerError)
			return
		}
		rw.Success("data saved", http.StatusOK)
		return
	}
}
