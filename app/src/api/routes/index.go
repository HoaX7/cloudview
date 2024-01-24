package router

import (
	"cloudview/app/src/api/controllers"
	"cloudview/app/src/database"

	"github.com/gorilla/mux"
)

type APIServerInterface interface {
	GetControllers() controllers.Controller
	GetDB() *database.DB
}

/*
This package is responsible to handle all routes.
Importing all routes under `index.go` for clean code and readability.
*/
func RegisterRoutes(s APIServerInterface) *mux.Router {
	controller := s.GetControllers()
	db := s.GetDB()
	router := mux.NewRouter()
	// subrouter := router.PathPrefix("/v1").Subrouter()
	homeRouter(router)

	// Database needs to be passed down from the global struct
	// to be able to access them in controllers
	oauthRouter(router, controller.AuthController, db)
	usersRouter(router, controller.UsersController, db)
	projectsRouter(router, controller.ProjectsController, db)
	providerAccountsRouter(router, controller.ProviderAccountsController, db)
	projectMembersRouter(router, controller.ProjectMembersController, db)
	servicesRouter(router, controller.ServicesController, db)

	return router
}
