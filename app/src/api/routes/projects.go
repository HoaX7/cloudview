package router

import (
	"cloudview/app/src/api/controllers"
	"cloudview/app/src/api/middleware"
	"cloudview/app/src/database"

	"github.com/gorilla/mux"
)

func projectsRouter(r *mux.Router, controller *controllers.ProjectsController, db *database.DB) {
	subrouter := r.PathPrefix("/projects").Subrouter()

	subrouter.HandleFunc("", middleware.Authenticate(controller.CreateProject(db), db)).Methods("POST")
	subrouter.HandleFunc("", middleware.Authenticate(controller.GetProject(db), db)).Methods("GET")
	subrouter.HandleFunc("/{id}", middleware.Authenticate(controller.GetProjectById(db), db)).Methods("GET")
	subrouter.HandleFunc("/{id}", middleware.Authenticate(controller.Update(db), db)).Methods("PATCH")
}
