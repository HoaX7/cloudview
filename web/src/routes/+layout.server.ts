import { COOKIE_NAME } from "$src/helpers/constants";
import type { LayoutServerLoad } from "./$types";

/**
 * By adding `ssr = false` we have converted our entire app
 * into SPA. In most cases it is not recommended.
 * This must be added in a layout.ts file.
 * `+layout.server.ts` file always only runs on the server, while
 * `+layout.ts` file runs on both client & server.
 * https://kit.svelte.dev/docs/page-options#ssr
 */
// export const ssr = false;

export const load = (({ cookies, url, params }) => {
	const cookie = cookies.get(COOKIE_NAME);
	const code = url.searchParams.get("code");
	return {
		code,
		params,
		cookie,
	};
}) satisfies LayoutServerLoad;