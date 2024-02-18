<script lang="ts">
	import clsx from "clsx";
	import Button from "../common/Button/Button.svelte";
	import Typography from "../common/Typography/Typography.svelte";
	import { OAUTH_URLS } from "$src/helpers/constants";
	import FullPageLoader from "../common/Loaders/FullPageLoader.svelte";
  import Brandlogo from "./Brandlogo.svelte";

	let isLoading = false;
</script>

<main class="container mx-auto h-full">
	{#if isLoading}
		<FullPageLoader />
	{/if}
	<div class="flex flex-col items-center justify-center h-full">
		<div class="w-80 flex flex-col justify-center text-center items-center">
			<Brandlogo font={32} />
			<Typography variant="p" weight="medium" font={16} classname="mt-3">
				Visualize and Monitor your resources deployed on AWS, GCP, Azure and more in one place
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
				<img alt="cloudfriendly, login with google" width={24} 
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
				<img alt="cloudfriendly, login with github" width={24} 
                src="/assets/images/github-icon.svg" class="mr-2" /> Login with
				Github
			</Button>
		</div>
	</div>
</main>
