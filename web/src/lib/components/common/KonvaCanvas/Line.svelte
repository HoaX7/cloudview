<!--
@component
The Line component needs to be placed either inside a Layer or Group component. 

### Usage:
```tsx
<Line config={{ points: [0, 0, 60, 30, 300, 90, 30, 100], stroke: "blue", strokeWidth: 10 }} />
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

	export let config: Konva.LineConfig;
	if (config) {
		config.perfectDrawEnabled = false;
	}
	export let handle = new Konva.Line(config);
	handle.cache({
		x: config.x,
		y: config.y,
		width: config.width,
		height: config.height
	});
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