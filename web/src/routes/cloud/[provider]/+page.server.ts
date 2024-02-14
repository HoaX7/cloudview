import { getApiGatewayIntegrations } from "$src/api/aws";
import { getResourceData } from "$src/api/services";
import type {
	ApiGatewayV2Props,
	ApiGatewayWithIntegrationProps,
	ResourceDataReturnType,
} from "$src/customTypes/services";
import {
	AWS_SERVICES,
	COOKIE_NAME,
	DEFAULT_REGION,
} from "$src/helpers/constants";
import { stripIntegrationUriIp } from "$src/helpers/konva";
import type { PageServerLoad } from "./$types";

export const load = (async ({ url, cookies, params }) => {
	const providerAccountId = url.searchParams.get("providerAccountId");
	const region = url.searchParams.get("region") || DEFAULT_REGION;
	const projectId = url.searchParams.get("projectId");
	if (!projectId || !providerAccountId) {
		return { missingParams: true };
	}
	const cookie = cookies.get(COOKIE_NAME);
	let res: ResourceDataReturnType = [];
	let error = "";
	if (cookie) {
		try {
			const result = await getResourceData(
				{
					projectId,
					providerAccountId,
					region,
				},
				{ cookie }
			);
			if (result.error || !result.data) throw result;
			res = result.data;
			const idex = res.findIndex((r) => r.name === AWS_SERVICES.APIGATEWAYV2);
			if (idex >= 0) {
				const apigateway = res[idex] as any;
				// Fetch integrations and attach it to apigateway data
				// This also makes it easier and faster to build arrow connectors
				// to show on canvas
				const apiGatewayWithIntegrations = await Promise.all(
					apigateway.result.Items.map((item: ApiGatewayV2Props["Items"][0]) => {
						return getApiGatewayIntegrations(
							{
								projectId,
								providerAccountId,
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
	}
	return {
		projectId,
		providerAccountId,
		region,
		metricData: res,
		error,
		cloudProvider: params.provider,
		showRegionDD: true,
	};
}) satisfies PageServerLoad;
