<script lang="ts">
  import { createEventDispatcher, onMount } from "svelte";
  import Group from "../../../common/KonvaCanvas/Group.svelte";
  import Image from "../../../common/KonvaCanvas/Image.svelte";
  import { delay } from "$src/helpers";
  import Text from "../../../common/KonvaCanvas/Text.svelte";
  import Datastore from "$src/store/data";
  import Ec2Data from "./ec2Data.svelte";
  import type { Ec2Props } from "$src/customTypes/services";
  import { getProportions } from "$src/helpers/konva/index";
  import { LEGEND_NAMES, STATUS_COLORS } from "$src/helpers/constants";
  import type { HighLightProps, LegendProps } from "$src/customTypes/konva";
  import { COLOR_SCHEME } from "$src/colorConfig";
  import type { GroupConfig } from "konva/lib/Group";
  import { getImageRect } from "../shapeCache";
  import ServiceGroupWithLabel from "../ServiceGroupWithLabel.svelte";
  import PreviewData from "../../views/previewData.svelte";
  import Circle from "$src/lib/components/common/KonvaCanvas/Circle.svelte";
  import StatusIcon from "../../views/statusIcon.svelte";

  export let data: Ec2Props;
  export let projectId: string;
  export let providerAccountId: string;
  export let region: string;
  export let idx: number = 0;
  export let setLegend: (legend: LegendProps[]) => void;
  export let highlights: HighLightProps;

  const datastore = Datastore.getDatastore();
  const legend: LegendProps[] = [];

  // ec2 is an array
  let imageData = (data?.Reservations || [])
  	.map((reservation, i) => {
  		return reservation.Instances.map((instance) => {
  			if (instance.SubnetId) {
  				legend.push({
  					id: instance.SubnetId,
  					name: LEGEND_NAMES.SUBNET,
  					count: 1,
  					highlight: [ instance.PublicIpAddress ],
  					colors: [ COLOR_SCHEME.VM ],
  				});
  			}
  			if (instance.VpcId) {
  				legend.push({
  					id: instance.VpcId,
  					name: LEGEND_NAMES.VPC,
  					count: 1,
  					highlight: [ instance.PublicIpAddress ],
  					colors: [ COLOR_SCHEME.VM ],
  				});
  			}
  			if (instance.SecurityGroups.length > 0) {
  				legend.push(
  					...instance.SecurityGroups.map((it) => ({
  						id: it.GroupId,
  						name: LEGEND_NAMES.SECURITY_GROUP,
  						count: 1,
  						highlight: [ instance.PublicIpAddress ],
  						colors: [ COLOR_SCHEME.VM ],
  					}))
  				);
  			}
  			// Save node positions instead of random position each time.
  			let x = 0,
  				y = 0;
  			if ($datastore.konvaConnectableNodes) {
  				const node = $datastore.konvaConnectableNodes.find(
  					(nd) => nd.id === instance.PublicIpAddress
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
  				text: instance.PublicIpAddress,
  				instance,
  				config: {
  					draggable: false,
  					// We will use the public ip address as ID
  					// to draw connecting arrows from api-gateway
  					id: instance.PublicIpAddress,
  					x,
  					y,
  					label: `Ec2 ${instance.PublicIpAddress}`,
  				} as GroupConfig,
  			};
  		});
  	})
  	.flat();

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

  if (legend.length > 0) {
  	setLegend(legend);
  }

  const dispatch = createEventDispatcher();

  const nodeConfigs = imageData.map((img) => img.config);

  // Ec2 image
  let imageEl: any = null;
  const imageWidth = 80;
  const imageHeight = 80;
  onMount(() => {
  	const img = document.createElement("img");
  	img.src = "/assets/images/aws/ec2.png";
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
  	instance: null,
  	showModal: false,
  	volume: null,
  	showPreview: false,
  	previewProportions: {
  		x: 0,
  		y: 0,
  	},
  	previewData: null,
  };

  // ##NOTICE - Elements are currently not draggable. If you wish to make them draggable,
  // make sure to also update the position of the image background rectangle.

  const resetPreview = () => {
  	state.showPreview = false;
      	state.previewProportions = {
      		x: 0,
      		y: 0,
      	};
      	state.previewData = null;
  };
</script>

<Ec2Data
  {projectId}
  {providerAccountId}
  {region}
  instance={state.instance}
  onClose={async () => {
  	state.showModal = false;

  	// setting a delay to allow drawer to close
  	await delay(700);
  	state.instance = null;
  }}
  showModal={state.showModal}
  volume={state.volume}
/>
{#if state.showPreview}
  <PreviewData
    color={COLOR_SCHEME.VM}
    proportions={state.previewProportions}
    data={state.previewData}
  />
{/if}
<ServiceGroupWithLabel
  label={{
  	text: "Ec2 Instances",
  	fill: COLOR_SCHEME.VM,
  }}
  borderColor={COLOR_SCHEME.VM}
  {idx}
>
  {#each imageData as item, index (index)}
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
      	state.previewData = [
      		{
      			name: "Public Ip",
      			value: item.instance.PublicIpAddress,
      		},
      		{
      			name: "vCpu",
      			value: item.instance.CpuOptions.CoreCount,
      		},
      		{
      			name: "RAM",
      			value: item.instance.InstanceType,
      		},
      		{
      			name: "State",
      			value: item.instance.State.Name || "",
      		},
      	];
      	state.previewProportions = {
      		x: (item.config.x || 0) - imageWidth / 2,
      		y: (item.config.y || 0) - (imageHeight + 40),
      	};
      	state.showPreview = true;
      }}
      on:mouseleave={(e) => {
      	dispatch("mouseleave", e);
      	console.log("Mouse leave...");
      }}
      on:dragend={() => {
      	dispatch("dragend", item.config);
      }}
      on:mouseout={resetPreview}
      on:click={() => {
      	resetPreview();
      	dispatch("click", {
      		...item.config,
      		width: imageWidth,
      		height: imageHeight + 20,
      	});
      	state.instance = item.instance;
      	const volume = (data?.Volumes || []).find(
      		(vol) =>
      			vol.VolumeId === item.instance.BlockDeviceMappings[0].Ebs?.VolumeId
      	);
      	if (volume) state.volume = volume;

      	state.showModal = true;
      }}
      on:dragmove={() => {
      	dispatch("dragmove", item.config);
      }}
      getHandler={(handle) => {
      	const rect = getImageRect({ fill: COLOR_SCHEME.VM });
      	handle.add(rect);
      }}
    >
      <StatusIcon
        status={item.instance.State.Name === "running"
        	? "RUNNING"
        	: item.instance.State.Name === "stopped"
        		? "STOPPED"
        		: "UNKNOWN"}
      />
      <Image
        config={{ image: imageEl }}
        position={{
        	draggable: false,
        	x: 12,
        	y: 12,
        	listening: false,
        }}
      />
      <Text
        config={{
        	text: item.text || "Ec2",
        	draggable: false,
        	y: -20,
        	x: 0,
        	listening: false,
        	// fill: COLOR_SCHEME.VM,
        	fontStyle: "bold",
        }}
      />
    </Group>
  {/each}
</ServiceGroupWithLabel>
