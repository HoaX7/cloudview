const { app_service } = require("./src/configs/appService/index");
var cookieParser = require("cookie-parser");
const bodyParser = require("body-parser");
const gateway = require("fast-gateway");
const { captureIp, rateLimitApi } = require("./src/helpers/rateLimiter");
const { initLoggerContext, setLoggerContext } = require("./src/loggers/context");
const loggers = require("./src/loggers");

const allowedDomains = (process.env.CORS_WHITELIST || "").split(",");

const server = gateway({
	middlewares: [
		require("cors")({
			origin: allowedDomains,
			credentials: true
		}),
		cookieParser(),
		bodyParser.raw({
			type: "multipart/form-data",
			limit: 10e6,
		}),
		bodyParser.json(),
		bodyParser.urlencoded({ extended: true }),
		captureIp,
		rateLimitApi,
		(req, res, next) => {
			loggers.info("Incoming request");
			initLoggerContext(() => {
				setLoggerContext({ trackingId: new Date().getTime().toString() });
				next();
			});
		}
	],
	routes: [ app_service ],
	restana: {
		errorHandler(err, req, res) {
			res.send(
				{
					error: true,
					message: err.message,
				},
				err.status || 300
			);
			return res.end();
		},
	},
});

server.get("/", (req, res) => {
	return res.send({
		success: true,
		message: "Api gateway service"
	});
});

module.exports = server;
