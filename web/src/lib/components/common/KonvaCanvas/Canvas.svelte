<script lang="ts">
  import Stage from "./Stage.svelte";
  import Layer from "./Layer.svelte";
  import clsx from "clsx";
  import { NAVBAR_HEIGHT } from "$src/helpers/constants";
  import Konva from "konva";
  import Datastore from "$src/store/data";
  import {
  	drawSVGPath,
  	getArrowHeadPoints,
  	getConnectorPointsByPosition,
  } from "$src/helpers/konva/index";
  import { afterUpdate, createEventDispatcher, onDestroy, onMount } from "svelte";
  import type {
  	ConnectableNodeProps,
  	HighLightProps,
  	LegendProps,
  	TargetFromNodeProps,
  } from "$src/customTypes/Konva";
  import Icon from "../Image/index.svelte";
  import Legend from "../../services/views/legend.svelte";
  import type { PathConfig } from "konva/lib/shapes/Path";
  import type { CircleConfig } from "konva/lib/shapes/Circle";
  import { COLOR_SCHEME } from "$src/colorConfig";
  import type { ArrowConfig } from "konva/lib/shapes/Arrow";
  import { clone, delay, isMobile } from "$src/helpers";
  import Minimap from "./minimap/Minimap.svelte";
  import GridLayer from "./GridLayer.svelte";
  import SettingStore from "$src/store/settings";
  import Settings from "./Settings.svelte";
  // import { onMount } from "svelte";

  // To draw connecting arrows between nodes
  /**
   * Make sure that datastore is only being used by canvas.
   * Otherwise it will cause to-many re-renders.
   *
   * If there's other data to be stored in store create a different store
   * for it.
   *
   * `auth` store is used by profile only.
   */
  const datastore = Datastore.getDatastore();
  const settingStore = SettingStore.getStore();

  export let legend: LegendProps[] = [];
  export let highlights: HighLightProps;
  $: {
  	if ((highlights.lines || []).length > 0) highlightLines(highlights.lines);
  }

  let state = {
  	centerX: 0,
  	centerY: 0,
  	grow: false,
  };
  let stage: Konva.Stage;
  let layer: Konva.Layer;

  /**
   * Caching static shapes for better performance.
   */
  let svgPath = new Konva.Path({
  	listening: false,
  	perfectDrawEnabled: false,
  	stroke: COLOR_SCHEME.CONNECTOR,
  	strokeWidth: 3,
  	draggable: false,
  	data: "M0 0 L10 10",
  	zIndex: 0,
  });
  let circle = new Konva.Circle({
  	listening: false,
  	draggable: false,
  	radius: 6,
  	stroke: COLOR_SCHEME.CONNECTOR,
  	fill: "white",
  	zIndex: 999,
  	strokeWidth: 2,
  	perfectDrawEnabled: false,
  });
  let arrow = new Konva.Arrow({
  	draggable: false,
  	listening: false,
  	perfectDrawEnabled: false,
  	points: [],
  	stroke: COLOR_SCHEME.CONNECTOR,
  	strokeWidth: 2,
  	fill: COLOR_SCHEME.CONNECTOR,
  });
  arrow.cache();
  circle.cache();
  svgPath.cache();

  /**
   * Need to calculate centerX and centerY
   * to focus the selected group at the center of the canvas.
   * @param newX
   * @param newY
   * @param width
   * @param height
   */
  export const handleRepositionStage = (
  	newX: number,
  	newY: number,
  	width: number,
  	height: number
  ) => {
  	const centerX = (window.innerWidth - width) / 2;
  	// The height of Nav bar is 64px
  	const centerY = (window.innerHeight - NAVBAR_HEIGHT - height) / 2;
  	newX = -(newX - centerX);
  	newY = -(newY - centerY);
  	if (!state.grow) {
  		stage.to({
  			x: newX,
  			y: newY,
  			duration: 0.5,
  		});
  	} else {
  		stage.to({
  			x: 0,
  			y: 0,
  			duration: 1,
  		});
  	}
  };

  type L = {
    circle: CircleConfig;
    line: PathConfig;
    arrow: ArrowConfig;
  };
  let linesToDraw: L[] = [];
  let linesArray: L[] = [];
  let minimapRef: any = null;
  let gridLayerRef: any = null;

  const isDesktop = !isMobile();

  const clearLines = (lines: L[]) => {
  	const idsToDelete = lines.map((ln) => ln.line.id);
  	const children = layer
  		?.getChildren()
  		.filter((child) => idsToDelete.includes(child.attrs.id));
  	children?.forEach((child) => child.destroy());
  };

  const updateConnectorShapes = (
  	lines: L[],
  	highlight = false,
  	animate = true
  ) => {
  	clearLines(lines);
  	lines.forEach((item) => {
  		const pathClone = svgPath.clone(item.line);
  		const circleClone = circle.clone(item.circle);
  		const arrowClone = arrow.clone(item.arrow);
  		if (animate && $settingStore.animate) {
  			arrowClone.visible(false);
  			const pathLen = pathClone.getLength();
  			pathClone.dash([ pathLen ]);
  			pathClone.dashOffset(pathLen);
  			const anim = new Konva.Animation(function (frame) {
  				const dashLen = pathLen - (frame?.time || 0) / 5;
  				pathClone.dashOffset(dashLen);
  				if (dashLen < 0) {
  					anim.stop();
  					arrowClone.visible(true);
  				}
  			}, layer);
  			anim.start();
  		}
  		layer?.add(pathClone);
  		if (highlight) {
  			pathClone.moveToTop();
  		} else {
  			pathClone.moveToBottom();
  		}
  		layer?.add(circleClone);
  		layer?.add(arrowClone);
  	});
  	layer?.batchDraw();
  };

  const drawConnectorLines = (
  	target: ConnectableNodeProps,
  	node: TargetFromNodeProps
  ) => {
  	const { from, to } = getConnectorPointsByPosition(
  		{
  			x: node.x,
  			y: node.y,
  		},
  		{
  			x: target.x,
  			y: target.y,
  		}
  	);
  	const arrowHeadPoints = getArrowHeadPoints(to);
  	const path = drawSVGPath(from, to);
  	const id = `line-${node.from}-${node.to}`;
  	const idx = linesArray.findIndex((ln) => {
  		return ln.line.id === id;
  	});
  	if (idx >= 0) {
  		// linesArray[idx].points = points;
  		linesArray[idx].line.data = path;
  		linesArray[idx].circle.x = from.x;
  		linesArray[idx].circle.y = from.y;
  		linesArray[idx].arrow.points = arrowHeadPoints;
  	} else {
  		linesArray.push({
  			circle: {
  				id,
  				zIndex: 999,
  				x: from.x,
  				y: from.y,
  				...(node.lineStyle ? { stroke: node.lineStyle.stroke } : {}),
  			},
  			line: {
  				data: path,
  				id,
  				...(node.lineStyle || {}),
  			},
  			arrow: {
  				id,
  				points: arrowHeadPoints,
  				...(node.lineStyle
  					? {
  						stroke: node.lineStyle.stroke,
  						fill: node.lineStyle.stroke,
  					}
  					: {}),
  			},
  		});
  	}
  };

  // Update arrow point position
  // function takes id (of the node which has moved)
  // lineId (which comes from target.from)
  export const updateConnector = (
  	targetFromNodes?: TargetFromNodeProps[],
  	connectableNodes?: ConnectableNodeProps[],
  	dragNodeId?: string
  ) => {
  	let targetFromArray = targetFromNodes || [];
  	const connectableNodeArray = connectableNodes || [];

  	const findAndSetTargetNodeArray = (id: string, useFilter = false) => {
  		if (useFilter) {
  			const arr = targetFromArray.filter(
  				(tg) => tg.id === id || tg.to.includes(id)
  			);
  			if (arr.length > 0) {
  				targetFromArray = arr;
  				return true;
  			}
  		}
  		const tgNodes = targetFromArray.filter((tg) => tg.id === id);
  		if (tgNodes.length > 0) {
  			targetFromArray = tgNodes;
  			return true;
  		}
  		return false;
  	};
  	if (dragNodeId) {
  		const targetArraySet = findAndSetTargetNodeArray(dragNodeId);
  		if (!targetArraySet) {
  			const node = connectableNodeArray.find((nd) => nd.id === dragNodeId);
  			if (node) {
  				findAndSetTargetNodeArray(dragNodeId, true);
  			}
  		}
  	}
  	/**
     * Optimization - Only pick the source and target node
     * when the elelents are dragged.
     * There's no need to redraw other static line points.
     */
  	targetFromArray.forEach((node) => {
  		/**
       * We are checking for includes because when connecting
       * `cloudfront` it only has the domain name mapped (id of the instance)
       * there is no direct way to get the exact id of the instance. Therefore,
       * we msut infer the associations.
       */
  		const target = connectableNodeArray.find(
  			(nd) => nd.id === node.to || node.to.includes(nd.id)
  		);
  		if (!target) return;
  		// const obstacles = connectableNodeArray.filter(
  		// 	(nd) => nd.id !== target?.id && nd.id !== node.id
  		// );
  		drawConnectorLines(target, node);
  		return;
  	});

  	linesToDraw = linesArray;
  	updateConnectorShapes(linesArray);
  	console.log("lines to draw:", linesArray.length);
  };

  // Update lines
  const listener = datastore.subscribe((newVal) => {
  	if (newVal.fetchData) {
  		linesArray = [];
  		clearLines(linesToDraw);
  	}
  	updateConnector(
  		newVal.konvaTargetFromNodes,
  		newVal.konvaConnectableNodes,
  		newVal.dragNodeId || undefined
  	);
  });

  onDestroy(() => {
  	listener();
  });

  const dispatch = createEventDispatcher();

  const highlightLines = (lineIds: string[]) => {
  	if (linesToDraw.length > 0 && lineIds.length > 0) {
  		const res = clone(linesToDraw);
  		res.map((ln) => {
  			if (lineIds.includes(ln.line?.id || "")) {
  				ln.line.opacity = 1;
  				ln.circle.opacity = 1;
  				ln.arrow.opacity = 1;
  				ln.circle.zIndex = 999;
  				ln.line.zIndex = 20;
  				ln.arrow.zIndex = 999;
  				return ln;
  			}
  			ln.line.opacity = 0.1;
  			ln.circle.opacity = 0.1;
  			ln.arrow.opacity = 0.1;
  			ln.line.zIndex = 0;
  			ln.arrow.zIndex = 0;
  			ln.circle.zIndex = 10;
  			return ln;
  		});
  		linesToDraw = res;
  		updateConnectorShapes(res, true, false);
  	}
  };
  export const resetLineHighlights = () => {
  	const res = clone(linesToDraw);
  	res.map((ln) => {
  		ln.line.opacity = 1;
  		ln.circle.opacity = 1;
  		ln.arrow.opacity = 1;
  		ln.circle.zIndex = 999;
  		ln.line.zIndex = 0;
  		ln.arrow.zIndex = 1;
  		return ln;
  	});
  	linesToDraw = res;
  	updateConnectorShapes(res, false, false);
  };
  let stageConfig = { draggable: true };
