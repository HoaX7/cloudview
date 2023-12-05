const { API_PORT_HTTP } = require("./env");
const server = require("./app");

server.start(API_PORT_HTTP || 5000, "0.0.0.0").then(() => {
	console.log("Listening on port", API_PORT_HTTP);
});