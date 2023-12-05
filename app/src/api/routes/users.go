package router

import (
	"cloudview/app/src/api/controllers"
	"cloudview/app/src/database"

	"github.com/gorilla/mux"
)

func usersRouter(r *mux.Router, controller *controllers.UsersController, db *database.DB) {
	subrouter := r.PathPrefix("/users").Subrouter()

	subrouter.HandleFunc("", controller.CreateUser(db)).Methods("POST")
	subrouter.HandleFunc("", controller.GetUserByEmail(db)).Methods("GET")

}
