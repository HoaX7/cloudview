package router

import (
	"cloudview/app/src/api/controllers"
	"cloudview/app/src/database"

	"github.com/gorilla/mux"
)

/*
Creating a subrouter acts as a prefix to the url.
Example: `/oauth/google` or `/oauth/logout`
Notice how you do not have to specify `oauth` while initialising different routes,
it is automatically taken care by subrouter.
*/
func oauthRouter(r *mux.Router, controller *controllers.AuthController, db *database.DB) {
	subrouter := r.PathPrefix("/oauth").Subrouter()

	subrouter.HandleFunc("/{provider}", controller.OauthLogin(db)).Methods("POST")
	subrouter.HandleFunc("/logout", controller.Logout).Methods("DELETE")
}
