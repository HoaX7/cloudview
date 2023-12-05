import { PUBLIC_API_SERVICE } from "$env/static/public";
import { logout } from "$src/api/oauth";
import Auth from "$src/store/auth";
import Datastore from "$src/store/data";
import axios from "axios";
import type {
	AxiosRequestConfig,
	AxiosRequestHeaders,
	AxiosResponse,
} from "axios";

export type PaginationMetadata = {
    current_page: number;
    per_page: number;
    total_count: number;
    total_pages: number;
};
export type ApiResponse<T> = {
    success?: boolean;
    error?: boolean;
    data?: T;
    metadata?: {
        pagination?: PaginationMetadata;
		expiresIn?: {
			time: number;
			unit: string;
		};
    };
};

export type ApiResponsePromise<T> = Promise<ApiResponse<T>>;

export type RequestHeaders = AxiosRequestHeaders & { cookie?: string; };
type ResponseType = "RAW" | "API_RESPONSE";
export interface RequestParams<D, RT> {
    url: string;
    method: "GET" | "PUT" | "POST" | "DELETE" | "OPTIONS" | "PATCH";
    data: D;
    service?: "app";
    version?: "v1" | "v2";
    isFile?: boolean;
    headers?: RequestHeaders;
    timeout?: number;
    responseType?: RT;
}

const isServer = typeof window === "undefined";

type ObjectType<T, U> = T extends "RAW"
    ? AxiosResponse<ApiResponse<U>>
    : ApiResponse<U>;

const requester = async <D, T, RT extends ResponseType = "API_RESPONSE">({
	data = {} as D,
	url,
	method,
	headers,
	isFile,
	timeout = 60000,
	responseType,
	service = "app",
	version = "v1",
}: RequestParams<D, RT>): Promise<ObjectType<RT, T>> => {
	const request = {
		headers: {
			"content-type": "application/json",
			...headers,
		},
		timeout,
		baseURL: `${PUBLIC_API_SERVICE}/api/${service}/${version}`,
		responseType: "json",
		withCredentials: true,
	} as AxiosRequestConfig;
	const Requester = axios.create(request);

	request.url = url;
	request.method = method || "GET";
	if (method === "GET") {
		request.params = data;
	} else {
		request.data = data;
	}

	Requester.interceptors.request.use((request) => {
		/**
		 * This will not work because cookie is not
		 * accessible by client.
		 * `Set-Cookie` header is not detected by response interceptor.
		 * To be able to access cookie you must set `httpOnly` to false
		 * in the backend. But, it is not recommended.
		 */
		if (request.headers?.cookie) {
			request.headers.authorization = `Bearer ${request.headers.cookie}`;
		}

		return request;
	});

	// Interceptors to handle unauthorized routes
	// ('POST', 'DELETE', 'PUT', 'PATCH')
	Requester.interceptors.response.use(
		(response) => {
			if (response.headers["set-cookie"]) {
				response.headers.cookie = response.headers["set-cookie"][0];
			}
			return response;
		},
		async (err) => {
			const statusCode = err.response?.status;
			if (statusCode === 401) {
				console.log("Please re-login");
				if (isServer === false) {
					window.location.href = "/";
					Auth.logout();
					Datastore.clear();
				}
				// call logout to clear cookie and session
				await logout();
			}

			throw err;
		},
	);

	console.log("Requester: ", request.method, request.url, request.baseURL, request.data);
	return Requester(request)
		.then((res) => (responseType === "RAW" ? res : res.data))
		.catch((err) => {
			console.log("Error occured", err);
			if (axios.isCancel(err)) return;
			if (axios.isAxiosError(err)) {
				// console.error("Request Failed: Axios Error ", err);
				throw err.response?.data;
			} else {
				console.error("Request Failed: unexpected error occured", err);
				throw err;
			}
		});
};

export { axios as default, requester };
