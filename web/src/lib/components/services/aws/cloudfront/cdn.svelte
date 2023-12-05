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
  import { createEventDispatcher, onMount } from "svelte";
  import CdnData from "./cdnData.svelte";
  import { delay } from "$src/helpers";
  import type { MetricDataReturnType } from "$src/customTypes/Services";
  import { AWS_SERVICES } from "$src/helpers/constants";
  import { getProportions } from "$src/helpers/konva/index";
  import Rect from "$src/lib/components/common/KonvaCanvas/Rect.svelte";
  import { COLOR_SCHEME } from "$src/colorConfig";
  import type Konva from "konva";

  export let data: CloudFrontProps;
  export let setLineTargets: (data: TargetFromNodeProps[]) => void;
  export let externalGroup: MetricDataReturnType;
  export let highlights: HighLightProps;

  const datastore = Datastore.getDatastore();

  const apigateways = externalGroup.find(
  	(et) => et.name === AWS_SERVICES.APIGATEWAYV2
  );

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
  	} else {
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
  	return {
  		Id: item.Id,
  		Name: item.DomainName,
  		config: {
  			x,
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
  			(nd) => nd.includes(it.config?.id || "") || nd === it.config?.id
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
  };
  const imageWidth = 80;
  const imageHeight = 80;
  let group: Konva.Group | null = null;
  let borderConfig = {
  	draggable: false,
  	zIndex: 0,
  	opacity: 0,
  	x: 0,
  	y: 0,
  	width: 0,
  	height: 0,
  	cornerRadius: 5,
  };
  const tm = setTimeout(() => {
  	clearTimeout(tm);
  	if (group) {
  		const proportions = group.getClientRect();
  		borderConfig.x = proportions.x - 10;
  		borderConfig.y = proportions.y - 10;
  		borderConfig.width =
        proportions.width + (imageWidth - (imageEl?.width || 0));
  		borderConfig.height =
        proportions.height + 10 + (imageHeight - (imageEl?.height || 0)) / 2;
  	}
  }, 100);
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
<Rect bind:config={borderConfig} />
<Group
  getHandler={(handle) => {
  	group = handle;
  }}
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
      on:dragmove={() => {
      	dispatch("dragmove", item.config);
      }}
      on:mouseenter={(e) => {
      	dispatch("mouseenter", {
      		id: item.Id,
      		highlights: item.lineTargets,
      	});
      }}
      on:mouseleave={(e) => {
      	dispatch("mouseleave", e);
      }}
      on:dragend={() => {
      	dispatch("dragend", item.config);
      }}
    >
      <Rect
        config={{
        	width: imageWidth,
        	height: imageHeight,
        	cornerRadius: 5,
        	fill: COLOR_SCHEME.CDN,
        	x: 0,
        	y: 0,
        }}
      />
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
        	fill: COLOR_SCHEME.CDN,
        	fontStyle: "bold",
        }}
      />
    </Group>
  {/each}
</Group>
