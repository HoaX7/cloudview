package controllers

import (
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/api/middleware"
	"cloudview/app/src/database"
	"cloudview/app/src/helpers"
	"cloudview/app/src/models"
	project_members_model "cloudview/app/src/models/project_members"
	projects_model "cloudview/app/src/models/projects"
	users_model "cloudview/app/src/models/users"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// TODO - add pagination
func (c *ProjectMembersController) GetMembersByProjectId(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.CustomResponseWriter(w)
		projectId := r.URL.Query().Get("projectId")
		valid, uid := helpers.IsValidUUID(projectId)
		if projectId == "" || !valid {
			rw.Error("Invalid `projectId` provided", http.StatusBadRequest)
			return
		}
		projectUid := *uid
		authenticatedUser, err := rw.User(db, r)
		if err != nil {
			c.Logger.Error("ProjectMembersController.GetMembersByProjectId: ERROR", err)
			rw.Forbidden()
			return
		}
		projectData, err := projects_model.GetByIdAndUserId(db, projectUid, authenticatedUser.ID)
		if err != nil {
			if errors.Is(err, custom_errors.NoDataFound) {
				rw.Error("Please select a valid Project.", http.StatusNotFound)
				return
			}
			c.Logger.Error("ProjectMembersController.GetMembersByProjectId: ERROR", err)
			rw.Error("Unable to fetch data", http.StatusInternalServerError)
			return
		}
		result, err := project_members_model.GetMembersByProjectId(db, projectData.ID)
		if err != nil {
			if errors.Is(err, custom_errors.NoDataFound) {
				rw.Success([]interface{}{}, http.StatusOK)
				return
			}
		}
		rw.Success(result, http.StatusOK)
		return
	}
}

type CreateMemberStruct struct {
	ProjectID uuid.UUID `json:"projectId"`
	Email     string    `json:"email"`
}

type CreateMemberReturnStruct struct {
	*models.ProjectMembers
	User *users_model.Users `json:"user"`
}

func (c *ProjectMembersController) CreateMember(db *database.DB) http.HandlerFunc {
	c.Logger.SetName(c.Name + ".CreateMember")
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.CustomResponseWriter(w)
		authenticatedUser := rw.SessionUser
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			c.Logger.Error("Error reading request body:", err)
			rw.Error("Bad request", http.StatusUnprocessableEntity)
			return
		}
		var request CreateMemberStruct
		if err := json.Unmarshal(body, &request); err != nil {
			c.Logger.Error("Error parsing request body:", err)
			rw.Error(custom_errors.UnknownError.Error(), http.StatusUnprocessableEntity)
			return
		}
		if request.Email == "" {
			rw.Error("Invalid `email` provided", http.StatusBadRequest)
			return
		}
		isUUIDValid, _ := helpers.IsValidUUID(request.ProjectID.String())
		if !isUUIDValid {
			rw.Error("Invalid project ID of type uuid provided", http.StatusBadRequest)
			return
		}

		// By fetching project directly from `projects` table you can also verify if the
		// authenticated user is the project owner.
		project, err := projects_model.GetByIdAndUserId(db, request.ProjectID, authenticatedUser.ID)
		if err != nil || project.OwnerID != &authenticatedUser.ID {
			c.Logger.Error("Project not found: ERROR project owner mismatch, project:", request.ProjectID, "auth user:", authenticatedUser.ID)
			rw.Forbidden()
			return
		}
		// Fetch user from DB if not exists create the user.
		var user users_model.Users
		/*
			TODO - After a new user is created send them a welcome email.
			Notify members when they are added to a new Project Team.
		*/
		existingUser, err := users_model.GetByEmail(db, request.Email)
		if err != nil {
			if errors.Is(err, custom_errors.NoDataFound) {
				c.Logger.Log("user not found in system, creating new user...")
				data, err := users_model.Create(db, users_model.Users{
					Email:    request.Email,
					Username: request.Email,
				})
				if err != nil {
					c.Logger.Error("Unable to create user in system", err)
					rw.Error("Unable to invite member, Please try again later.", http.StatusInternalServerError)
					return
				}
				user = data
			} else {
				c.Logger.Error("Unknown ERROR", err)
				rw.Error("Unable to invite member, Please try again later.", http.StatusInternalServerError)
				return
			}
		} else {
			c.Logger.Log("user found:", existingUser.ID)
			user = existingUser
		}

		c.Logger.Log("checking if member exists in Team:", user.ID)
		// Check if the member is already part of the Team.
		projectMember, err := project_members_model.GetProjectByIdAndUserId(db, request.ProjectID, user.ID)
		if err != nil {
			// if user is not already a member send invite.
			if !errors.Is(err, custom_errors.NoDataFound) {
				c.Logger.Error("Unable to invite user", err)
				rw.Error("Unable to invite member, Please try again later.", http.StatusInternalServerError)
				return
			}
		}
		/*
			TODO - Doing a dirty check to see if member is already part of the Team
			need to check how to validate uuid correctly.
		*/
		if projectMember.ID.String() != "00000000-0000-0000-0000-000000000000" {
			c.Logger.Log("Member is already part of the Team:", user.ID)
			rw.Error("Member has already joined your Team. Upgrade your account to invite more members.", http.StatusConflict)
			return
		}
		result, err := project_members_model.Create(db, models.ProjectMembers{
			ProjectID: request.ProjectID,
			UserID:    user.ID,
			IsOwner:   false,
			// TODO - send appropriate permissions from frontend to assign users
			// Permissions: permissions.SetPermissions([]string{}),
		})
		if err != nil {
			c.Logger.Error("Unable to invite member", err)
			rw.Error("Unable to invite member, Please try again later.", http.StatusInternalServerError)
			return
		}
		c.Logger.Log("create success", result.ID)
		res := CreateMemberReturnStruct{
			&result,
			&user,
		}
		rw.Success(res, http.StatusOK)
		return
	}
}

