import isEmpty from "lodash/isEmpty";

export const isEmptyObject = (val = {}) => {
	return isEmpty(val);
};