<script lang="ts">
  import type { PanelProps } from "$src/customTypes/metricPanels";
  import { getRandomNumber } from "$src/helpers";
  import { onDestroy, onMount } from "svelte";
  import Usage from "../../Usage.svelte";
  import Panel from "./panel.svelte";
  import Icon from "$lib/components/common/Image/index.svelte";

  export let panels: {
    [key: string]: (PanelProps & {
      id: string;
      pollInterval?: string;
      unit?: "seconds" | "hour"
      ds?: {
        timestamp: string;
        value: number;
      }[];
    })[];
  };

  const categories = Object.keys(panels);

  const dt = new Date();
  const timeseries = Array(60)
  	.fill(0)
  	.map(() => {
  		dt.setSeconds(dt.getSeconds() - 1);
  		return {
  			timestamp: dt.toLocaleTimeString(),
  			value: getRandomNumber(60),
  		};
  	})
  	.reverse();

    const bts = Array(24).fill(0).map(() => {
    	dt.setHours(dt.getHours() - 1);
    	return {
    		timestamp: dt.toLocaleString(),
    		value: getRandomNumber(30),
    	};
    }).reverse();

  panels["cpu-usage"][0].ds = timeseries;
  panels["cpu-usage"][0].pollInterval = "1s";
  panels["bandwidth-usage"][0].pollInterval = "1h";
  panels["bandwidth-usage"][0].ds = bts;
  panels["bandwidth-usage"][0].unit = "hour";

  let interval: any;
  onMount(() => {
  	interval = setInterval(() => {
  		const dt = new Date();
  		timeseries.push({
  			value: getRandomNumber(60),
  			timestamp: dt.toLocaleTimeString(),
  		});
  		timeseries.splice(0, 1);
  		panels["cpu-usage"][0].ds = timeseries;
  	}, 1000);
  });

  onDestroy(() => {
  	clearInterval(interval);
  });
</script>

{#each categories as key, i (i)}
  <Panel containerClass="mt-5" title={key} showDD>
    <div class="grid grid-cols-12">
      {#each panels[key] as item, i (i)}
        <div class="col-span-12 md:col-span-6">
          <div class="border-2 border-gray-100 shadow-md mt-3">
            <div
              class="flex justify-between hover:cursor-move hover:bg-gray-100 p-1 px-2"
            >
              <div class="w-4">
                {(item?.ds || [])[(item?.ds || []).length - 1]?.value || 0}%
              </div>
              <div>
                {item.id || "N/A"}
              </div>
              <div class="flex items-center gap-1">
                  <Icon src="/assets/images/refresh.svg" width="24" alt="refresh" />
                  {item.pollInterval || "1s"}
              </div>
            </div>
            <div class="mt-3">
              <Usage
                datasets={[
                	{
                		fill: true,
                		backgroundColor: "rgba(86, 204, 242, 0.4)",
                		borderColor: "#2d9cdb",
                		data:
                      item.ds?.map((t) => ({
                      	x: t.timestamp,
                      	y: t.value,
                      })) || [],
                		label: key,
                	},
                ]}
                labels={item.ds?.map((t) => t.timestamp) || []}
                unit={item.unit}
              />
            </div>
          </div>
        </div>
      {/each}
    </div>
  </Panel>
{/each}
