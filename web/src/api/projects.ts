import type { ProjectProps, ProjectWithServiceProps } from "$src/customTypes/projects";
import { requester, type ApiResponsePromise, type RequestHeaders } from "$src/helpers/requester";

export const getProjects = (params: { cookie?: string; }): ApiResponsePromise<ProjectProps[]> => {
	return requester({
		url: "/projects",
		data: {},
		method: "GET",
		headers: (params.cookie ? { cookie: params.cookie } : {}) as RequestHeaders
	});
};

export const getProjectById = async (params: { id: string; }) => {
	return requester({
		url: `/projects/${params.id}`,
		data: {},
		method: "GET"
	});
};

export const createProject = (
	data: Pick<ProjectProps, "name" | "description" | "type">
): ApiResponsePromise<ProjectProps> => {
	return requester({
		url: "/projects",
		data,
		method: "POST"
	});
};

export const updateProject = (
	id: string,
	data: Partial<Pick<ProjectProps, "name" | "description" | "email" | "isDeleted">>
): ApiResponsePromise<ProjectProps> => {
	return requester({
		url: `/projects/${id}`,
		data,
		method: "PATCH"
	});
};

type Params = {
    name: string;
    description?: string;
    accessKeyId: string;
    accessKeySecret: string;
    provider: string;
}
export const createProjectWithService = (data: Params): ApiResponsePromise<ProjectWithServiceProps> => {
	return requester({
		url: "/projects/createWithService",
		data,
		method: "POST"
	});
};
