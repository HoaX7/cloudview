<!--
@component
The Layer component needs to be placed inside a konva Stage component. 

### Usage:
```tsx
<Stage config={{ width: 1000, height: 1000 }}>
	<Layer>
		Place your components here
	</Layer>

	You also can add multiple Layers to a Stage
</Stage>
```
-->
<script lang="ts">
	import Konva from "konva";
	import { onMount, onDestroy, createEventDispatcher } from "svelte";
	import { type Writable, writable } from "svelte/store";
	import { Container, getParentStage, setContainerContext } from "$src/lib/utils/konva/context";
	import { registerEvents } from "$src/lib/utils/konva/events";

	export let config: Konva.LayerConfig = {};
	config.perfectDrawEnabled = false;
	export let handle: Konva.Layer = new Konva.Layer(config);
    export let getHandler: (layer: Konva.Layer) => void = () => null;

	let layer: HTMLDivElement;

	let inner = writable<null | Konva.Layer>(null);
	let isReady = false;

	let parent: Writable<null | Konva.Stage> = getParentStage();
    let dispatcher = createEventDispatcher();

	onMount(() => {
		$parent!.add(handle);
		if (typeof getHandler === "function") getHandler(handle);
		inner.set(handle);
		isReady = true;
        
		registerEvents(dispatcher, handle);
	});

	onDestroy(() => {
		if (handle) {
			handle.destroy();
		}
	});

	setContainerContext(Container.Layer, inner);
</script>

<div bind:this={layer} {...$$restProps}>
	{#if isReady}
	<slot />
{/if}
</div>