const app_service_allowed_routes = {
	GET: [],
	POST: [ "/oauth/github" ],
	PATCH: [],
	DELETE: [ "/oauth/logout" ],
	PUT: []
};

module.exports = { app_service_allowed_routes };
