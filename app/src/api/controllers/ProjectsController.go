package controllers

import (
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/api/middleware"
	"cloudview/app/src/database"
	"cloudview/app/src/helpers"
	"cloudview/app/src/models"
	project_members_model "cloudview/app/src/models/project_members"
	projects_model "cloudview/app/src/models/projects"
	"cloudview/app/src/permissions"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func (c *ProjectsController) CreateProject(db *database.DB) http.HandlerFunc {
	c.Logger.SetName(c.Name + ".CreateProject")
	/*
		TODO - Add middleware to check if the user
		has the correct subscription plan to be able to
		create more than 1 project.

		Pricing Model also includes
		cost per contributor.

	*/
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.CustomResponseWriter(w)
		authenticatedUser := rw.SessionUser
		perms := authenticatedUser.Permissions
		canContinue := false
		if perms != nil {
			canContinue = permissions.VerifyPermissions([]string{
				permissions.USER_MODIFY_PROJECT,
			}, *perms)
			canContinue = false
		}
		if !canContinue {
			rw.ErrorMessage = "Your account does not have the right permissions, please contact us at vivekrajsr.96@gmail.com to resolve the issue."
			rw.Forbidden()
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			c.Logger.Error("Error reading request body:", err)
			rw.Error("Bad request", http.StatusUnprocessableEntity)
			return
		}
		var request models.Projects
		if err := json.Unmarshal(body, &request); err != nil {
			c.Logger.Error("Error parsing request body:", err)
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
		c.Logger.Log("creating new project for owner: ", request.OwnerID)
		result, err := projects_model.Create(db, request)
		if err != nil {
			c.Logger.Error("Unable to create project", err.Error())
			rw.Error("Something went wrong, Please try again", http.StatusUnprocessableEntity)
			return
		}
		memberPerms := permissions.SetPermissions(permissions.AllProjectMemberPermissions)
		// create project members
		project_members_model.Create(db, models.ProjectMembers{
			ProjectID:   result.ID,
			UserID:      authenticatedUser.ID,
			IsOwner:     true,
			Permissions: memberPerms,
		})
		rw.Success(result, http.StatusOK)
		return
	}
}

func (c *ProjectsController) GetProject(db *database.DB) http.HandlerFunc {
	c.Logger.SetName(c.Name + ".GetProject")
	/*
		TODO - Add pagination
	*/
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.CustomResponseWriter(w)
		authenticatedUser := rw.SessionUser
		// Fetch projects from `project_members` table to also
		// fetch the projects that the user is a member of.
		// result, err := models.GetByOwnerId(db, authenticatedUser.ID)

		fmt.Printf("%v", authenticatedUser)
		result, err := project_members_model.GetProjectsByUserId(db, authenticatedUser.ID)
		if err != nil {
			c.Logger.Error("unable to fetch data", err)
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
	c.Logger.SetName(c.Name + "._createWithService")
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.CustomResponseWriter(w)
		authenticatedUser := rw.SessionUser
		body, err := io.ReadAll(r.Body)
		if err != nil {
			c.Logger.Error("Error reading request body:", err)
			rw.Error("Bad request", http.StatusUnprocessableEntity)
			return
		}
		var request projects_model.CreateWithServiceProps
		if err := json.Unmarshal(body, &request); err != nil {
			c.Logger.Error("Error parsing request body:", err)
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
	c.Logger.SetName(c.Name + ".GetProjectById")
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.CustomResponseWriter(w)
		authenticatedUser := rw.SessionUser
		id := mux.Vars(r)["id"]
		isUUIDValid, uuid := helpers.IsValidUUID(id)
		if !isUUIDValid {
			rw.Error("Invalid project ID provided", http.StatusBadRequest)
			return
		}
		// Fetching projects from `project_members` table mapping.
		// Oriignal projects table only has `owner_id`.
		// result, err := models.GetByIdAndUserId(db, uuid, authenticatedUser.ID)
		result, err := project_members_model.GetProjectByIdAndUserId(db, *uuid, authenticatedUser.ID)
		if err != nil {
			c.Logger.Error("Unable to fetch data", err)
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
	c.Logger.SetName(c.Name + ".Update")
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.CustomResponseWriter(w)
		authenticatedUser := rw.SessionUser
		perms := authenticatedUser.Permissions
		canContinue := false
		if perms != nil {
			canContinue = permissions.VerifyPermissions([]string{
				permissions.USER_MODIFY_PROJECT,
			}, *perms)
		}
		if !canContinue {
			rw.ErrorMessage = "Your account does not have the right permissions, please contact us at vivekrajsr.96@gmail.com to resolve the issue."
			rw.Forbidden()
			return
		}
		id := mux.Vars(r)["id"]
		isUUIDValid, uuid := helpers.IsValidUUID(id)
		if !isUUIDValid {
			rw.Error("Invalid project ID provided", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			c.Logger.Error("Error reading request body:", err)
			rw.Error("Bad request", http.StatusUnprocessableEntity)
			return
		}
		var request UpdateDTO
		if err := json.Unmarshal(body, &request); err != nil {
			c.Logger.Error("Error parsing request body:", err)
			rw.Error(custom_errors.UnknownError.Error(), http.StatusUnprocessableEntity)
			return
		}
		if err := projects_model.Update(db, *uuid, authenticatedUser.ID, models.Projects{
			Name:        request.Name,
			Description: request.Description,
			Email:       request.Email,
			IsDeleted:   request.IsDeleted,
		}); err != nil {
			c.Logger.Error("unable to save", err)
			rw.Error("Unable to save data", http.StatusInternalServerError)
			return
		}
		rw.Success("data saved", http.StatusOK)
		return
	}
}
