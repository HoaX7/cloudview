const { APP_SERVICE_DOMAIN_PROTOCOL, APP_SERVICE_DOMAIN } = require("../../../env");
const hooks = require("./hooks");
const app_service = {
	prefix: "/api/app/v1",
	target: `${APP_SERVICE_DOMAIN_PROTOCOL}://${APP_SERVICE_DOMAIN}`,
	hooks
};
module.exports = { app_service };
