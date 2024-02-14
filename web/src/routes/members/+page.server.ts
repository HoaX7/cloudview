import { getProjectMembers } from "$src/api/projectMembers";
import type { ProjectMemberApiProps } from "$src/customTypes/projectMembers";
import { COOKIE_NAME } from "$src/helpers/constants";
import type { PageServerLoad } from "./$types";

export const load = (async ({ url, cookies }) => {
	const projectId = url.searchParams.get("projectId");
	const cookie = cookies.get(COOKIE_NAME);
	let result: ProjectMemberApiProps[] = [];
	let error = "";
	if (cookie) {
		try {
			if (projectId) {
				const res = await getProjectMembers({ projectId }, { cookie });
				if (res.error || !res.data) throw res;
				result = res.data;
			}
		} catch (err: any) {
			console.log("Unable to fetch members", err);
			error = err.message;
		}
	}
	return {
		projectId,
		result,
		error
	};
}) satisfies PageServerLoad;