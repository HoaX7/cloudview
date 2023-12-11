<script lang="ts">
  import type Konva from "konva";
  import Stage from "../Stage.svelte";
  import Layer from "../Layer.svelte";
  import { createEventDispatcher, onMount } from "svelte";
  import type { StageConfig } from "konva/lib/Stage";
  import Rect from "../Rect.svelte";
  import type { RectConfig } from "konva/lib/shapes/Rect";
  import Konvastore from "$src/store/konva";
  import ServiceGroupWithLabel from "$src/lib/components/services/aws/ServiceGroupWithLabel.svelte";
  import { delay } from "$src/helpers";

  export let scale = 1;
  const konvastore = Konvastore.getStore();
  let stage: Konva.Stage;

  const minimapW = window.innerWidth / 7;
  const minimapH = window.innerHeight / 7;
  let stageConfig = {
  	draggable: false,
  	listening: true,
  	width: minimapW,
  	height: minimapH,
  	x: 0,
  	y: 0
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
  	const scaleX = (minimapW / mainCanvasProportions.width) - 0.015; // padding
  	const scaleY = (minimapH / mainCanvasProportions.height) - 0.015;
  	const viewScale = Math.min(scaleX, scaleY);
  	stageConfig.scaleX = viewScale;
  	stageConfig.scaleY = viewScale;
  	// viewBoxConfig.x = mainCanvasProportions.x / scale;
  	// viewBoxConfig.y = mainCanvasProportions.y / scale;
  };

  onMount(async () => {
  	await delay(200);
  	drawMap(stage?.getClientRect());
  });


  const dispatch = createEventDispatcher();

</script>

<Stage
  bind:config={stageConfig}
  class="bg-gray-100 shadow absolute top-5 right-5 rounded"
  customDimensions
  getHandler={(handle) => {
  	stage = handle;
  }}
>
    <Layer>
      {#each Array($konvastore.rowCount.external) as _, idx (idx)}
      <ServiceGroupWithLabel idx={idx} label={{
      	text: "",
      	fill: "" 
      }}  borderColor="#2d9cdb" externalService={true} />
      {/each}
      {#each Array($konvastore.rowCount.internal) as _, idx (idx)}
      <ServiceGroupWithLabel idx={idx} label={{
      	text: "",
      	fill: "" 
      }}  borderColor="#2d9cdb" externalService={false} />
      {/each}
    </Layer>
    <Layer>
      <Rect bind:config={viewBoxConfig} on:dragmove={() => {
      	dispatch("dragmove", viewBoxConfig);
      }} />
    </Layer>
</Stage>
