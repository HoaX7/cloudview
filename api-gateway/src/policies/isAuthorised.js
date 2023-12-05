const jwt = require("jsonwebtoken");
const { JWT_KEY } = require("../../env");
const { Unauthorised } = require("../plugins/responses");

module.exports = async function (token) {
	try {
		const decoded = jwt.verify(token, JWT_KEY);
		if (!decoded) {
			throw new Unauthorised();
		}
		return decoded;
	} catch (err) {
		throw new Unauthorised();
	}
};
