<!--
@component
The Group component needs to be placed inside a konva Layer or Group component. 

The Group component automatically groups all components that are placed inside it.

### Usage:
```tsx
<Group>
	Place components that should be grouped here
</Group>
```
-->
<script lang="ts">
	import { getParentContainer, setContainerContext, Container, type KonvaParent } from "$src/lib/utils/konva/context";
	import { registerEvents, type KonvaEvents } from "$src/lib/utils/konva/events";
	import { copyExistingKeys } from "$src/lib/utils/konva/object";
	import Konva from "konva";
	import { onMount, onDestroy, createEventDispatcher } from "svelte";
	import { type Writable, writable } from "svelte/store";

	interface $$Events extends KonvaEvents {}

	export let config: Konva.GroupConfig = {};
	config.perfectDrawEnabled = false;
	export let handle: Konva.Group = new Konva.Group(config);
	export let staticConfig = false;
	export let getHandler: (handler: Konva.Group) => void = () => {
		return;
	};

	let inner = writable<null | Konva.Group>(null);

	let dispatcher = createEventDispatcher();

	let isReady = false;

	$: if (handle) {
		handle.setAttrs(config);
	}

	let parent: Writable<null | KonvaParent> = getParentContainer();

	onMount(() => {
		$parent!.add(handle);

		if (typeof getHandler === "function") getHandler(handle);

		if (!staticConfig) {
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
		}

		registerEvents(dispatcher, handle);

		inner.set(handle);
		isReady = true;
	});

	onDestroy(() => {
		if (handle) {
			handle.destroy();
		}
	});

	setContainerContext(Container.Group, inner);
</script>

{#if isReady}
	<slot />
{/if}
