<script lang="ts">
  import type { EFSProps } from "$src/customTypes/Services";
  import Group from "$src/lib/components/common/KonvaCanvas/Group.svelte";
  import Image from "$src/lib/components/common/KonvaCanvas/Image.svelte";
  import Text from "$src/lib/components/common/KonvaCanvas/Text.svelte";
  import Datastore from "$src/store/data";
  import { createEventDispatcher, onMount } from "svelte";
  import EfsData from "./efsData.svelte";
  import { delay } from "$src/helpers";
  import { getProportions, truncateResourceLabel } from "$src/helpers/konva/index";
  import { COLOR_SCHEME } from "$src/colorConfig";
  import type { GroupConfig } from "konva/lib/Group";
  import type { HighLightProps } from "$src/customTypes/Konva";
  import { getImageRect } from "../shapeCache";
  import ServiceGroupWithLabel from "../ServiceGroupWithLabel.svelte";

  export let data: EFSProps;
  export let idx: number = 0;
  export let highlights: HighLightProps;
  const datastore = Datastore.getDatastore();

  let imageData = (data?.FileSystems || []).map((fs, i) => {
  	let x = 0,
  		y = 0;
  	if ($datastore.konvaConnectableNodes) {
  		const node = $datastore.konvaConnectableNodes.find(
  			(nd) => nd.id === fs.FileSystemId
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
  	return {
  		config: {
  			x,
  			y,
  			draggable: true,
  			id: fs.FileSystemId,
  			label: `EFS ${fs.FileSystemId}`,
  		} as GroupConfig,
  		id: fs.FileSystemId,
  		fileSystem: fs,
  		name: truncateResourceLabel(fs.Name) || "Efs",
  		data: fs,
  	};
  });
  $: {
  	imageData = imageData.map((it) => {
  		const node = (highlights?.nodes || []).find((nd) => nd?.includes(it.config?.id || "") || nd === it.config?.id);
  		if (highlights.nodes && highlights.nodes.length > 0 && !node) {
  			it.config.opacity = .3;
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
  onMount(() => {
  	const img = document.createElement("img");
  	img.src = "/assets/images/aws/efs.png";
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
  	data: null,
  	showModal: false,
  };
  const imageWidth = 80;
  const imageHeight = 80;
</script>

<EfsData
  showModal={state.showModal}
  data={state.data}
  closeModal={async () => {
  	state.showModal = false;
  	await delay(700);
  	state.data = null;
  }}
/>

<ServiceGroupWithLabel label={{
	text: "Elastic File Systems",
	fill: COLOR_SCHEME.FILE_SYSTEMS
}} borderColor={COLOR_SCHEME.FILE_SYSTEMS} {idx}>
  {#each imageData as item (item.id)}
    <Group
      bind:config={item.config}
      on:dragmove={() => {
      	dispatch("dragmove", item.config);
      }}
      on:mouseenter={(e) => {
      	const targets = $datastore.konvaTargetFromNodes.filter((tg) => tg.to.includes(item.config?.id || ""));
      	dispatch("mouseenter", {
      		id: item.config.id,
      		highlights: targets,
      		extras: targets.map((tg) => tg.id)
      	});
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
	  getHandler={(handle) => {
	  	const rect = getImageRect({ fill: COLOR_SCHEME.FILE_SYSTEMS });
	  	handle.add(rect);
	  }}
    >
      <Image
        config={{ image: imageEl }}
        position={{
        	draggable: false,
        	x: 6,
        	y: 8,
        }}
      />
      <Text
        config={{
        	text: item.name || "Efs",
        	draggable: false,
        	y: -20,
        	x: 0,
        	listening: false,
        	fontStyle: "bold",
        	// fill: COLOR_SCHEME.FILE_SYSTEMS,
        }}
      />
    </Group>
  {/each}
</ServiceGroupWithLabel>