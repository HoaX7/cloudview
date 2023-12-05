const loggers = {
	info: (...args) => {
		console.log(args.join(" -> "));
	},
	error: (...args) => {
		console.log(args.join(" -> "));
	}
};

module.exports = loggers;