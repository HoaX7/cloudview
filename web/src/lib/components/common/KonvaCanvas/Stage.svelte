<!--
@component
The Stage component is the entry point and parent for all other konva components.
Under the hood it creates a div element where the html canvas is attached to.  

### Usage:
```tsx
<Stage config={{ width: 1000, height: 1000 }}>
	Place your Layers here
</Stage>
```
-->
<script lang="ts">
	import Konva from "konva";
	import { onMount, onDestroy, createEventDispatcher } from "svelte";
	import { writable } from "svelte/store";
	import { Container, setContainerContext } from "$src/lib/utils/konva/context";
	import { registerEvents } from "$src/lib/utils/konva/events";
  import { NAVBAR_HEIGHT } from "$src/helpers/constants";

	export let config: Konva.ContainerConfig = {
		width: 0,
		height: 0,
		draggable: true,
		x: 0,
		y: 0,
	};
	let stage: HTMLDivElement;
	export let handle: Konva.Stage | null = null;
    export let getHandler: (stage: Konva.Stage) => void = () => null;
	export let customDimensions = false;

	let inner = writable<null | Konva.Stage>(null);

	$: if (handle) {
		handle.setAttrs(config);
	}
	let isReady = false;

    let dispatcher = createEventDispatcher();

	onMount(() => {
		if (!customDimensions) {
			config.width = window.innerWidth;
			config.height = window.innerHeight - NAVBAR_HEIGHT;
		}
		handle = new Konva.Stage({
			container: stage,
			...config
		});

		if (typeof getHandler === "function") getHandler(handle);

		const stageComponent = handle.getContent();
		handle.on("dragstart", () => {
			stageComponent.classList.add("gbr");
		});
		handle.on("dragend", () => {
			stageComponent.classList.remove("gbr");
		});
		// stageComponent.style.backgroundImage = "radial-gradient(black 1px, transparent 0)";
		// stageComponent.style.backgroundSize = "40px 40px";
		// stageComponent.style.backgroundPosition = "-19px -19px";

		inner.set(handle);
		isReady = true;

		registerEvents(dispatcher, handle);
	});

	onDestroy(() => {
		if (handle) {
			handle.destroy();
		}
	});

	setContainerContext(Container.Stage, inner);
</script>

<div bind:this={stage} {...$$restProps}>
	{#if isReady}
		<slot />
	{/if}
</div>
