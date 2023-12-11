<!--
    @component
    Grid layer to use as stage background for canvas.
    The Component is extremely lagging if the width and height of the gridbox is small

    DO NOT USE THIS COMPONENT
    
-->
<script lang="ts">
  import { getParentStage } from "$src/lib/utils/konva/context";
  import Konva from "konva";
  import Layer from "./Layer.svelte";
  import type { LineConfig } from "konva/lib/shapes/Line";
  import { onMount } from "svelte";

  const gridOutsizePercentage = .2;
  const fieldSize = 50;

  const staticLineConfig = {
  	draggable: false,
  	listening: false,
  	stroke: "rgba(190,190,190,0.1)",
  	strokeWidth: 2,
  	lineCap: "round",
  	perfectDrawEnabled: false,
  } as LineConfig;
  let startX: number, endX: number, startY: number, endY: number;

  let layer: Konva.Layer;

  const line = new Konva.Line(staticLineConfig);
  line.cache();

  const stage = getParentStage();

  onMount(() => {
  	updateGrid();
  	$stage?.on("dragend", updateGrid);
  });

  // calculate position of the grid lines and draw them
  function updateLines() {
  	// dont draw grid if zoomed out to wide
  	if (Math.pow(($stage?.scaleX() || 1), -1) > 10) return;

  	// calculate coordinates of the top left corner of the stage
  	let transform = $stage?.getAbsoluteTransform().copy().invert();
  	if (!transform) return;
  	let topLeftCorner = transform.point({
  		x: 0,
  		y: 0,
  	});

  	// calculate the width and height of the grid
  	calculatePoints(topLeftCorner);

  	// draw horizontal and vertical lines
  	createLines(startX, endX, true);
  	createLines(startY, endY, false);
  }
  export function updateGrid() {
  	layer.destroyChildren();
  	updateLines();
  	layer.batchDraw();
  }

  // calculate and save the current dimensions of the bord
  function calculatePoints(corner: { x: number; y: number }) {
  	if (!$stage) return;
  	startX =
      Math.floor(
      	(corner.x -
          $stage.width() *
            gridOutsizePercentage *
            Math.pow($stage.scaleX(), -1)) /
          fieldSize
      ) * fieldSize;
  	endX =
      Math.floor(
      	(corner.x +
          $stage.width() *
            (1 + gridOutsizePercentage) *
            Math.pow($stage.scaleX(), -1)) /
          fieldSize
      ) * fieldSize;

  	startY =
      Math.floor(
      	(corner.y -
          $stage.height() *
            gridOutsizePercentage *
            Math.pow($stage.scaleY(), -1)) /
          fieldSize
      ) * fieldSize;
  	endY =
      Math.floor(
      	(corner.y +
          $stage.height() *
            (1 + gridOutsizePercentage) *
            Math.pow($stage.scaleY(), -1)) /
          fieldSize
      ) * fieldSize;
  }

  // draw lines
  function createLines(start: number, end: number, horizontalLines: boolean) {
  	let clone;

  	for (let n = start; n < end; n += fieldSize) {
  		clone = line.clone({ points: horizontalLines ? [ n, startY, n, endY ] : [ startX, n, endX, n ], });

  		// some config options to improve performance
  		clone.perfectDrawEnabled(false);
  		clone.shadowForStrokeEnabled(false);
  		clone.transformsEnabled("position");

  		layer.add(clone);
  	}
  }
</script>

<Layer config={{ listening: false }} getHandler={(handle) => {
	layer = handle;
}}>
</Layer>
