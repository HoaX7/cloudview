/*
This is the backend app service home router mainly used to ping the service.
*/
package router

import (
	"net/http"

	"cloudview/app/src/api/middleware"

	"github.com/gorilla/mux"
)

func homeRouter(r *mux.Router) {
	/*
		routes: This matches '/route'
		If you use / the route only matches /route/

		If you want to serve both routes, need to duplicate urls
	*/
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.RegisterResponses(w)
		rw.Success("App service", http.StatusOK)
	}).Methods("GET")
}
