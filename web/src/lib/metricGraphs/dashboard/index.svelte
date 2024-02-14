<script lang="ts">
  import type { ResourceDataReturnType } from "$src/customTypes/services";
  import Panels from "./panels/index.svelte";
  import type {
  	MetricPanelProps,
  	PanelProps,
  } from "$src/customTypes/metricPanels";
  import CreatePanel from "./panels/createPanel.svelte";
  import clsx from "clsx";

  export const resourceList: ResourceDataReturnType = [];
  export let data: MetricPanelProps[];

  const panels = data?.reduce((acc, r) => {
  	r.panels?.forEach((p) => {
  		acc[p.type] = (acc[p.type] || []).concat({
  			...p,
  			id: r.instanceId,
  		});
  	});
  	return acc;
  }, {} as { [key: string]: (PanelProps & { id: string })[] });

  let state = { showModal: false };
</script>

<div class="mt-5">
  <button
    class={clsx("rounded-full py-2 px-3 bg-gradient text-white text-sm font-bold",
    	"disabled:opacity-75 disabled:cursor-not-allowed")}
    on:click={() => {
    	state.showModal = !state.showModal;
    }}
    disabled
  >
    {state.showModal ? "Back" : "Add Panel"}
  </button>
</div>

{#if state.showModal}
  <CreatePanel />
{:else}
<div class="pb-5">
    <Panels {panels} />
</div>
{/if}
