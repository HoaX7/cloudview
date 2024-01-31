package router

import (
	"cloudview/app/src/api/controllers"
	"cloudview/app/src/api/middleware"
	"cloudview/app/src/database"

	"github.com/gorilla/mux"
)

func servicesRouter(r *mux.Router, controller *controllers.ServicesController, db *database.DB) {
	subrouter := r.PathPrefix("/services").Subrouter()

	subrouter.HandleFunc("/getData", middleware.Authenticate(controller.GetServiceData(db), db)).Methods("GET")

	// This route is for apigatewayv2 to fetch integrations
	// to map ec2/lambda
	subrouter.HandleFunc("/aws/getApiGatewayV2Integrations", middleware.Authenticate(controller.GetApiGatewayV2Integrations(db), db)).Methods("GET")

	// This route is explicitly used for aws to fetch usage/bandwidth data
	subrouter.HandleFunc("/{provider}/getUsage", middleware.Authenticate(controller.GetUsage(db), db)).Methods("GET")
}
