<!--
@component
The Rect component needs to be placed either inside a konva Layer or Group component. 

### Usage:
```tsx
<Rect config={{ x: 100, y: 100, width: 100, height: 50, fill: "blue" }} />
```
-->
<script lang="ts">

	import Konva from "konva";
	import { onMount, onDestroy, createEventDispatcher } from "svelte";
	import type { Writable } from "svelte/store";
	import { getParentContainer, type KonvaParent } from "$src/lib/utils/konva/context";
	import { registerEvents } from "$src/lib/utils/konva/events";
	import { copyExistingKeys } from "$src/lib/utils/konva/object";


	export let config: Konva.RectConfig;
	if (config) {
		config.perfectDrawEnabled = false;
		// config.listening = false;
	}
	export let handle = new Konva.Rect(config);
    export let getHandler: (rect: Konva.Rect) => void = () => null;

	let parent: Writable<null | KonvaParent> = getParentContainer();
    let dispatcher = createEventDispatcher();

    $: handle.setAttrs(config);

	onMount(() => {
		$parent!.add(handle);

		/**
         * Need to update the config object with new position
        */
		handle.on("transformend", () => {
			copyExistingKeys(config, handle.getAttrs());
			config = config;
		});

		handle.on("dragend", () => {
			copyExistingKeys(config, handle.getAttrs());
			config = config;
		});
		handle.on("dragmove", () => {
			copyExistingKeys(config, handle.getAttrs());
			config = config;
		});
		registerEvents(dispatcher, handle);

		if (typeof getHandler === "function") {
			getHandler(handle);
		}
	});

	onDestroy(() => {
		handle.destroy();
	});
</script>