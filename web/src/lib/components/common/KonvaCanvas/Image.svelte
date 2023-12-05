<!--
@component
The Image component needs to be placed either inside a Layer or Group component. 

### Usage:
```tsx
<Image position={{ x: 100, y: 100, width: 100, height: 100 }} config={{ image: imageObj }} />
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

	export let config: Pick<Konva.ImageConfig, "image">;
    export let position: Omit<Konva.ImageConfig, "image">;
	export let handle = new Konva.Image({
		image: config.image,
		...position,
		perfectDrawEnabled: false,
		listening: false,
	});
	export let staticConfig = false;

	let parent: Writable<null | KonvaParent> = getParentContainer();
	let dispatcher = createEventDispatcher();

	$: handle.setAttrs(config);

	onMount(() => {
		$parent!.add(handle);

		if (!staticConfig) {
			handle.on("transformend", () => {
				copyExistingKeys(position, handle.getAttrs());
				position = position;
			});

			handle.on("dragend", () => {
				copyExistingKeys(position, handle.getAttrs());
				position = position;
			});
		}

		registerEvents(dispatcher, handle);
	});

	onDestroy(() => {
		handle.destroy();
	});
</script>