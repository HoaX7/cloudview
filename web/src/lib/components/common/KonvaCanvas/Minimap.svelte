<script lang="ts">
  import type Konva from "konva";
  import Stage from "./Stage.svelte";
  import Layer from "./Layer.svelte";
  import { createEventDispatcher, onDestroy } from "svelte";
  import type { StageConfig } from "konva/lib/Stage";
  import Rect from "./Rect.svelte";
  import type { RectConfig } from "konva/lib/shapes/Rect";

  export let layer: Konva.Layer;
  export let scale = 1;
  let previewLayer: null | Konva.Layer = null;

  const minimapW = window.innerWidth / 5;
  const minimapH = window.innerHeight / 6;
  let stageConfig = {
  	draggable: false,
  	listening: true,
  	width: minimapW,
  	height: minimapH,
  	scaleX: 1 / 6,
  	scaleY: 1 / 6,
  } as StageConfig;
  let viewBoxConfig = {
  	draggable: true,
  	listening: true,
  	width: window.innerWidth / scale,
  	height: window.innerHeight / scale,
  	x: 0,
  	y: 0,
  	fill: "rgba(86, 204, 242, 0.1)",
  	stroke: "#2d9cdb",
  	strokeWidth: 20,
  } as RectConfig;

  export const drawMap = (mainCanvasProportions: any) => {
  	const scaleX = minimapW / mainCanvasProportions.width - 0.015; // padding
  	const scaleY = minimapH / mainCanvasProportions.height - 0.015;
  	const viewScale = Math.min(scaleX, scaleY);
  	stageConfig.scaleX = viewScale;
  	stageConfig.scaleY = viewScale;
  	previewLayer = layer.clone({ listening: false });
  	viewBoxConfig.x = mainCanvasProportions.x / scale;
  	viewBoxConfig.y = mainCanvasProportions.y / scale;
  };

  onDestroy(() => {
  	previewLayer?.destroy();
  });

  const dispatch = createEventDispatcher();
</script>

<Stage
  bind:config={stageConfig}
  class="bg-gray-100 shadow absolute top-5 right-5 rounded"
  customDimensions
  on:click={() => {
  	if (previewLayer) {
  		const proportions = previewLayer.getRelativePointerPosition();
  		if (!proportions) return;
  		viewBoxConfig.x = proportions.x - ((viewBoxConfig?.width || 0) / 2);
  		viewBoxConfig.y = proportions.y - ((viewBoxConfig?.height || 0) / 2);
  		dispatch("dragmove", {
  			x: viewBoxConfig.x,
  			y: viewBoxConfig.y
  		});
  	}
  }}
>
  {#if previewLayer}
    <Layer handle={previewLayer} />
    <Layer>
      <Rect bind:config={viewBoxConfig} on:dragmove={() => {
      	dispatch("dragmove", viewBoxConfig);
      }} />
    </Layer>
  {/if}
</Stage>
