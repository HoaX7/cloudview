/*
This is the backend app service home router mainly used to ping the service.
*/
package router

import (
	"cloudview/app/src/api/controllers"
	"cloudview/app/src/database"

	"github.com/gorilla/mux"
)

func metricPanelsRouter(r *mux.Router, controller *controllers.MetricPanelsController, db *database.DB) {
	subrouter := r.PathPrefix("/metricPanels").Subrouter()
	/*
		routes: This matches '/route'
		If you use / the route only matches /route/

		If you want to serve both routes, need to duplicate urls
	*/
	subrouter.HandleFunc("", controller.CreateMetricPanel(db)).Methods("POST")
	subrouter.HandleFunc("/{id}", controller.UpdateMetricPanel(db)).Methods("PATCH")
}
