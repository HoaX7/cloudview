package controllers

/*
Responsible for defining all the controller methods available.
Make sure to add all the controller pointers to the `Controller` struct.
By doing this, it allows us to access these controllers from
`APIServer` struct.

## All functions/methods used in routes or another file that want to
## use DB MUST use controller methods and MUST NOT directly access models.
*/
type AuthController struct{}
type UsersController struct{}
type ProjectsController struct{}
type ServiceController struct{}
type ProjectMembersController struct{}

type Controller struct {
	*AuthController
	*UsersController
	*ProjectsController
	*ServiceController
	*ProjectMembersController
}
