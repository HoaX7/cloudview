<script lang="ts">
  import type {
  	ApiGatewayV2IntegrationProps,
  	ELBV2Props,
  	ResourceDataReturnType,
  } from "$src/customTypes/ervices;
  import Group from "$src/lib/components/common/KonvaCanvas/Group.svelte";
  import Image from "$src/lib/components/common/KonvaCanvas/Image.svelte";
  import Text from "$src/lib/components/common/KonvaCanvas/Text.svelte";
  import Datastore from "$src/store/data";
  import { createEventDispatcher, onMount } from "svelte";
  import ElbV2Data from "./elbV2Data.svelte";
  import { delay } from "$src/helpers";
  import { AWS_SERVICES, LEGEND_NAMES } from "$src/helpers/constants";
  import { getProportions, truncateResourceLabel } from "$src/helpers/konva/index";
  import type { HighLightProps, LegendProps } from "$src/customTypes/konva";
  import { COLOR_SCHEME } from "$src/colorConfig";
  import Rect from "$src/lib/components/common/KonvaCanvas/Rect.svelte";
  import type Konva from "konva";
  import type { GroupConfig } from "konva/lib/Group";
  import { getImageRect } from "../shapeCache";
  import ServiceGroupWithLabel from "../ServiceGroupWithLabel.svelte";
  import type { CloudFrontProps } from "$src/customTypes/aws/cloudfront";
  import PreviewData from "../../views/previewData.svelte";
  import StatusIcon from "../../views/statusIcon.svelte";
  import KonvaStore from "$src/store/konva";

  export let data: ELBV2Props;
  export let externalGroup: ResourceDataReturnType;
  export let setLegend: (legend: LegendProps[]) => void;
  export let highlights: HighLightProps;
  export let idx: number;
  const datastore = Datastore.getDatastore();
  const konvastore = KonvaStore.getStore();

  let legend: LegendProps[] = [];

  const apigateways = externalGroup.find(
  	(et) => et.name === AWS_SERVICES.APIGATEWAYV2
  ) as {
	name: string;
	result: ApiGatewayV2IntegrationProps[];
  } | undefined;
  const cdns = externalGroup.find((et) => et.name === AWS_SERVICES.CLOUDFRONT) as {
	name: string;
	result: CloudFrontProps;
  } | undefined;
  let offset = 0;
  if (apigateways) {
  	offset = apigateways.result.length;
  }
  if (cdns) {
  	offset = offset + cdns.result.Items.length;
  }

  let imageData = (data?.LoadBalancers || []).map((lb, i) => {
  	legend.push({
  		name: LEGEND_NAMES.VPC,
  		id: lb.VpcId,
  		highlight: [ lb.DNSName ],
  		colors: [ COLOR_SCHEME.LOADBALANCER ],
  		count: 1,
  	});
  	lb.SecurityGroups.map((sg) => {
  		legend.push({
  			name: LEGEND_NAMES.SECURITY_GROUP,
  			id: sg,
  			highlight: [ lb.DNSName ],
  			colors: [ COLOR_SCHEME.LOADBALANCER ],
  			count: 1,
  		});
  	});
  	let x = 0,
  		y = 0;
  	if ($datastore.konvaConnectableNodes) {
  		const node = $datastore.konvaConnectableNodes.find(
  			(nd) => nd.id === lb.DNSName
  		);
  		if (node) {
  			(x = node.x), (y = node.y);
  		}
  	}
  	if (x === 0 || y === 0) {
  		const proportions = getProportions(offset, i, "external");
  		x = proportions.x;
  		y = proportions.y;
  	}
  	// 55 is imagewidth plus padding
	  const midX = ($konvastore.externalBoundingRect.width / 2) + $konvastore.externalBoundingRect.x - 55;
  	return {
  		config: {
  			x: midX,
  			y,
  			id: lb.DNSName,
  			draggable: true,
  			label: `Load Balancer ${lb.CanonicalHostedZoneId}`,
  		} as GroupConfig,
  		id: lb.DNSName,
  		name: truncateResourceLabel(lb.LoadBalancerName),
  		data: lb,
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

  if (legend.length > 0) setLegend(legend);

  const dispatch = createEventDispatcher();

  const nodeConfigs = imageData.map((it) => it.config);

  let imageEl: any = null;
  onMount(async () => {
  	const img = document.createElement("img");
  	img.src = "/assets/images/aws/elbv2.png";
  	img.onload = () => {
  		imageEl = img;
  	};

  	dispatch("initialPosition", nodeConfigs);
  });

  const state: any = {
  	showModal: false,
  	data: null,
  	showPreview: false,
  	previewData: null,
  	previewProportions: {
  		x: 0,
  		y: 0
  	}
  };
  const imageWidth = 80;
  const imageHeight = 80;
</script>

<ElbV2Data
  showModal={state.showModal}
  closeModal={async () => {
  	state.showModal = false;
  	await delay(700);
  	state.data = null;
  }}
  data={state.data}
/>

{#if state.showPreview && state.previewData}
<PreviewData proportions={state.previewProportions} data={state.previewData} color={COLOR_SCHEME.LOADBALANCER} />
{/if}

<ServiceGroupWithLabel
  label={{
  	text: "Loadbalancers",
  	fill: COLOR_SCHEME.LOADBALANCER,
  }}
  borderColor={COLOR_SCHEME.LOADBALANCER}
  externalService
  {idx}
>
  {#each imageData as item (item.id)}
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
      on:dragend={() => {
      	dispatch("dragend", item.config);
      }}
      on:mouseenter={(e) => {
      	const targets = $datastore.konvaTargetFromNodes.filter((tg) => tg.to.includes(item.config?.id || ""));
      	dispatch("mouseenter", {
      		id: item.config.id,
      		highlights: targets,
      		extras: targets.map((tg) => tg.id)
      	});

      	state.previewData = [ {
      		name: "Scheme",
      		value: item.data.Scheme
      	}, {
      		name: "State",
      		value: item.data.State?.Code || "Unknown"
      	} ];
      	state.previewProportions = {
      		x: (item.config.x || 0) - 220,
      		y: item.config.y
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
	  getHandler={(handle) => {
	  	const rect = getImageRect({ fill: COLOR_SCHEME.LOADBALANCER });
	  	handle.add(rect);
	  }}
    >
	<StatusIcon status={item.data.State?.Code === "active" ? "RUNNING" : "STOPPED"} />
      <Image config={{ image: imageEl }} position={{
      	draggable: false,
      	x: 12,
      	y: 12 
      }} />
      <Text
        config={{
        	text: item.name || "Load Balancer",
        	draggable: false,
        	y: -20,
        	x: 0,
        	listening: false,
        	// fill: COLOR_SCHEME.LOADBALANCER,
        	fontStyle: "bold",
        }}
      />
    </Group>
  {/each}
</ServiceGroupWithLabel>
