import type { ProjectMemberApiProps } from "$src/customTypes/projectMembers";
import {
	requester,
	type ApiResponsePromise,
	type RequestHeaders,
} from "$src/helpers/requester";

export const getProjectMembers = async (
	params: { projectId: string },
	headers?: { cookie?: string }
): ApiResponsePromise<ProjectMemberApiProps[]> => {
	return requester({
		url: "/projectMembers",
		data: params,
		method: "GET",
		headers: (headers?.cookie
			? { cookie: headers.cookie }
			: {}) as RequestHeaders,
	});
};

export const inviteProjectMember = async (data: {
    projectId: string;
    email: string;
}): ApiResponsePromise<ProjectMemberApiProps> => {
	return requester({
		url: "/projectMembers",
		data,
		method: "POST"
	});
};

export const updateProjectMember = async (params: { id: string; }, data: {
    isActive?: boolean;
    isDeleted?: boolean;
    projectId: string;
}) => {
	return requester({
		url: `/projectMembers/${params.id}`,
		data,
		method: "PATCH"
	});
};