<!--
@component
The Transformer component needs to be placed inside a konva Layer or Group component. 

In order to add shapes to the transformer you need to access the underlying 
Konva Transformer by binding the `handle` prop. 
Then use the `nodes()` function to add any shapes to the Transformer.

### Usage:
```tsx
<script>
	let transformer;

	transformer.nodes([someShape, otherShape]);
</script>

<Transformer bind:handle={transformer} />
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

	export let config: Konva.TransformerConfig = {};
	export let handle = new Konva.Transformer(config);
	export let staticConfig = false;

	let parent: Writable<null | KonvaParent> = getParentContainer();
	let dispatcher = createEventDispatcher();

	$: if (handle) {
		handle.setAttrs(config);
	}

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
		if (handle) {
			handle.destroy();
		}
	});
</script>