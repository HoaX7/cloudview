<script lang="ts">
	import clsx from "clsx";
	import css from "./Loader.module.css";

	export let title = "";
	export let description = "";
</script>

<div
	class={clsx(
		"fixed top-0 left-0 right-0 bottom-0 w-full h-screen z-50 overflow-hidden",
		"bg-gray-700 opacity-75 flex flex-col items-center justify-center"
	)}
>
	<div
		class={clsx("ease-linear rounded-full border-4",
			"border-t-4 h-12 w-12 mb-4", css["fullpage-loader"])}
	/>
	<h2 class="text-center text-white text-xl font-semibold">
		{title}
	</h2>
	{#if description}
		<p class="w-1/3 text-center text-white">
			{description}
		</p>
	{/if}
</div>
