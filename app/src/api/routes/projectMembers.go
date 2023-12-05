package router

import (
	"cloudview/app/src/api/controllers"
	"cloudview/app/src/database"

	"github.com/gorilla/mux"
)

func projectMembersRouter(r *mux.Router, controller *controllers.ProjectMembersController, db *database.DB) {
	subrouter := r.PathPrefix("/projectMembers").Subrouter()

	subrouter.HandleFunc("", controller.GetMembersByProjectId(db)).Methods("GET")
	subrouter.HandleFunc("", controller.CreateMember(db)).Methods("POST")
	subrouter.HandleFunc("/{id}", controller.ToggleMemberAccess(db)).Methods("PATCH")
}
