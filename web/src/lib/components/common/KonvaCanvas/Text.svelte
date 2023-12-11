<!--
@component
The Text component needs to be placed either inside a Layer or Group component. 

### Usage:
```tsx
<Text config={{ x: 100, y: 100, text: "some text", fontSize: 25 }} />
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

	export let config: Konva.TextConfig;
	if (config) {
		config.perfectDrawEnabled = false;
		config.listening = false;
	}
	export let handle = new Konva.Text(config);
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
		dispatcher("ready", handle);
	});

	onDestroy(() => {
		handle.destroy();
	});
</script>