const { APP_SERVICE_DOMAIN_PROTOCOL, APP_SERVICE_DOMAIN } = require("../../../env");
const loggers = require("../../loggers");
const hooks = require("./hooks");

const target = `${APP_SERVICE_DOMAIN_PROTOCOL}://${APP_SERVICE_DOMAIN}`;
loggers.info("init_app_service:", target);
const app_service = {
	prefix: "/api/app/v1",
	target,
	hooks
};
module.exports = { app_service };
