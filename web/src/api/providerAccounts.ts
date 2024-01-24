import type {
	ProviderAccountProps,
	ProviderAccountWithProjectProps,
} from "$src/customTypes/ProviderAccounts";
import {
	requester,
	type ApiResponsePromise,
	type RequestHeaders,
} from "$src/helpers/requester";

// TODO - Add Pagination
export const getProviderAccountsByProjectId = async (
	projectId: string,
	cookieParams: { cookie?: string }
): ApiResponsePromise<ProviderAccountProps[]> => {
	return requester({
		url: "/provider_accounts",
		data: { projectId },
		method: "GET",
		headers: (cookieParams.cookie
			? { cookie: cookieParams.cookie }
			: {}) as RequestHeaders,
	});
};

export const editProviderAccount = async (
	id: string,
	data: Partial<ProviderAccountProps>
) => {
	return requester({
		url: `/provider_accounts/${id}`,
		data,
		method: "PATCH",
	});
};

export const createProviderAccount = async (params: {
  provider: "AWS";
  projectId: string;
  accessKeyId: string;
  accessKeySecret: string;
  name: string;
}): ApiResponsePromise<ProviderAccountProps> => {
	return requester({
		url: "/provider_accounts",
		data: params,
		method: "POST",
	});
};

export const getProviderAccountDetails = async (
	id: string,
	projectId: string,
	cookie?: string
): ApiResponsePromise<ProviderAccountWithProjectProps> => {
	return requester({
		url: `/provider_accounts/${id}`,
		data: { projectId },
		method: "GET",
		headers: (cookie ? { cookie } : {}) as RequestHeaders,
	});
};
