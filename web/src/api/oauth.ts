import type { SessionProps } from "$src/customTypes/user";
import { requester, type ApiResponsePromise } from "$src/helpers/requester";

export const Login = async (params: {
	provider: "google" | "github";
	code: string;
}): ApiResponsePromise<SessionProps> => {
	return requester({
		url: `/oauth/${params.provider}`,
		data: { code: params.code },
		method: "POST"
	});
};

export const getUserSession = async (): ApiResponsePromise<SessionProps> => {
	return requester({
		url: "/oauth/session",
		data: {},
		method: "GET"
	});
};

export const logout = async () => {
	return requester({
		url: "/oauth/logout",
		data: {},
		method: "DELETE"
	});
};
