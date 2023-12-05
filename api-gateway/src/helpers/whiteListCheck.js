const { match } = require("path-to-regexp");
const { Unexpected } = require("../plugins/responses");
module.exports = (service, method, params = undefined) => {
	try {
		const url_list = service[method];
		if (url_list && (url_list.includes(params.wild))) {
			return true;
		} else {
			for (let pos = 0; pos < url_list.length; pos++) {
				const check = match(url_list[pos], { decode: decodeURIComponent });
				const matched_obj = check("/" + params.wild);
				if (matched_obj) {
					return true;
				}
			}
		}
		return false;
	} catch (err) {
		throw new Unexpected(err.message, err.code || err.status);
	}
	
};