import { getApiGatewayIntegrations } from "$src/api/aws";
import { getMetricData } from "$src/api/services";
import type {
	ApiGatewayV2Props,
	ApiGatewayWithIntegrationProps,
	MetricDataReturnType,
} from "$src/customTypes/Services";
import {
	AWS_SERVICES,
	COOKIE_NAME,
	DEFAULT_REGION,
} from "$src/helpers/constants";
import { stripIntegrationUriIp } from "$src/helpers/konva";
import type { PageServerLoad } from "./$types";

export const load = (async ({ url, cookies, params }) => {
	const serviceId = url.searchParams.get("serviceId");
	const region = url.searchParams.get("region") || DEFAULT_REGION;
	const projectId = url.searchParams.get("projectId");
	if (!projectId || !serviceId) {
		return { missingParams: true };
	}
	const cookie = cookies.get(COOKIE_NAME);
	let res: MetricDataReturnType = [];
	let error = "";
	try {
		const result = await getMetricData(
			{
				projectId,
				serviceId,
				region,
			},
			{ cookie }
		);
		if (result.error || !result.data) throw result;
		res = result.data;
		const idex = res.findIndex((r) => r.name === AWS_SERVICES.APIGATEWAYV2);
		if (idex >= 0) {
			const apigateway = res[idex];
			// Fetch integrations and attach it to apigateway data
			// This also makes it easier and faster to build arrow connectors
			// to show on canvas
			const apiGatewayWithIntegrations = await Promise.all(
				apigateway.result.Items.map((item: ApiGatewayV2Props["Items"][0]) => {
					return getApiGatewayIntegrations(
						{
							projectId,
							serviceId,
							region,
							apiId: item.ApiId,
						},
						{ cookie }
					).then((res) => {
						return {
							...item,
							integrations: res.data?.Items || [],

							/**
                             * Strip the integrationUri to get the ip address
                             * of Ec2 VMs and lambda functions.
                             *
                             * We will be using this as our target to draw connecting arrows.
                             */
							lineTargets: (res.data?.Items || []).map((it) => ({
								from: item.ApiId,
								to: stripIntegrationUriIp(it.IntegrationUri),
							})),
						};
					});
				})
			);
			res[idex].result =
        apiGatewayWithIntegrations as ApiGatewayWithIntegrationProps[];
		}
	} catch (err) {
		console.error("Unable to fetch services:", err);
		error = "Unable to fetch data";
	}
	return {
		projectId,
		serviceId,
		region,
		metricData: res,
		error,
		cloudProvider: params.provider,
	};
}) satisfies PageServerLoad;
