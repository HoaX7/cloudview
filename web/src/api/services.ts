import type { MetricDataReturnType } from "$src/customTypes/Services";
import {
	requester,
	type ApiResponsePromise,
	type RequestHeaders,
} from "$src/helpers/requester";

export const getMetricData = async (params: {
    providerAccountId: string;
    projectId: string;
    region: string;
}, headers?: { cookie?: string }): ApiResponsePromise<MetricDataReturnType> => {
	return requester({
	    url: "/services/getData",
	    data: params,
	    method: "GET",
	    headers: (headers?.cookie ? { cookie: headers?.cookie } : {}) as RequestHeaders
	});
};
