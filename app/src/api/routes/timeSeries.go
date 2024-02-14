/*
This is the backend app service home router mainly used to ping the service.
*/
package router

import (
	"cloudview/app/src/api/controllers"
	"cloudview/app/src/api/middleware"
	"cloudview/app/src/database"

	"github.com/gorilla/mux"
)

func timeSeriessRouter(r *mux.Router, controller *controllers.TimeSeriesController, db *database.DB) {
	subrouter := r.PathPrefix("/timeSeries").Subrouter()

	subrouter.HandleFunc("", controller.SaveSeries(db)).Methods("POST")
	subrouter.HandleFunc("", middleware.Authenticate(controller.GetByMetricId(db), db)).Methods("GET")
}
