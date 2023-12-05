<!--
@component
The Path component needs to be placed either inside a konva Layer or Group component. 

### Usage:
```tsx
<Path config={{ x: 100, y: 100, width: 100, height: 100, fill: "blue", 
data: "M213.1,6.7c-32.4-14.4-73.7,0-88.1,30.6C110.6,4.9"
",67.5-9.5,36.9,6.7C2.8,22.9-13.4,62.4,13.5,110.9C33.3,145.1"
",67.5,170.3,125,217c59.3-46.7,93.5-71.9,111.5-106.1C263.4,64.2,247.2,22.9,213.1,6.7z" }} />
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

  export let config: Konva.PathConfig;
  if (config) {
  	config.perfectDrawEnabled = false;
  }
  export let handle = new Konva.Path(config);
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
