import { defineConfig } from "astro/config";
import tailwind from "@astrojs/tailwind";
import sitemap from "@astrojs/sitemap";
import cloudflare from "@astrojs/cloudflare";
import compressor from "astro-compressor";
import purgecss from "astro-purgecss";

// https://astro.build/config
export default defineConfig({
	integrations: [ tailwind(), sitemap({
		priority: .8,
		changefreq: "weekly",
		lastmod: new Date("2023-02-01")
	}), compressor(), purgecss() ],
	site: "https://getcloudfriendly.com",
	trailingSlash: "ignore",
	// required for cloudflare to remove trailing slash ('/')
	build: { format: "file" },
	output: "hybrid",
	adapter: cloudflare(),
	image: {
		domains: [ "astro.build" ],
		service: { entrypoint: "astro/assets" }
	}
});