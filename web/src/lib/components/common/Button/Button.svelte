<script lang="ts">
	import clsx from "clsx";
	import { createEventDispatcher } from "svelte";

	export let id: string = "";
	export let type: "button" | "submit" | "reset" | null | undefined = null;
	export let classname: string = "";

	const dispatch = createEventDispatcher();

</script>

<button
	{id}
	{type}
	class={clsx("px-10 disabled:opacity-75 disabled:pointer-events-none",
		"p-3 rounded-full font-semibold", classname)}
	{...$$restProps}
	on:click={(e) => dispatch("click", e)}
>
	<slot />
</button>
