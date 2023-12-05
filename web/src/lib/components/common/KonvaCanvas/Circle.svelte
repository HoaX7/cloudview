<!--
@component
The Circle component needs to be placed either inside a konva Layer or Group component. 

### Usage:
```tsx
<Circle config={{ x: 100, y: 100, radius: 50, fill: "blue" }} />
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

	export let config: Konva.CircleConfig;
	if (config) {
		config.perfectDrawEnabled = false;
		config.listening = false;
	}
	export let handle = new Konva.Circle(config);
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
