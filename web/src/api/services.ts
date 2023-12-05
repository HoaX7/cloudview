import type { MetricDataReturnType, ServiceProps } from "$src/customTypes/Services";
import {
	requester,
	type ApiResponsePromise,
	type RequestHeaders,
} from "$src/helpers/requester";

export const getServicesByProjectId = async (
	projectId: string,
	cookieParams: { cookie?: string }
): ApiResponsePromise<ServiceProps[]> => {
	return requester({
		url: "/services",
		data: { projectId },
		method: "GET",
		headers: (cookieParams.cookie
			? { cookie: cookieParams.cookie }
			: {}) as RequestHeaders,
	});
};


export const editService = async (id: string, data: Partial<ServiceProps>) => {
	return requester({
		url: `/services/${id}`,
		data,
		method: "PATCH"
	});
};

export const getMetricData = async (params: {
    serviceId: string;
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

export const createService = async (params: {
    provider: "AWS";
    projectId: string;
    accessKeyId: string;
    accessKeySecret: string;
    name: string;
}): ApiResponsePromise<ServiceProps> => {
	return requester({
		url: "/services",
		data: params,
		method: "POST"
	});
};
