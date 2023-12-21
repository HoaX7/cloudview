import { requester } from "../helpers/requester";

export const requestDemoApi = async (data: { email: string; name: string; }) => {
	return requester({
		url: "register",
		method: "POST",
		data
	});
};