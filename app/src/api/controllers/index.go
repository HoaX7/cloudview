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
type ProviderAccountsController struct{}
type ProjectMembersController struct{}
type ServicesController struct{}
type MetricPanelsController struct{}

type Controller struct {
	*AuthController
	*UsersController
	*ProjectsController
	*ProviderAccountsController
	*ProjectMembersController
	*ServicesController
	*MetricPanelsController
}

func (m *AuthController) Name() string {
	return "AuthController"
}
func (m *UsersController) Name() string {
	return "UsersController"
}
func (m *ProjectsController) Name() string {
	return "ProjectsController"
}
func (m *ProviderAccountsController) Name() string {
	return "ProviderAccountsController"
}
func (m *ProjectMembersController) Name() string {
	return "ProjectMembersController"
}
func (m *ServicesController) Name() string {
	return "ServicesController"
}
func (m *MetricPanelsController) Name() string {
	return "MetricPanelsController"
}
