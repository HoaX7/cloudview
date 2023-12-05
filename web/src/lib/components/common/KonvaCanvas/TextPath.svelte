<!--
@component
The TextPath component needs to be placed either inside a Layer or Group component. 

### Usage:
```tsx
<TextPath config={{ x: 100, y: 100, fill: "#333", text: "some text", fontSize: 25,
 data: "M10 10 C0 0 10 150 100 100 S300 150 5.0.300" }} />
```
-->
<script lang="ts">
	import { getParentContainer, type KonvaParent } from "$src/lib/utils/konva/context";
	import { registerEvents, type KonvaEvents } from "$src/lib/utils/konva/events";
	import { copyExistingKeys } from "$src/lib/utils/konva/object";


	import Konva from "konva";
	import { onMount, onDestroy, createEventDispatcher } from "svelte";
	import type { Writable } from "svelte/store";

	interface $$Events extends KonvaEvents {}

	export let config: Konva.TextPathConfig;
	export let handle = new Konva.TextPath(config);
	export let staticConfig = false;

	let parent: Writable<null | KonvaParent> = getParentContainer();
	let dispatcher = createEventDispatcher();

	$: handle.setAttrs(config);

	onMount(() => {
		$parent!.add(handle);

		if (!staticConfig) {
			handle.on("transformend", () => {
				copyExistingKeys(config, handle.getAttrs());
				config = config;
			});

			handle.on("dragend", () => {
				copyExistingKeys(config, handle.getAttrs());
				config = config;
			});
		}

		registerEvents(dispatcher, handle);
	});

	onDestroy(() => {
		handle.destroy();
	});
</script>
