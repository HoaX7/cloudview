<script lang="ts">
	import { createEventDispatcher } from "svelte";
	import css from "./input.module.css";
	import clsx from "clsx";
	export let classname = "";
	export let value = "";
	const dispatch = createEventDispatcher();
</script>

<input
	type="text"
	bind:value
	class={clsx(
		classname,
		css["input"],
		"focus:outline-gray-300",
		"rounded w-full border-2 disabled:opacity-75 disabled:cursor-not-allow"
	)}
	{...$$restProps}
    on:input={() => dispatch("input", value)}
/>
