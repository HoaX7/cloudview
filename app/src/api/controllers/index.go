package controllers

import "cloudview/app/src/api/middleware/logger"

/*
Responsible for defining all the controller methods available.
Make sure to add all the controller pointers to the `Controller` struct.
By doing this, it allows us to access these controllers from
`APIServer` struct.

## All functions/methods used in routes or another file that want to
## use DB MUST use controller methods and MUST NOT directly access models.
*/
type AuthController struct {
	Name   string
	Logger *logger.LoggerStruct
}
type UsersController struct {
	Name   string
	Logger *logger.LoggerStruct
}
type ProjectsController struct {
	Name   string
	Logger *logger.LoggerStruct
}
type ProviderAccountsController struct {
	Name   string
	Logger *logger.LoggerStruct
}
type ProjectMembersController struct {
	Name   string
	Logger *logger.LoggerStruct
}
type ServicesController struct {
	Name   string
	Logger *logger.LoggerStruct
}
type MetricPanelsController struct {
	Name   string
	Logger *logger.LoggerStruct
}
type TimeSeriesController struct {
	Name   string
	Logger *logger.LoggerStruct
}

type Controller struct {
	*AuthController
	*UsersController
	*ProjectsController
	*ProviderAccountsController
	*ProjectMembersController
	*ServicesController
	*MetricPanelsController
	*TimeSeriesController
}

func InitControllers() *Controller {
	return &Controller{
		AuthController: &AuthController{
			Name:   "AuthController",
			Logger: logger.NewLogger(),
		},
		UsersController: &UsersController{
			Name:   "UsersController",
			Logger: logger.NewLogger(),
		},
		ProjectsController: &ProjectsController{
			Name:   "ProjectsController",
			Logger: logger.NewLogger(),
		},
		ProviderAccountsController: &ProviderAccountsController{
			Name:   "ProviderAccountsController",
			Logger: logger.NewLogger(),
		},
		ProjectMembersController: &ProjectMembersController{
			Name:   "ProjectMembersController",
			Logger: logger.NewLogger(),
		},
		ServicesController: &ServicesController{
			Name:   "ServicesController",
			Logger: logger.NewLogger(),
		},
		MetricPanelsController: &MetricPanelsController{
			Name:   "MetricPanelsController",
			Logger: logger.NewLogger(),
		},
		TimeSeriesController: &TimeSeriesController{
			Name:   "TimeSeriesController",
			Logger: logger.NewLogger(),
		},
	}
}
