import type { CreateMetricPanelProps, MetricPanelProps } from "$src/customTypes/metricPanels";
import { requester, type ApiResponsePromise, type RequestHeaders } from "$src/helpers/requester";

type Params = {
    page: number;
    limit: number;
    providerAccountId: string;
}
export const getMetricPanels = async (params: Params, cookie?: string) => {
	return requester({
		method: "GET",
		url: "/metricPanels",
		data: params,
		headers: (cookie ? { cookie } : {}) as RequestHeaders
	});
};

export const createMetricPanel = async (data: CreateMetricPanelProps): ApiResponsePromise<MetricPanelProps> => {
	return requester({
		method: "POST",
		url: "/metricPanels",
		data,
	});
};