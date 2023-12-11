<script lang="ts">
  import type { LambdaFunctionProps } from "$src/customTypes/Services";
  import { createEventDispatcher, onMount } from "svelte";
  import Datastore from "$src/store/data";
  import { delay } from "$src/helpers";
  import Group from "../../../common/KonvaCanvas/Group.svelte";
  import Image from "../../../common/KonvaCanvas/Image.svelte";
  import Text from "../../../common/KonvaCanvas/Text.svelte";
  import LambdaData from "./lambdaData.svelte";
  import { getProportions } from "$src/helpers/konva/index";
  import Rect from "$src/lib/components/common/KonvaCanvas/Rect.svelte";
  import type Konva from "konva";
  import { COLOR_SCHEME } from "$src/colorConfig";
  import type { GroupConfig } from "konva/lib/Group";
  import type { HighLightProps } from "$src/customTypes/Konva";
  import { getImageRect } from "../shapeCache";
  import ServiceGroupWithLabel from "../ServiceGroupWithLabel.svelte";
  import PreviewData from "../../views/previewData.svelte";
  import Circle from "$src/lib/components/common/KonvaCanvas/Circle.svelte";
  import StatusIcon from "../../views/statusIcon.svelte";

  export let data: LambdaFunctionProps;
  export let idx: number = 0;
  export let highlights: HighLightProps;

  const datastore = Datastore.getDatastore();

  const dispatch = createEventDispatcher();

  let imageData = data.Functions.map((fn, i) => {
  	let x = 0,
  		y = 0;
  	if ($datastore.konvaConnectableNodes) {
  		const node = $datastore.konvaConnectableNodes.find(
  			(nd) => nd.id === fn.FunctionArn
  		);
  		if (node) {
  			(x = node.x), (y = node.y);
  		}
  	}
  	if (x === 0 || y === 0) {
  		const proportions = getProportions(idx, i, "internal");
  		x = proportions.x;
  		// Navbar height is 64px and should not be included while placing icons on the canvas
  		y = proportions.y;
  	}
  	return {
  		text: fn.FunctionName,
  		data: fn,
  		config: {
  			draggable: true,
  			id: fn.FunctionArn,
  			x,
  			y,
  		} as GroupConfig,
  	};
  });
  $: {
  	imageData = imageData.map((it) => {
  		const node = (highlights?.nodes || []).find(
  			(nd) => nd?.includes(it.config?.id || "") || nd === it.config?.id
  		);
  		if (highlights.nodes && highlights.nodes.length > 0 && !node) {
  			it.config.opacity = 0.3;
  			return it;
  		}
  		it.config.opacity = 1;
  		return it;
  	});
  }

  const nodeConfigs = imageData.map((img) => img.config);

  // lambda image
  let imageEl: any = null;
  onMount(() => {
  	const img = document.createElement("img");
  	img.src = "/assets/images/aws/lambda.png";
  	img.onload = () => {
  		imageEl = img;
  	};

  	/**
     *  This event returns initial vector2D positions
     * of instance nodes. This event gives us pre-computed position of the
     * node to be connected to by another instance.
     *
     * This is a common event that can be used by other components
     * to show connecting arrow.
     */
  	dispatch("initialPosition", nodeConfigs);
  });

  const state: any = {
  	showModal: false,
  	data: null,
  	previewData: null,
  	showPreview: false,
  	previewProportions: {
  		x: 0,
  		y: 0
  	}
  };
  const imageWidth = 80;
  const imageHeight = 80;
</script>

<LambdaData
  showModal={state.showModal}
  closeModal={async () => {
  	state.showModal = false;
  	await delay(700);
  	state.data = null;
  }}
  data={state.data}
/>
{#if state.showPreview && state.previewData}
<PreviewData data={state.previewData} proportions={state.previewProportions} color={COLOR_SCHEME.SERVERLESS} />
{/if}
<ServiceGroupWithLabel borderColor={COLOR_SCHEME.SERVERLESS} label={{
	text: "Lambda Functions",
	fill: COLOR_SCHEME.SERVERLESS
}} {idx}>
  {#each imageData as item, index (index)}
    <Group
      bind:config={item.config}
      on:click={() => {
      	dispatch("click", {
      		...item.config,
      		width: imageWidth,
      		height: imageHeight + 20,
      	});
      	state.data = item.data;
      	state.showModal = true;
      }}
      on:dragmove={() => {
      	dispatch("dragmove", item.config);
      }}
      on:mouseenter={(e) => {
      	const targets = $datastore.konvaTargetFromNodes.filter((tg) =>
      		tg.to.includes(item.config?.id || "")
      	);
      	dispatch("mouseenter", {
      		id: item.config.id,
      		highlights: targets,
      		extras: targets.map((tg) => tg.id),
      	});
      	state.previewData = [ {
      		name: "Run Time",
      		value: item.data.Runtime
      	}, {
      		name: "Storage",
      		value: item.data.EphemeralStorage.Size + " MB"
      	}, {
      		name: "Memory",
      		value: item.data.MemorySize + " MB"
      	}, {
      		name: "State",
      		value: item.data.State || "Unknown"
      	} ];
      	state.previewProportions = {
      		x: (item.config.x || 0) - (imageWidth / 2),
      		y: (item.config.y || 0) - (imageHeight + 40)
      	};
      	state.showPreview = true;
      }}
	  on:mouseout={() => {
	  	state.showPreview = false;
	  	state.previewProportions = {
	  		x: 0,
	  		y: 0
	  	};
	  	state.previewData = null;
	  }}
      on:mouseleave={(e) => {
      	dispatch("mouseleave", e);
      }}
      on:dragend={() => {
      	dispatch("dragend", item.config);
      }}
	  getHandler={(handle) => {
	  	const rect = getImageRect({ fill: COLOR_SCHEME.SERVERLESS });
	  	handle.add(rect);
	  }}
    >
	<StatusIcon status={"UNKNOWN"} />
      <Image
        config={{ image: imageEl }}
        position={{
        	draggable: false,
        	x: 12,
        	y: 12,
        }}
      />
      <Text
        config={{
        	text: item.text,
        	draggable: false,
        	y: -20,
        	x: 0,
        	listening: false,
        	fontStyle: "bold",
        	fill: COLOR_SCHEME.SERVERLESS,
        }}
      />
    </Group>
  {/each}
</ServiceGroupWithLabel>