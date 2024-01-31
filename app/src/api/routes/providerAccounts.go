package router

import (
	"cloudview/app/src/api/controllers"
	"cloudview/app/src/api/middleware"
	"cloudview/app/src/database"

	"github.com/gorilla/mux"
)

func providerAccountsRouter(r *mux.Router, controller *controllers.ProviderAccountsController, db *database.DB) {
	subrouter := r.PathPrefix("/provider_accounts").Subrouter()

	// route will be deprecated - in favor of using 'cross account access'
	// which is more secure.
	subrouter.HandleFunc("", middleware.Authenticate(controller.StoreAccessKey(db), db)).Methods("POST")
	subrouter.HandleFunc("/{id}", middleware.Authenticate(controller.UpdateProviderAccount(db), db)).Methods("PATCH")
	subrouter.HandleFunc("", middleware.Authenticate(controller.GetByProject(db), db)).Methods("GET")
	subrouter.HandleFunc("/{id}", middleware.Authenticate(controller.GetById(db), db)).Methods("GET")
}
