import type { ResourceDataReturnType } from "$src/customTypes/services";
import {
	requester,
	type ApiResponsePromise,
	type RequestHeaders,
} from "$src/helpers/requester";

export const getResourceData = async (params: {
    providerAccountId: string;
    projectId: string;
    region: string;
}, headers?: { cookie?: string }): ApiResponsePromise<ResourceDataReturnType> => {
	return requester({
	    url: "/services/getData",
	    data: params,
	    method: "GET",
	    headers: (headers?.cookie ? { cookie: headers?.cookie } : {}) as RequestHeaders
	});
};
