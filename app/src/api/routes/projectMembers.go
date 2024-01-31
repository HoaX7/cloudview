package router

import (
	"cloudview/app/src/api/controllers"
	"cloudview/app/src/api/middleware"
	"cloudview/app/src/database"

	"github.com/gorilla/mux"
)

func projectMembersRouter(r *mux.Router, controller *controllers.ProjectMembersController, db *database.DB) {
	subrouter := r.PathPrefix("/projectMembers").Subrouter()

	subrouter.HandleFunc("", middleware.Authenticate(controller.GetMembersByProjectId(db), db)).Methods("GET")
	subrouter.HandleFunc("", middleware.Authenticate(controller.CreateMember(db), db)).Methods("POST")
	subrouter.HandleFunc("/{id}", middleware.Authenticate(controller.ToggleMemberAccess(db), db)).Methods("PATCH")
}
