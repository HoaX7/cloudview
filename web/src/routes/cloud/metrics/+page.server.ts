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
	let providerAcc;
	try {
		const [ result, providerAccDetails ] = await Promise.all([
			getResourceData(
				{
					projectId,
					providerAccountId,
					region,
				},
				{ cookie }
			),
			getProviderAccountDetails(providerAccountId, projectId, cookie)
		]);
		if (result.error || !result.data) throw result;
		if (providerAccDetails.error || !providerAccDetails.data) throw providerAccDetails;
		res = result.data;
		providerAcc = providerAccDetails.data;
	} catch (err) {
		console.error("Unable to fetch metrics:", err);
		error = "Unable to fetch data";
	}
	return {
		projectId,
		providerAccountId,
		region,
		metricData: res,
		error,
		showRegionDD: true,
		providerAcc
	};
}) satisfies PageServerLoad;