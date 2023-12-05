<script lang="ts">
  import type { MetricDataReturnType, } from "$src/customTypes/Services";
  import Datastore from "$src/store/data";
  import InfiniteCanvas from "../common/KonvaCanvas/Canvas.svelte";
  import { clone, debounce, uniqueArray } from "$src/helpers";
  import FullPageLoader from "../common/Loaders/FullPageLoader.svelte";
  import type {
  	ConnectableNodeProps,
  	HighLightProps,
  	LegendProps,
  	TargetFromNodeProps,
  } from "$src/customTypes/Konva";
  import AwsIndex from "./aws/awsIndex.svelte";
  import type Konva from "konva";
  import { includes } from "lodash";

  // Contains ec2, dynamodb etc objects.
  export let result: MetricDataReturnType;
  export let projectId: string;
  export let serviceId: string;
  let loading = false;

  let canvas: any = null;
  let Stage: Konva.Stage;
  let awsRef: any = null;

  let initialPositions: ConnectableNodeProps[] = [];
  let lineTargets: TargetFromNodeProps[] = [];

  let legend: LegendProps[] = [];

  $: debounce((data, targets) => {
  	console.log("Debounced initial positions & line targets after 200ms");
  	// handleInitPosition(data);
  	// setTargetFromNodes(targets);
  	setInitPositionAndTargetNodes(data, targets);
  	awsRef?.updateBorder();
  }, 200)(initialPositions, lineTargets);

  const setInitPositionAndTargetNodes = (data: ConnectableNodeProps[], targets: TargetFromNodeProps[]) => {
  	if ($datastore) {
  		const res = clone($datastore);
  		data.map((item) => {
  			const idx = (res.konvaConnectableNodes || []).findIndex(
  				(nd) => nd.id === item.id
  			);
  			if (idx >= 0) {
  				res.konvaConnectableNodes[idx] = item;
  			} else {
  				res.konvaConnectableNodes = [
  					...(res.konvaConnectableNodes || []),
  					item,
  				];
  			}
  		});
  		res.konvaTargetFromNodes = uniqueArray([
  			...(res.konvaTargetFromNodes || []),
  			...targets,
  		]);
  		$datastore = res;
  	}
  };

  // @deprecated - Use 'setInitPositionAndTargetNodes'
  const setTargetFromNodes = (targets: TargetFromNodeProps[]) => {
  	if ($datastore) {
  		$datastore.konvaTargetFromNodes = uniqueArray([
  			...($datastore.konvaTargetFromNodes || []),
  			...targets,
  		]);
  	}
  };

  // on:click={() => stage?.handleRepositionStage(imageConfig.x, imageConfig.y)}

  const datastore = Datastore.getDatastore();
  let _region = $datastore.selectedRegion;

  $: {
  		if ($datastore.fetchData) {
  		_region = $datastore.selectedRegion;
  			// fetch data depending on provider
  		if (awsRef) {
  				awsRef.refetchData(_region);
  			}
  		initialPositions = [];
  		lineTargets = [];
  		const res = clone($datastore);
  		res.konvaConnectableNodes = [];
  		res.konvaTargetFromNodes = [];
  		res.fetchData = false;
  		$datastore = res;
  	}
  }

  // Need to fetch data according to provider (Aws, gcp ...etc)
  

  // @deprecated - Use 'setInitPositionAndTargetNodes'
  const handleInitPosition = (detail: any[]) => {
  	// This allows us to pin-point where the
  	// node could be, to draw lines
  	if ($datastore) {
  		const res = clone($datastore);
  		detail.map((item) => {
  			const idx = (res.konvaConnectableNodes || []).findIndex(
  				(nd) => nd.id === item.id
  			);
  			if (idx >= 0) {
  				res.konvaConnectableNodes[idx] = item;
  			} else {
  				res.konvaConnectableNodes = [
  					...(res.konvaConnectableNodes || []),
  					item,
  				];
  			}
  		});
  		$datastore = res;
  	}
  };

  // update the datastore positions when the instances are moved around
  const handleDragInstance = (detail: any) => {
  	if ($datastore) {
  		const res = clone($datastore);
  		res.dragNodeId = detail.id;
  		res.konvaConnectableNodes = (
  			res.konvaConnectableNodes || []
  		).map((nd) => {
  			if (nd.id === detail.id) {
  				nd.x = detail.x;
  				nd.y = detail.y;
  			}
  			return nd;
  		});

  		res.konvaTargetFromNodes = (
  			res.konvaTargetFromNodes || []
  		).map((tg) => {
  			if (tg.id === detail.id) {
  				tg.x = detail.x;
  				tg.y = detail.y;
  			}
  			return tg;
  		});

  		$datastore = res;
  	}
  	// stage?.updateConnector();
  };

  let highlights: HighLightProps = {} as HighLightProps;
  type D = {
	id: string;
	highlights: {
		from: string;
		to: string;
	}[];
	extras: string[];
  }
  const setHighlights = (detail: D) => {
  	const data = {
  		nodes: [ ...new Set([
  			detail.id, ...(detail?.highlights || []).map((hl) => hl.to), ...(detail.extras || [])
  		]) ],
  		lines: (detail?.highlights || []).map((hl) => `line-${hl.from}-${hl.to}`)
  	};
  	// We are setting a dummy line to disable all lines when no nodes are connected.
  	if (data.nodes.length > 0 && data.lines.length <= 0) {
  		data.lines = [ "dummy-line" ];
  	}
  	highlights = data;
  };

  const highlightNodes = (data: string[]) => {
  	let lines: string[] = [];
  	let nodes: string[] = clone(data);
  	data.forEach((nd) => {
  		const target = $datastore.konvaTargetFromNodes.find((tg) => tg.to.includes(nd));
  		if (target) {
  			lines.push(`line-${target.from}-${target.to}`);
  			nodes.push(target.id);
  		}
  	});
  	if (nodes.length > 0 && lines.length <= 0) {
  		lines = [ "dummy-line" ];
  	}
  	 highlights = {
  		nodes,
  		lines
  	 };
  };
