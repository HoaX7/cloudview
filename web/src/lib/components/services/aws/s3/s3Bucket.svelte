<script lang="ts">
  import type { S3Props } from "$src/customTypes/services";
  import Group from "$src/lib/components/common/KonvaCanvas/Group.svelte";
  import Image from "$src/lib/components/common/KonvaCanvas/Image.svelte";
  import Text from "$src/lib/components/common/KonvaCanvas/Text.svelte";
  import Datastore from "$src/store/data";
  import { createEventDispatcher, onMount } from "svelte";
  import S3Data from "./s3Data.svelte";
  import { bytesToMegaBytes, delay } from "$src/helpers";
  import { getProportions, truncateResourceLabel } from "$src/helpers/konva/index";
  import type Konva from "konva";
  import Rect from "$src/lib/components/common/KonvaCanvas/Rect.svelte";
  import { COLOR_SCHEME } from "$src/colorConfig";
  import type { GroupConfig } from "konva/lib/Group";
  import type { HighLightProps } from "$src/customTypes/konva";
  import { getImageRect } from "../shapeCache";
  import ServiceGroupWithLabel from "../ServiceGroupWithLabel.svelte";
  import PreviewData from "../../views/previewData.svelte";

  export let data: S3Props;
  export let idx: number = 0;
  export let highlights: HighLightProps;

  const datastore = Datastore.getDatastore();
  const imageHeight = 80,
  	imageWidth = 80;

  let imageData = (data?.Data?.Buckets || []).map((bucket, i) => {
  	let x = 0,
  		y = 0;
  	if ($datastore.konvaConnectableNodes) {
  		const node = $datastore.konvaConnectableNodes.find(
  			(nd) => nd.id === bucket.Name
  		);
  		if (node) {
  			(x = node.x), (y = node.y);
  		}
  	}
  	if (x === 0 || y === 0) {
  		const proportions = getProportions(idx, i, "internal");
  		x = proportions.x;
  		y = proportions.y;
  	}
  	if (bucket.Name === "zappa-contentai") {
  		bucket.isPublic = true;
  	}
  	return {
  		config: {
  			draggable: true,
  			x,
  			y,
  			id: bucket.Name,
  			label: `S3 ${bucket.Name}`,
  		} as GroupConfig,
  		id: bucket.Name,
  		name: truncateResourceLabel(bucket.Name),
  		data: bucket,
  	};
  });
  $: {
  	imageData = imageData.map((it) => {
  		const node = (highlights?.nodes || []).find((nd) => nd?.includes(it.config?.id || "") || nd === it.config?.id);
  		if (
  			highlights.nodes &&
        highlights.nodes.length > 0 &&
        !node
  		) {
  			it.config.opacity = 0.3;
  			return it;
  		}
  		it.config.opacity = 1;
  		return it;
  	});
  }
  const dispatch = createEventDispatcher();

  const nodeConfigs = imageData.map((img) => img.config);

  // Ec2 image
  let imageEl: any = null;
  let publicImageEl: any = null;
  onMount(() => {
  	const img = document.createElement("img");
  	img.src = "/assets/images/aws/s3.png";
  	img.onload = () => {
  		imageEl = img;
  	};

  	const publicImg = document.createElement("img");
  	publicImg.src = "/assets/images/public.svg";
  	publicImg.onload = () => {
  		publicImageEl = publicImg;
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
  	metric: null,
  	bucket: null,
  	acl: null,
  	previewData: null,
  	showPreview: false,
  	previewProportions: {
  		x: 0,
  		y: 0
  	}
  };

</script>

<S3Data
  showModal={state.showModal}
  metric={state.metric}
  bucket={state.bucket}
  onClose={async () => {
  	state.showModal = false;
  	await delay(700);
  	state.metric = null;
  	state.data = null;
  	state.acl = null;
  }}
  acl={state.acl}
/>
{#if state.showPreview && state.previewData}
<PreviewData proportions={state.previewProportions} data={state.previewData} color={COLOR_SCHEME.OBJECT_STORAGE} />
{/if}
<ServiceGroupWithLabel
  label={{
  	text: "S3 Buckets",
  	fill: COLOR_SCHEME.OBJECT_STORAGE
  }}
  borderColor={COLOR_SCHEME.OBJECT_STORAGE}
  {idx}
>
  {#each imageData as item (item.id)}
    <Group
      bind:config={item.config}
      on:mouseenter={(e) => {
      	const targets = $datastore.konvaTargetFromNodes.filter((tg) => tg.to.includes(item.config?.id || ""));
      	dispatch("mouseenter", {
      		id: item.config.id,
      		highlights: targets,
      		extras: targets.map((tg) => tg.id)
      	});
      	const metric = data.Metrics.find((mt) => mt.Name === item.name);
      	state.previewData = [ {
      		name: "Bucket",
      		value: item.data.Name
      	}, {
      		name: "Size",
      		value: bytesToMegaBytes(metric?.Statistics.Datapoints[0]?.Sum || 0) + " MB"
      	} ];
      	if (item.data.isPublic) {
      		state.previewData.push({
      			name: "Public",
      			value: "Yes"
      		});
      	}
      	state.previewProportions = {
      		x: (item.config.x || 0) - (imageWidth / 2),
      		y: (item.config.y || 0) - (imageHeight + 40)
      	};
      	state.showPreview = true;
      }}
      on:mouseleave={(e) => {
      	dispatch("mouseleave", e);
      }}
	  on:mouseout={() => {
	  	state.showPreview = false;
	  	state.previewPosition = {
	  		x: 0,
	  		y: 0
	  	};
	  	state.previewData = null;
	  }}
      on:click={() => {
      	dispatch("click", {
      		...item.config,
      		width: imageWidth,
      		height: imageHeight + 20,
      	});
      	state.bucket = item.data;
      	const metric = data.Metrics.find((mt) => mt.Name === item.name);
      	const acl = data.ACLList.find((acl) => acl.Bucket === item.name);
      	state.acl = acl?.ACL;
      	state.metric = metric;
      	state.showModal = true;
      }}
      on:dragend={() => {
      	dispatch("dragend", item.config);
      }}
      on:dragmove={() => {
      	dispatch("dragmove", item.config);
      }}
	  getHandler={(handle) => {
	  	const rect = getImageRect({ fill: COLOR_SCHEME.OBJECT_STORAGE });
	  	handle.add(rect);
	  }}
    >
	{#if item.data.isPublic}
	  <Image 
	  	config={{ image: publicImageEl }}
		position={{
			draggable: false,
			x: 3,
			y: 3,
			width: 12,
			height: 12
		}}
	  />
	{/if}
      <Image
        config={{ image: imageEl }}
        position={{
        	draggable: false,
        	x: 13,
        	y: 12,
        }}
      />
      <Text
        config={{
        	text: item.name || "S3",
        	draggable: false,
        	y: -20,
        	x: 0,
        	listening: false,
        	fontStyle: "bold",
        	// fill: COLOR_SCHEME.OBJECT_STORAGE,
        }}
      />
    </Group>
  {/each}
</ServiceGroupWithLabel>
