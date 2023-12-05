const requestIp = require("request-ip");
const rateLimit = require("express-rate-limit");

module.exports = {
	captureIp(req, res, next) {
		req.ip = requestIp.getClientIp(req);
		return next();
	},
	rateLimitApi(req, res, next) {
		rateLimit({
			windowMs: 1000 * 60,
			max: 30,
			handler: (req, res) => res.send("You are being rate limited", 429)
		});
		return next();
	}
};