</script>

{#if loading}
  <FullPageLoader />
{/if}
<!-- TODO: Display Services based on Provider.  -->
<!-- Currently only AWS provider is considered -->
<InfiniteCanvas
  bind:this={canvas}
  on:init={(e) => {
  	Stage = e.detail;
  }}
  {legend}
  {highlights}
  on:highlight-nodes={(e) => {
  	highlightNodes(e.detail);
  }}
>
<AwsIndex
  	bind:this={awsRef}
    {projectId}
    {serviceId}
    region={_region}
    {result}
	{highlights}
	setLoading={(bool) => loading = bool}
    on:mouseenter={(e) => {
    	if (e.detail.id) {
    		setHighlights(e.detail);
    	}
    	Stage.container().style.cursor = "pointer";
    }}
    on:mouseleave={() => {
    	Stage.container().style.cursor = "default";
    	highlights = {
    		nodes: [],
    		lines: [] 
    	};
    	canvas?.resetLineHighlights();
    }}
    on:click={(e) => {
    	const { x, y, width, height } = e.detail;
    	canvas?.handleRepositionStage(x, y, width, height);
    }}
    setLineTargets={(dataTargets) => {
    	/**
       * It is recommended to debounce this data before
       * updating it to store to avoid too many update calls.
       */
    	lineTargets.push(...dataTargets);
    }}
    on:initialPosition={(e) => {
    	/**
       * It is recommended to debounce this data before
       * updating it to store to avoid too many update calls.
       */
    	initialPositions.push(...e.detail);
    }}
    on:dragmove={(e) => {
    	// has to be updated real-time for smoother experience
    	handleDragInstance(e.detail);
    }}
	on:dragend={() => {
		$datastore.dragNodeId = null;
	}}
	setLegend={(data) => {
		data.forEach((item) => {
			const idx = legend.findIndex((l) => l.id === item.id);
			if (idx >= 0) {
				legend[idx].highlight = [
					...new Set([ ...legend[idx].highlight, ...item.highlight ])
				];
				legend[idx].count = legend[idx].highlight.length;
				legend[idx].colors = [
					...new Set([ ...legend[idx].colors, ...item.colors ]),
				];
			} else {
				legend.push({
					...item,
					count: 1,
				});
			}
		});
	}}
  />
</InfiniteCanvas>
