<!--
    @component
    Grid layer to use as stage background for canvas.
    The Component is extremely lagging if the width and height of the gridbox is small

    DO NOT USE THIS COMPONENT
    
-->
<script lang="ts">
	import { getParentStage } from "$src/lib/utils/konva/context";
	import type Konva from "konva";
	import { onMount } from "svelte";
	import type { Writable } from "svelte/store";
	import Layer from "./Layer.svelte";
	import Rect from "./Rect.svelte";

	const width = 100;
	const height = 100;

	let parent: Writable<null | Konva.Stage> = getParentStage();
    let components: { x: number; y: number; }[] = [];

	onMount(() => {
		const stagePos = $parent?.getPosition();
		if (stagePos) {
			const startX = Math.floor((-stagePos.x - window.innerWidth) / width) * width;
			const endX = Math.floor((-stagePos.x + window.innerWidth * 2) / width) * width;

			const startY = Math.floor((-stagePos.y - window.innerHeight) / height) * height;
			const endY = Math.floor((-stagePos.y + window.innerHeight * 2) / height) * height;
			for (var x = startX; x < endX; x += width) {
				for (var y = startY; y < endY; y += height) {
					components.push({
						x,
						y 
					});
				}
			}
		}
	});
</script>

<Layer>
    {#each components as grid, index (index)}
        <Rect config={{
        	x: grid.x,
        	y: grid.y,
        	width,
        	height,
        	draggable: false,
        	fill: "white",
        	stroke: "#ccc"
        }} />
    {/each}
</Layer>
