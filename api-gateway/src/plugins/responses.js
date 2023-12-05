class Unauthorised {
	constructor() {
		this.status = 401;
		this.error = true;
		this.message = "Unauthorised";
	}
}

class Unexpected {
	constructor(message, status) {
		this.status = status || 500;
		this.error = true;
		this.message = message;
	}
}

module.exports = {
	Unauthorised,
	Unexpected 
};