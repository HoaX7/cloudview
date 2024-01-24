package router

import (
	"cloudview/app/src/api/controllers"
	"cloudview/app/src/database"

	"github.com/gorilla/mux"
)

func servicesRouter(r *mux.Router, controller *controllers.ServicesController, db *database.DB) {
	subrouter := r.PathPrefix("/services").Subrouter()

	subrouter.HandleFunc("/getData", controller.GetServiceData(db)).Methods("GET")

	// This route is for apigatewayv2 to fetch integrations
	// to map ec2/lambda
	subrouter.HandleFunc("/aws/getApiGatewayV2Integrations", controller.GetApiGatewayV2Integrations(db)).Methods("GET")

	// This route is explicitly used for aws to fetch usage/bandwidth data
	subrouter.HandleFunc("/{provider}/getUsage", controller.GetUsage(db)).Methods("GET")
}
