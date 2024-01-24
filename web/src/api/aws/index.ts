import type { ApiGatewayV2IntegrationProps, UsageProps } from "$src/customTypes/Services";
import {
	requester,
	type ApiResponsePromise,
	type RequestHeaders,
} from "$src/helpers/requester";

export const getApiGatewayIntegrations = async (
	params: {
    projectId: string;
    providerAccountId: string;
    region: string;
    apiId: string;
  },
	headers?: { cookie?: string }
): ApiResponsePromise<ApiGatewayV2IntegrationProps> => {
	return requester({
		url: "/services/aws/getApiGatewayV2Integrations",
		data: params,
		method: "GET",
		headers: (headers?.cookie
			? { cookie: headers.cookie }
			: {}) as RequestHeaders,
	});
};

export const getEc2UsageData = async (params: {
  instanceId: string;
  projectId: string;
  providerAccountId: string;
  region: string;
  instance: "ec2";
}): ApiResponsePromise<UsageProps> => {
	return requester({
		url: "/services/aws/getUsage",
		data: params,
		method: "GET"
	});
};
