import { getProviderAccountsByProjectId } from "$src/api/providerAccounts";
import type { ProviderAccountProps } from "$src/customTypes/providerAccounts";
import { COOKIE_NAME } from "$src/helpers/constants";
import type { PageServerLoad } from "./$types";

export const load = (async ({ url, params, cookies }) => {
	const providerAccountId = url.searchParams.get("providerAccountId");
	const projectId = params.id;
	const cookie = cookies.get(COOKIE_NAME);
	let accounts: ProviderAccountProps[] = [];
	let error = "";
	try {
		if (!cookie) {
			console.error("No cookie found. re-login required.");
			return {
				projectId,
				accounts,
				providerAccountId,
				error
			};
		}
		const result = await getProviderAccountsByProjectId(params.id, { cookie });
		if (result.data) accounts = result.data;
	} catch (err) {
		console.error("Unable to fetch accounts:", err);
		error = "Unable to fetch data";
	}
	return {
		projectId,
		accounts,
		providerAccountId,
		error
	};
}) satisfies PageServerLoad;