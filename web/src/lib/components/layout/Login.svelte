<script lang="ts">
	import clsx from "clsx";
	import Button from "../common/Button/Button.svelte";
	import Typography from "../common/Typography/Typography.svelte";
	import { OAUTH_URLS } from "$src/helpers/constants";
	import FullPageLoader from "../common/Loaders/FullPageLoader.svelte";

	let isLoading = false;
</script>

<main class="container mx-auto h-full">
	{#if isLoading}
		<FullPageLoader />
	{/if}
	<div class="flex flex-col items-center justify-center h-full">
		<div class="w-80 flex flex-col justify-center text-center items-center">
			<Typography variant="h1" weight="medium" font={32}>Cloudview</Typography>
			<Typography variant="p" weight="medium" font={16} classname="text-gray-500 mt-3">
				Manage and Monitor your micro services deployed on AWS, GCP, Azure and more in one place
			</Typography>
			<Button
				classname={clsx(
					"hover:bg-opacity-50 hover:bg-gray-100",
					"border border-gray-100 shadow flex items-center mt-3",
					"disabled:opacity-50"
				)}
				type="button"
				disabled={true}
			>
				<img alt="cloudview, login with google" width={24} 
                src="/assets/images/google-icon.svg" class="mr-2" /> Login with
				Google
			</Button>
			<Button
				classname={clsx(
					"hover:bg-opacity-50 hover:bg-gray-100 mt-3",
					"border border-gray-100 shadow flex items-center mt-3"
				)}
				type="button"
                on:click={(e) => {
                	isLoading = true;
                	window.location.assign(OAUTH_URLS.GITHUB);
                }}
			>
				<img alt="cloudview, login with github" width={24} 
                src="/assets/images/github-icon.svg" class="mr-2" /> Login with
				Github
			</Button>
		</div>
	</div>
</main>
