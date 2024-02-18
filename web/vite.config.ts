import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vite";

export default defineConfig({
	plugins: [ sveltekit() ],
	server: { port: 5000, },
	// use this host to allow connections from cloudrun
	// `npm run preview` is used to run in prod
	preview: {
		port: 5000,
		strictPort: true,
		host: "0.0.0.0"
	}
});
