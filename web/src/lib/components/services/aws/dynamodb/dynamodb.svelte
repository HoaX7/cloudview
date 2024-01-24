<script lang="ts">
  import { COLOR_SCHEME } from "$src/colorConfig";
  import type { HighLightProps } from "$src/customTypes/Konva";
  import type { DynamoDBProps } from "$src/customTypes/Services";
  import { getProportions, truncateResourceLabel } from "$src/helpers/konva";
  import Group from "$src/lib/components/common/KonvaCanvas/Group.svelte";
  import Rect from "$src/lib/components/common/KonvaCanvas/Rect.svelte";
  import Datastore from "$src/store/data";
  import type Konva from "konva";
  import type { GroupConfig } from "konva/lib/Group";
  import { createEventDispatcher, onMount } from "svelte";
  import { getImageRect } from "../shapeCache";
  import Image from "$src/lib/components/common/KonvaCanvas/Image.svelte";
  import Text from "$src/lib/components/common/KonvaCanvas/Text.svelte";
  import ServiceGroupWithLabel from "../ServiceGroupWithLabel.svelte";
  import Table from "$src/lib/components/common/Table/Table.svelte";
  import DynamodbData from "./dynamodbData.svelte";
  import PreviewData from "../../views/previewData.svelte";
  import { bytesToMegaBytes } from "$src/helpers";
  import StatusIcon from "../../views/statusIcon.svelte";

  let datastore = Datastore.getDatastore();
  export let data: DynamoDBProps;
  export let idx: number = 0;
  export let highlights: HighLightProps;
  console.log({ data });
  let imageData = (data || []).map(({ Table }, i) => {
  	let x = 0,
  		y = 0;
  	// TODO - change the identifier to access node.
  	// if ($datastore.konvaConnectableNodes) {
  	// 	const node = $datastore.konvaConnectableNodes.find(
  	// 		(nd) => nd.id === name
  	// 	);
  	// 	if (node) {
  	// 		(x = node.x), (y = node.y);
  	// 	}
  	// } else {
  	// 	const proportions = getProportions(idx, i, "internal");
  	// 	x = proportions.x;
  	// 	y = proportions.y;
  	// }

  	const proportions = getProportions(idx, i, "internal");
  		x = proportions.x;
  		y = proportions.y;
  	return {
  		config: {
  			draggable: false,
  			x,
  			y,
  			id: Table.TableId,
  			label: Table.TableName,
  		} as GroupConfig,
  		data: Table,
  		text: truncateResourceLabel(Table.TableName),
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

  const dispatch = createEventDispatcher();

  const nodeConfigs = imageData.map((img) => img.config);

  let imageEl: any = null;
  const imageWidth = 80;
  const imageHeight = 80;
  onMount(() => {
  	const img = document.createElement("img");
  	img.src = "/assets/images/aws/dynamodb.png";
  	img.onload = () => {
  		imageEl = img;
  	};

  	dispatch("initialPosition", nodeConfigs);
  });

  const state: any = {
  	data: null,
  	showModal: false,
  	previewData: null,
  	showPreview: false,
  	previewProportions: {
  		x: 0,
  		y: 0
  	}
  };

</script>

<DynamodbData 
    showModal={state.showModal}
    data={state.data}
    closeModal={() => {
    	state.showModal = false;
    	state.data = null;
    }}
/>

{#if state.showPreview && state.previewData}
<PreviewData proportions={state.previewProportions} data={state.previewData} color={COLOR_SCHEME.DB} />
{/if}

<ServiceGroupWithLabel
    borderColor={COLOR_SCHEME.DB}
    label={{
    	fill: COLOR_SCHEME.DB,
    	text: "Dynamo DB"
    }}
    {idx}
>
  {#each imageData as item, index (index)}
    <Group
      bind:config={item.config}
      on:mouseenter={(e) => {
      	const targets = $datastore.konvaTargetFromNodes.filter((tg) => tg.to.includes(item.config?.id || ""));
      	dispatch("mouseenter", {
      		id: item.config.id,
      		highlights: targets,
      		extras: targets.map((tg) => tg.id)
      	});
      	state.previewData = [ {
      		name: "Size",
      		value: bytesToMegaBytes(item.data.TableSizeBytes || 0) + " MB"
      	}, {
      		name: "Item Count",
      		value: (item.data.ItemCount || 0).toString()
      	}, {
      		name: "Status",
      		value: item.data.TableStatus || "Unknown"
      	} ];
      	state.previewProportions = {
      		x: (item.config.x || 0) - (imageWidth / 2),
      		y: (item.config.y || 0) - (imageHeight + 40)
      	};
      	state.showPreview = true;
      }}
      on:mouseout={() => {
      	state.showPreview = false;
      	state.previewData = null;
      	state.previewProportions = {
      		x: 0,
      		y: 0
      	};
      }}
      on:mouseleave={(e) => {
      	dispatch("mouseleave", e);
      }}
      on:dragend={() => {
      	dispatch("dragend", item.config);
      }}
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
	  getHandler={(handle) => {
	  	const rect = getImageRect({ fill: COLOR_SCHEME.DB });
	  	handle.add(rect);
	  }}
    >
    <StatusIcon status={item.data.TableStatus === "ACTIVE" ? "RUNNING" : "STOPPED"} />
      <Image
        config={{ image: imageEl }}
        position={{
        	draggable: false,
        	x: 16,
        	y: 12,
        	listening: false,
        }}
      />
      <Text
        config={{
        	text: item.text || "Dynamodb",
        	draggable: false,
        	y: -20,
        	x: 0,
        	listening: false,
        	// fill: COLOR_SCHEME.DB,
        	fontStyle: "bold",
        }}
      />
    </Group>
  {/each}
</ServiceGroupWithLabel>