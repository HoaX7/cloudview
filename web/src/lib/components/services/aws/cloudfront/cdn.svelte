<script lang="ts">
  import type {
  	HighLightProps,
  	TargetFromNodeProps,
  } from "$src/customTypes/Konva";
  import type { CloudFrontProps } from "$src/customTypes/aws/cloudfront";
  import Group from "$src/lib/components/common/KonvaCanvas/Group.svelte";
  import Image from "$src/lib/components/common/KonvaCanvas/Image.svelte";
  import Text from "$src/lib/components/common/KonvaCanvas/Text.svelte";
  import Datastore from "$src/store/data";
  import { createEventDispatcher, onMount, tick } from "svelte";
  import CdnData from "./cdnData.svelte";
  import { delay } from "$src/helpers";
  import type { ApiGatewayV2IntegrationProps, MetricDataReturnType } from "$src/customTypes/Services";
  import { AWS_SERVICES } from "$src/helpers/constants";
  import { getProportions, truncateResourceLabel } from "$src/helpers/konva/index";
  import { COLOR_SCHEME } from "$src/colorConfig";
  import { getImageRect } from "../shapeCache";
  import ServiceGroupWithLabel from "../ServiceGroupWithLabel.svelte";
  import PreviewData from "../../views/previewData.svelte";
  import StatusIcon from "../../views/statusIcon.svelte";
  import KonvaStore from "$src/store/konva";

  export let data: CloudFrontProps;
  export let setLineTargets: (data: TargetFromNodeProps[]) => void;
  export let externalGroup: MetricDataReturnType;
  export let highlights: HighLightProps;
  export let idx: number;

  const datastore = Datastore.getDatastore();
  const konvastore = KonvaStore.getStore();

  const apigateways = externalGroup.find(
  	(et) => et.name === AWS_SERVICES.APIGATEWAYV2
  ) as {
	name: string;
	result: ApiGatewayV2IntegrationProps[];
  } | undefined;

  let offset = 0;
  if (apigateways) {
  	offset = apigateways.result.length;
  }

  let imageData = (data?.Items || []).map((item, i) => {
  	let x = 0,
  		y = 0;
  	if ($datastore.konvaConnectableNodes) {
  		const node = $datastore.konvaConnectableNodes.find(
  			(nd) => nd.id === item.Id
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

  	const lineTargets = item.Origins.Items.map((org) => {
  		return {
  			from: item.Id,
  			to: org.Id,
  		};
  	});
  	const midX = ($konvastore.externalBoundingRect.width / 2) + $konvastore.externalBoundingRect.x - 55;
  	return {
  		Id: item.Id,
  		Name: truncateResourceLabel(item.DomainName),
  		config: {
  			x: midX,
  			y,
  			draggable: true,
  			id: item.Id,
  			label: `Cloudfront CDN ${item.Id}`,
  		} as any,
  		lineTargets,
  		data: item,
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

  const targets = imageData
  	.map((item) => {
  		return item.lineTargets.map((tg) => ({
  			id: item.Id,
  			from: tg.from,
  			to: tg.to,
  			x: item.config.x,
  			y: item.config.y,
  			label: `Cloudfront CDN ${item.Id}`,
  		}));
  	})
  	.flat();

  setLineTargets(targets);

  const dispatch = createEventDispatcher();


  let imageEl: any = null;
  const nodeConfigs = imageData.map((it) => it.config);

  onMount(async () => {
  	const img = document.createElement("img");
  	img.src = "/assets/images/aws/cloudfront.png";
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
  const imageWidth = 80;
  const imageHeight = 80;
</script>

<CdnData
  showModal={state.showModal}
  data={state.data}
  closeModal={async () => {
  	state.showModal = false;
  	await delay(700);
  	state.data = null;
  }}
/>
{#if state.showPreview && state.previewData}
<PreviewData proportions={state.previewProportions} data={state.previewData} color={COLOR_SCHEME.CDN} />
{/if}
<ServiceGroupWithLabel
  label={{
  	text: "Cloudfront CDN",
  	fill: COLOR_SCHEME.CDN,
  }}
  borderColor={COLOR_SCHEME.CDN}
  externalService
  {idx}
>
  {#each imageData as item (item.Id)}
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
		on:mouseout={() => {
			state.showPreview = false;
			state.previewProportions = {
				x: 0,
				y: 0
			};
			state.previewData = null;
		}}
      on:dragmove={() => {
      	dispatch("dragmove", item.config);
      }}
      on:mouseenter={(e) => {
      	dispatch("mouseenter", {
      		id: item.Id,
      		highlights: item.lineTargets,
      	});
      	state.previewData = [ {
      		name: "Http Version",
      		value: item.data.HttpVersion
      	}, {
      		name: "Status",
      		value: item.data.Status || "Unknown"
      	}, {
      		name: "Enabled",
      		value: item.data.Enabled ? "Yes" : "No"
      	} ];
      	state.previewProportions = {
      		x: item.config.x - 220,
      		y: item.config.y
      	};
      	state.showPreview = true;
      }}
      on:mouseleave={(e) => {
      	dispatch("mouseleave", e);
      }}
      on:dragend={() => {
      	dispatch("dragend", item.config);
      }}
      getHandler={(handle) => {
      	const rect = getImageRect({ fill: COLOR_SCHEME.CDN });
      	handle.add(rect);
      }}
    >
	  <StatusIcon status={item.data.Enabled ? "RUNNING" : "STOPPED"} />
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
        	text: item.Name || "Cloudfront",
        	draggable: false,
        	y: -20,
        	x: 0,
        	listening: false,
        	// fill: COLOR_SCHEME.CDN,
        	fontStyle: "bold",
        }}
		on:ready={(e) => {
			console.log("text", e.detail.width());
		}}
      />
    </Group>
  {/each}
</ServiceGroupWithLabel>
