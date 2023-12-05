const isAuthorised = require("../../policies/isAuthorised");
const isWhitelist = require("../../helpers/whiteListCheck");
// eslint-disable-next-line no-unused-vars
const { Unauthorised, Unexpected } = require("../../plugins/responses");
const loggers = require("../../loggers");
const { app_service_allowed_routes } = require("../../routes/app_service");
const { COOKIE_NAME } = require("../../helpers/constants");

module.exports = {
	async onRequest(req, res) {
		try {
			let flag = isWhitelist(
				app_service_allowed_routes,
				req.method,
				req.params
			);
			loggers.info("appService.hooks.onRequest: incoming request for app service: " + JSON.stringify({
				remoteAddress: req.socket.remoteAddress,
				forwardIp: req.headers["x-forwarded-for"],
				method: req.method,
				flag,
				url: req.url,
				params: req.params,
				query: req.query,
				data: req.body,
			}));
			if (!flag) {
				const cookie = req.cookies[COOKIE_NAME] || req.headers.cookie;
				if (!cookie) {
					throw new Unauthorised();
				}
				const user = await isAuthorised(cookie);
				if (req.url === "/oauth/session") {
					if (user) {
						/**
						 * User data is being set as `data` from app backend
						 * while setting session.
						 * Cookie name set is `cv-token`
						 */
						res.send({
							success: true,
							data: user.data
						});
					}
					return res.end();
				}
			}
		} catch (err) {
			console.log("appService.hooks.onRequest: ERROR", err);
			res.send(err, err.status || 500);
			return res.end();
		}
	},
};
