
interface Props {
    headers?: Record<string, unknown>,
    url: string;
    method: "GET" | "POST" | "PATCH" | "DELETE",
    data: any;
}
export const requester = async ({
	url,
	method,
	data = {},
	headers
}: Props) => {
	try {
		// const base = "http://127.0.0.1:6002";
		const base = "https://api_service_worker.getcloudfriendly.com";
		const options: RequestInit = {
			method,
			headers: {
				"accept": "application/json",
				"content-type": "application/json",
				...headers
			},
			referrerPolicy: "no-referrer"
		};
		if (method !== "GET") {
			options.body = JSON.stringify(data);
		}
		return fetch(`${base}/${url}`, options).then((res) => res.json())
			.catch(err => {
				throw err;
			});
	} catch (err) {
		console.error("helpers.requester: Failed", err);
		throw err;
	}
};