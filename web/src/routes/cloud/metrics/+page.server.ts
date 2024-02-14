import { getMetricPanels } from "$src/api/metricPanels";
import { getProviderAccountDetails } from "$src/api/providerAccounts";
import { getResourceData } from "$src/api/services";
import type { ProviderAccountWithProjectProps } from "$src/customTypes/providerAccounts";
import type { ResourceDataReturnType } from "$src/customTypes/services";
import { COOKIE_NAME, DEFAULT_REGION } from "$src/helpers/constants";
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
	let panels, providerAcc;
	if (cookie) {
		try {
			const [ result, metricPanels, accountDetails ] = await Promise.all([
				getResourceData(
					{
						projectId,
						providerAccountId,
						region,
					},
					{ cookie }
				),
				getMetricPanels({
					providerAccountId,
					page: 1,
					limit: 10
				}, cookie),
				getProviderAccountDetails(providerAccountId, projectId, cookie)
			]);
			if (result.error || !result.data) throw result;
			if (metricPanels.error || !metricPanels.data) throw metricPanels;
			if (accountDetails.error || !accountDetails.data) throw accountDetails;
			res = result.data;
			panels = metricPanels.data;
			providerAcc = accountDetails.data;
		} catch (err) {
			console.error("Unable to fetch metrics:", err);
			error = "Unable to fetch data";
		}
	} else {
		console.error("cookie not found. re-login required.");
	}
	return {
		projectId,
		providerAccountId,
		region,
		metricData: res,
		error,
		showRegionDD: true,
		panels,
		providerAcc
	};
}) satisfies PageServerLoad;