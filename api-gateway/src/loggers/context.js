const { createNamespace } = require("cls-hooked");

const LOGGER_CONTEXT = "logger_context";
/**
 * cls hooked uses call stack to map each context
 * to that call stack to maintain uniqueness.
 * So the interface set does not change with multiple
 * requests
 */
const ctx = createNamespace(LOGGER_CONTEXT);

const getLoggerContext = () =>
	ctx.get(LOGGER_CONTEXT) || {};

/**
 * 
 * @param {{trackingId: string}} context
 * @property trackingId content.trackingId
 * @returns 
 */
const setLoggerContext = (context) =>
	ctx.set(LOGGER_CONTEXT, context);

const initLoggerContext = (cb) => {
	return ctx.run(() => cb());
};

module.exports = {
	initLoggerContext,
	setLoggerContext,
	getLoggerContext
};