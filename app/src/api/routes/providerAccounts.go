package router

import (
	"cloudview/app/src/api/controllers"
	"cloudview/app/src/database"

	"github.com/gorilla/mux"
)

func providerAccountsRouter(r *mux.Router, controller *controllers.ProviderAccountsController, db *database.DB) {
	subrouter := r.PathPrefix("/provider_accounts").Subrouter()

	// route will be deprecated - in favor of using 'cross account access'
	// which is more secure.
	subrouter.HandleFunc("", controller.StoreAccessKey(db)).Methods("POST")
	subrouter.HandleFunc("/{id}", controller.UpdateProviderAccount(db)).Methods("PATCH")
	subrouter.HandleFunc("", controller.GetByProject(db)).Methods("GET")
	subrouter.HandleFunc("/{id}", controller.GetById(db)).Methods("GET")
}
