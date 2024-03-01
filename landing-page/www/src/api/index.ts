import { requester } from "../helpers/requester";

export const requestDemoApi = async (data: { email: string; name: string; callScheduledAt: Date; notes?: string; }) => {
	return requester({
		url: "register",
		method: "POST",
		data
	});
};