<script lang="ts">
  import { createEventDispatcher, onMount } from "svelte";
  import Image from "../../../common/KonvaCanvas/Image.svelte";
  import type { ApiGatewayWithIntegrationProps } from "$src/customTypes/ervices;
  import Group from "../../../common/KonvaCanvas/Group.svelte";
  import Text from "../../../common/KonvaCanvas/Text.svelte";
  import type { HighLightProps, TargetFromNodeProps } from "$src/customTypes/konva";
  import Datastore from "$src/store/data";
  import ApigatewayData from "./apigatewayData.svelte";
  import { delay } from "$src/helpers";
  import { getProportions, truncateResourceLabel } from "$src/helpers/konva/index";
  import Rect from "$src/lib/components/common/KonvaCanvas/Rect.svelte";
  import { COLOR_SCHEME } from "$src/colorConfig";
  import type Konva from "konva";
  import { getImageRect } from "../shapeCache";
  import ServiceGroupWithLabel from "../ServiceGroupWithLabel.svelte";
  import KonvaStore from "$src/store/konva";

  export let data: ApiGatewayWithIntegrationProps[];
  export let idx: number;
  const datastore = Datastore.getDatastore();
  const konvastore = KonvaStore.getStore();

  /**
   * Line Targets - Indicates the number of instances
   * to point to from Api-Gateway.
   *
   * This can use re-used in other components to
   * show a arrow-line from the instance image to another node.
   */
  export let setLineTargets: (data: TargetFromNodeProps[]) => void;
  export let highlights: HighLightProps;

  export const projectId: string = "";
  export const providerAccountId: string = "";
  export const region: string = "";

  /**
   * IMPORTANT - since External data is re-ordered we know the order of each
   * service. api gateway -> cloudfront -> elb -> route53.
   * We can now determine `offset` padding depending on how many apigateway ...etc
   * services are connected.
   */
  let offset = 0;

  let imageEl: any = null;
  let apiGateways = data.map((item, i) => {
  	// Save node positions instead of random position each time.
  	let x = 0,
  		y = 0;
  	if ($datastore.konvaConnectableNodes) {
  		const node = $datastore.konvaConnectableNodes.find(
  			(nd) => nd.id === item.ApiId
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
  	// 55 is image padding with width
  	const midX = ($konvastore.externalBoundingRect.width / 2) + $konvastore.externalBoundingRect.x - 55;
  	return {
  		ApiId: item.ApiId,
  		Name: truncateResourceLabel(item.Name),
  		config: {
  			x: midX,
  			y,
  			draggable: false,
  			id: item.ApiId,
  			label: `Api Gateway (${item.ApiId})`,
  		} as any,
  		lineTargets: item.lineTargets,
  		data: item,
  	};
  });
  $: {
  	apiGateways = apiGateways.map((it) => {
  		const node = (highlights?.nodes || []).find((nd) => nd?.includes(it.config?.id || "") || nd === it.config?.id);
  		if (highlights.nodes && highlights.nodes.length > 0 && !node) {
  			it.config.opacity = .3;
  			return it;
  		}
  		it.config.opacity = 1;
  		return it;
  	});
  }
  const targets = apiGateways
  	.map((d) =>
  		d.lineTargets.map((tg) => ({
  			id: d.ApiId,
  			from: tg.from,
  			to: tg.to,
  			x: d.config.x,
  			y: d.config.y,
  			label: `Api Gateway (${d.ApiId})`,
  		}))
  	)
  	.flat();
  setLineTargets(targets);

  const dispatch = createEventDispatcher();

  const nodeConfigs = apiGateways.map((it) => it.config);

  onMount(async () => {
  	const img = document.createElement("img");
  	img.src = "/assets/images/aws/api-gateway.png";
  	img.onload = () => {
  		imageEl = img;
  	};

  	dispatch("initialPosition", nodeConfigs);
  });

  const state: any = {
  	showModal: false,
  	data: null,
  };

  const imageWidth = 80;
  const imageHeight = 80;
</script>

<ApigatewayData
  data={state.data}
  showModal={state.showModal}
  closeModal={async () => {
  	state.showModal = false;
  	await delay(700);
  	state.data = null;
  }}
/>
<ServiceGroupWithLabel
  label={{
  	text: "Api Gateways",
  	fill: COLOR_SCHEME.GATEWAY
  }}
  borderColor={COLOR_SCHEME.GATEWAY}
  externalService
  {idx}
>
  {#each apiGateways as item (item.ApiId)}
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
      on:mouseenter={() => {
      	dispatch("mouseenter", {
      		id: item.ApiId,
      		highlights: item.lineTargets,
      	});
      }}
      on:mouseleave={(e) => {
      	dispatch("mouseleave", e);
      }}
      on:dragend={() => {
      	dispatch("dragend", item.config);
      }}
      on:dragmove={() => {
      	dispatch("dragmove", item.config);
      }}
	  getHandler={(handle) => {
	  	const rect = getImageRect({ fill: COLOR_SCHEME.GATEWAY });
	  	handle.add(rect);
	  }}
    >
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
        	text: item.Name,
        	draggable: false,
        	y: -20,
        	x: 0,
        	listening: false,
        	// fill: COLOR_SCHEME.GATEWAY,
        	fontStyle: "bold",
        	// fontFamily: "inherit"
        }}
      />
    </Group>
  {/each}
</ServiceGroupWithLabel>
