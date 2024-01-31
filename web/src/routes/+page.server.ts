import { getProjects } from "$src/api/projects";
import type { ProjectProps } from "$src/customTypes/projects";
import { COOKIE_NAME } from "$src/helpers/constants";
import type { PageServerLoad } from "./$types";

export const load = (async ({ cookies }) => {
	let result: ProjectProps[] = [];
	let error = "";
	try {
		const cookie = cookies.get(COOKIE_NAME);
		const res = await getProjects({ cookie });
		if (res.error || !res.data) throw res;
		result = res.data;
	} catch (err) {
		console.error("Unable to fetch projects", err);
		error = "Unable to fetch projects";
	}
	return {
		error,
		result
	};
}) satisfies PageServerLoad;