type ToggleMemberAccessStruct struct {
	IsActive  *bool     `json:"isActive"`
	ProjectID uuid.UUID `json:"projectId"`
	IsDeleted *bool     `json:"isDeleted"`
	// IsDeleted bool `json:"isDeleted"` // discuss if you should remove the member entirely.
}

func (c *ProjectMembersController) ToggleMemberAccess(db *database.DB) http.HandlerFunc {
	c.Logger.SetName(c.Name + ".ToggleMemberAccess")
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.CustomResponseWriter(w)
		authenticatedUser := rw.SessionUser
		id := mux.Vars(r)["id"]
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			c.Logger.Error("Error reading request body:", err)
			rw.Error("Bad request", http.StatusUnprocessableEntity)
			return
		}
		var request ToggleMemberAccessStruct
		if err := json.Unmarshal(body, &request); err != nil {
			c.Logger.Error("Error parsing request body:", err)
			rw.Error(custom_errors.UnknownError.Error(), http.StatusUnprocessableEntity)
			return
		}
		if id == "" {
			rw.NotFound()
			return
		}
		valid, uid := helpers.IsValidUUID(id)
		if !valid {
			rw.Error("Invalid ID provided", http.StatusBadRequest)
			return
		}
		isUUIDValid, _ := helpers.IsValidUUID(request.ProjectID.String())
		if !isUUIDValid {
			rw.Error("Invalid project ID of type uuid provided", http.StatusBadRequest)
			return
		}

		// By fetching project directly from `projects` table you can also verify if the
		// authenticated user is the project owner.
		project, err := projects_model.GetByIdAndUserId(db, request.ProjectID, authenticatedUser.ID)
		if err != nil || *project.OwnerID != authenticatedUser.ID {
			c.Logger.Error("Project not found: ERROR project owner mismatch, project:", request.ProjectID, "auth user:", authenticatedUser.ID)
			rw.Forbidden()
			return
		}
		if err := project_members_model.Update(db, *uid, models.ProjectMembers{
			IsActive:  request.IsActive,
			IsDeleted: request.IsDeleted,
		}); err != nil {
			c.Logger.Error("unable to update data", err)
			rw.Error("Unable to update data", http.StatusInternalServerError)
			return
		}
		rw.Success("Data updated", http.StatusOK)
		return
	}
}
