<script lang="ts">
  import type {
  	ELBV2Props,
  	MetricDataReturnType,
  } from "$src/customTypes/Services";
  import Group from "$src/lib/components/common/KonvaCanvas/Group.svelte";
  import Image from "$src/lib/components/common/KonvaCanvas/Image.svelte";
  import Text from "$src/lib/components/common/KonvaCanvas/Text.svelte";
  import Datastore from "$src/store/data";
  import { createEventDispatcher, onMount } from "svelte";
  import ElbV2Data from "./elbV2Data.svelte";
  import { delay } from "$src/helpers";
  import { AWS_SERVICES, LEGEND_NAMES } from "$src/helpers/constants";
  import { getProportions } from "$src/helpers/konva/index";
  import type { HighLightProps, LegendProps } from "$src/customTypes/Konva";
  import { COLOR_SCHEME } from "$src/colorConfig";
  import Rect from "$src/lib/components/common/KonvaCanvas/Rect.svelte";
  import type Konva from "konva";
  import type { GroupConfig } from "konva/lib/Group";

  export let data: ELBV2Props;
  export let externalGroup: MetricDataReturnType;
  export let setLegend: (legend: LegendProps[]) => void;
  export let highlights: HighLightProps;
  const datastore = Datastore.getDatastore();

  let legend: LegendProps[] = [];

  const apigateways = externalGroup.find(
  	(et) => et.name === AWS_SERVICES.APIGATEWAYV2
  );
  const cdns = externalGroup.find((et) => et.name === AWS_SERVICES.CLOUDFRONT);
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
  	} else {
  		const proportions = getProportions(offset, i, "external");
  		x = proportions.x;
  		y = proportions.y;
  	}
  	return {
  		config: {
  			x,
  			y,
  			id: lb.DNSName,
  			draggable: true,
  			label: `Load Balancer ${lb.CanonicalHostedZoneId}`,
  		} as GroupConfig,
  		id: lb.DNSName,
  		name: `${lb.LoadBalancerName} ${lb.CanonicalHostedZoneId}`,
  		data: lb,
  	};
  });
  $: {
  	imageData = imageData.map((it) => {
  		const node = (highlights?.nodes || []).find((nd) => nd.includes(it.config?.id || "") || nd === it.config?.id);
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

<ElbV2Data
  showModal={state.showModal}
  closeModal={async () => {
  	state.showModal = false;
  	await delay(700);
  	state.data = null;
  }}
  data={state.data}
/>
<Rect bind:config={borderConfig} />
<Group
  getHandler={(handle) => {
  	group = handle;
  }}
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
      }}
      on:mouseleave={(e) => {
      	dispatch("mouseleave", e);
      }}
    >
      <Rect
        config={{
        	width: imageWidth,
        	height: imageHeight,
        	cornerRadius: 5,
        	fill: COLOR_SCHEME.LOADBALANCER,
        	x: 0,
        	y: 0,
        }}
      />
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
        	fill: COLOR_SCHEME.LOADBALANCER,
        	fontStyle: "bold",
        }}
      />
    </Group>
  {/each}
</Group>
