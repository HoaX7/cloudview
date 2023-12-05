import { getServicesByProjectId } from "$src/api/services";
import type { ServiceProps } from "$src/customTypes/Services";
import { COOKIE_NAME } from "$src/helpers/constants";
import type { PageServerLoad } from "./$types";

export const load = (async ({ url, params, cookies }) => {
	const serviceId = url.searchParams.get("serviceId");
	const projectId = params.id;
	const cookie = cookies.get(COOKIE_NAME);
	let services: ServiceProps[] = [];
	let error = "";
	try {
		const result = await getServicesByProjectId(params.id, { cookie });
		if (result.data) services = result.data;
	} catch (err) {
		console.error("Unable to fetch services:", err);
		error = "Unable to fetch data";
	}
	return {
		projectId,
		services,
		serviceId,
		error
	};
}) satisfies PageServerLoad;