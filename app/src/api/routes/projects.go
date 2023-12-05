package router

import (
	"cloudview/app/src/api/controllers"
	"cloudview/app/src/database"

	"github.com/gorilla/mux"
)

func projectsRouter(r *mux.Router, controller *controllers.ProjectsController, db *database.DB) {
	subrouter := r.PathPrefix("/projects").Subrouter()

	subrouter.HandleFunc("", controller.CreateProject(db)).Methods("POST")
	subrouter.HandleFunc("", controller.GetProject(db)).Methods("GET")
	subrouter.HandleFunc("/{id}", controller.GetProjectById(db)).Methods("GET")
	subrouter.HandleFunc("/createWithService", controller.CreateWithService(db)).Methods("POST")
	subrouter.HandleFunc("/{id}", controller.Update(db)).Methods("PATCH")
}
