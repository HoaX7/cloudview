<!--
@component
The Arrow component needs to be placed either inside a Layer or Group component. 

### Usage:
```tsx
<Arrow config={{ x: 100, y: 100, points: [0, 0, 40, 40], 
    pointerLength: 20, pointerWidth: 20, fill: "blue", 
    stroke: "blue", strokeWidth: 6 }} />
```
-->
<script lang="ts">
  import {
  	getParentContainer,
  	type KonvaParent,
  } from "$src/lib/utils/konva/context";
  import {
  	registerEvents,
  	type KonvaEvents,
  } from "$src/lib/utils/konva/events";
  import { copyExistingKeys } from "$src/lib/utils/konva/object";

  import Konva from "konva";
  import { onMount, onDestroy, createEventDispatcher } from "svelte";
  import type { Writable } from "svelte/store";

  interface $$Events extends KonvaEvents {}

  export let config: Konva.ArrowConfig;
  if (config) {
  	config.perfectDrawEnabled = false;
  	config.listening = false;
  }
  export let handle = new Konva.Arrow(config);
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