</script>

<Stage
  class={clsx("focus:cursor-grabbing relative")}
  getHandler={(handle) => {
  	stage = handle;
  	dispatch("init", handle);
  }}
  bind:config={stageConfig}
>
  {#if isDesktop}
    <Minimap
      bind:this={minimapRef}
      on:dragmove={(e) => {
      	/**
         * This logic allows us to create an interactive minimap.
         * Users can change position of the canvas based on what they want to look at
         * by dragging the viewbox rectangle.
         */
      	const config = e.detail;
      	const scale = stage.scaleX();
      	// The reason we are using 'negative' values is the rect
      	// inverts the rect position when dragged.
      	const newXPos = -config.x * scale;
      	const newYPos = -config.y * scale;
      	gridLayerRef?.updateGrid();
      	stage.position({
      		x: newXPos,
      		y: newYPos,
      	});
      	stage.batchDraw();
      }}
    />
  {/if}
  <div
    class="md:items-center absolute bottom-0 w-full gap-4 flex-col-reverse md:flex-row flex justify-between mb-3 px-5"
  >

  {#if legend.length > 0}
    <div
      class="absolute right-5 md:left-2 bottom-[65px] md:relative md:bottom-0"
    >
      <Legend
        {legend}
        on:reset-highlight={() => {
        	resetLineHighlights();
        	dispatch("highlight-nodes", []);
        }}
        on:highlight={(e) => {
        	dispatch("highlight-nodes", e.detail);
        }}
      />
    </div>
	{:else}
	<div />
  {/if}
    <div class="bg-gray-100 shadow rounded px-2">
      New changes may take up to 15 minutes to reflect
    </div>
    <div class="flex items-center">
      <button
        class="help-text mr-2"
        on:click={() => {
        	// const node = ($datastore?.konvaConnectableNodes || [])[0];
        	// if (node) {
        	// 	handleRepositionStage(node.x, node.y, 0, 0);
        	// } else {
        	stage.to({
        		x: 0,
        		y: 0,
        		duration: 0.5,
        	});
        	// }
        }}
      >
        <Icon src="/assets/images/focus-center.svg" alt="center" width="20" />
      </button>
	  <Settings />
    </div>
  </div>
  <GridLayer bind:this={gridLayerRef} />
  <Layer
    bind:handle={layer}
  >
    <slot />
  </Layer>
</Stage>
