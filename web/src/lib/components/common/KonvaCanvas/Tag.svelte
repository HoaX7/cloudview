<!--
@component
The Tag component needs to be placed either inside a svelte-konva Layer or Group component. 

### Usage:
```tsx
<Tag config={{ x: 10, y: 20, fill: "black", 
pointerDirection: "down", pointerWidth: 10, pointerHeight: 10, lineJoin: "round" }} />
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

	export let config: Konva.TagConfig;
	export let handle = new Konva.Tag(config);
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