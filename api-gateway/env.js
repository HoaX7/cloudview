require("dotenv").config({ path: ".env" });
module.exports = {
	ENV: process.env.ENV,
	API_PORT_HTTP: process.env.API_PORT_HTTP,
	JWT_KEY: "yOushaLlnotPass",
	APP_SERVICE_DOMAIN: process.env.APP_SERVICE_DOMAIN,
	APP_SERVICE_DOMAIN_PROTOCOL: process.env.APP_SERVICE_DOMAIN_PROTOCOL
};