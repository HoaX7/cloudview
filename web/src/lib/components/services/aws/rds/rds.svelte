<script lang="ts">
  import type { RDSProps } from "$src/customTypes/aws/rds";
  import Group from "$src/lib/components/common/KonvaCanvas/Group.svelte";
  import Image from "$src/lib/components/common/KonvaCanvas/Image.svelte";
  import Text from "$src/lib/components/common/KonvaCanvas/Text.svelte";
  import Datastore from "$src/store/data";
  import { createEventDispatcher, onMount } from "svelte";
  import RdsData from "./rdsData.svelte";
  import { delay } from "$src/helpers";
  import { getProportions, truncateResourceLabel } from "$src/helpers/konva/index";
  import { COLOR_SCHEME } from "$src/colorConfig";
  import type { HighLightProps, LegendProps } from "$src/customTypes/Konva";
  import { LEGEND_NAMES } from "$src/helpers/constants";
  import type { GroupConfig } from "konva/lib/Group";
  import { getImageRect } from "../shapeCache";
  import ServiceGroupWithLabel from "../ServiceGroupWithLabel.svelte";
  import PreviewData from "../../views/previewData.svelte";

  export let data: RDSProps;
  export let idx: number = 0;
  export let setLegend: (data: LegendProps[]) => void;
  export let highlights: HighLightProps;
  const datastore = Datastore.getDatastore();
  const legend: LegendProps[] = [];

  let imageData = (data?.DBInstances || []).map((instance, i) => {
  	if (instance.DBSubnetGroup.Subnets.length > 0) {
  		legend.push(
  			...instance.DBSubnetGroup.Subnets.map((subnet) => ({
  				id: subnet.SubnetIdentifier,
  				name: LEGEND_NAMES.SUBNET,
  				highlight: [ instance.DBInstanceArn ],
  				count: 1,
  				colors: [ COLOR_SCHEME.RDB ],
  			}))
  		);
  	}
  	if (instance.DBSubnetGroup.VpcId) {
  		legend.push({
  			id: instance.DBSubnetGroup.VpcId,
  			name: LEGEND_NAMES.VPC,
  			highlight: [ instance.DBInstanceArn ],
  			count: 1,
  			colors: [ COLOR_SCHEME.RDB ],
  		});
  	}
  	if (instance.VpcSecurityGroups.length > 0) {
  		legend.push(
  			...instance.VpcSecurityGroups.map((group) => ({
  				id: group.VpcSecurityGroupId,
  				name: LEGEND_NAMES.SECURITY_GROUP,
  				count: 1,
  				highlight: [ instance.DBInstanceArn ],
  				colors: [ COLOR_SCHEME.RDB ],
  			}))
  		);
  	}
  	let x = 0,
  		y = 0;
  	if ($datastore.konvaConnectableNodes) {
  		const node = $datastore.konvaConnectableNodes.find(
  			(nd) => nd.id === instance.DBInstanceArn
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
  		id: instance.DBInstanceArn,
  		name: truncateResourceLabel(instance.DBInstanceIdentifier),
  		data: instance,
  		config: {
  			draggable: true,
  			id: instance.DBInstanceArn,
  			x,
  			y,
  			label: `RDS ${instance.DBInstanceIdentifier}`,
  		} as GroupConfig,
  	};
  });
  $: {
  	imageData = imageData.map((it) => {
  		const node = (highlights?.nodes || []).find(
  			(nd) =>
  				nd?.includes(it.config?.id || "") || nd === it.config?.id
  		);
  		if (highlights.nodes && highlights.nodes.length > 0 && !node) {
  			it.config.opacity = 0.3;
  			return it;
  		}
  		it.config.opacity = 1;
  		return it;
  	});
  }
  if (legend.length > 0) setLegend(legend);
  const dispatch = createEventDispatcher();

  const nodeConfigs = imageData.map((img) => img.config);

  // Ec2 image
  let imageEl: any = null;
  onMount(() => {
  	const img = document.createElement("img");
  	img.src = "/assets/images/aws/rds.png";
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
  	showPreview: false,
  	previewData: null,
  	previewPosition: {
  		x: 0,
  		y: 0,
  	},
  };
  const imageWidth = 80;
  const imageHeight = 80;
</script>

<RdsData
  showModal={state.showModal}
  data={state.data}
  closeModal={async () => {
  	state.showModal = false;
  	await delay(700);
  	state.data = null;
  }}
/>
{#if state.showPreview && state.previewData}
  <PreviewData
    proportions={state.previewPosition}
    data={state.previewData}
    color={COLOR_SCHEME.RDB}
  />
{/if}
<ServiceGroupWithLabel
  label={{
  	text: "Relational DB Services",
  	fill: COLOR_SCHEME.RDB,
  }}
  borderColor={COLOR_SCHEME.RDB}
  {idx}
>
  {#each imageData as item (item.id)}
    <Group
      bind:config={item.config}
      on:mouseenter={(e) => {
      	const targets = $datastore.konvaTargetFromNodes.filter((tg) =>
      		tg.to.includes(item.config?.id || "")
      	);
      	dispatch("mouseenter", {
      		id: item.config.id,
      		highlights: targets,
      		extras: targets.map((tg) => tg.id),
      	});
      	state.previewPosition = {
      		x: (item.config.x || 0) - imageWidth / 2,
      		y: (item.config.y || 0) - (imageHeight + 40),
      	};
      	state.previewData = [
      		{
      			name: "Engine",
      			value: `${item.data.Engine} x${item.data.EngineVersion}`,
      		},
      		{
      			name: "Class",
      			value: item.data.DBInstanceClass,
      		},
      		{
      			name: "Storage",
      			value: item.data.AllocatedStorage + " GB",
      		},
      		{
      			name: "Public Access",
      			value: item.data.PubliclyAccessible ? "Yes" : "No",
      		},
      	];
      	state.showPreview = true;
      }}
      on:mouseout={() => {
      	state.showPreview = false;
      	state.previewData = null;
      	state.previewProportions = {
      		x: 0,
      		y: 0,
      	};
      }}
      on:mouseleave={(e) => {
      	dispatch("mouseleave", e);
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
      on:dragend={() => {
      	dispatch("dragend", item.config);
      }}
      getHandler={(handle) => {
      	const rect = getImageRect({ fill: COLOR_SCHEME.RDB });
      	handle.add(rect);
      }}
    >
      <Image
        config={{ image: imageEl }}
        position={{
        	draggable: false,
        	y: 12,
        	x: 12,
        }}
      />
      <Text
        config={{
        	text: item.name || "Rds",
        	draggable: false,
        	y: -20,
        	x: 0,
        	listening: false,
        	fontStyle: "bold",
        	// fill: COLOR_SCHEME.RDB,
        }}
      />
    </Group>
  {/each}
</ServiceGroupWithLabel>
