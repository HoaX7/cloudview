<script lang="ts">
  import type { PanelProps } from "$src/customTypes/metricPanels";
  import { getRandomNumber } from "$src/helpers";
  import { onDestroy, onMount } from "svelte";
  import Usage from "../../Usage.svelte";
  import Panel from "./panel.svelte";
  import Icon from "$lib/components/common/Image/index.svelte";

  // export let panels: {
  //   [key: string]: (PanelProps & {
  //     id: string;
  //     pollInterval?: string;
  //     unit?: "seconds" | "hour"
  //     ds?: {
  //       timestamp: string;
  //       value: number;
  //     }[];
  //   })[];
  // };

  const panels: {
    [key: string]: (PanelProps & {
      id: string;
      pollInterval?: string;
      unit?: "seconds" | "hour";
      ds?: {
        timestamp: string;
        value: number;
      }[];
    })[];
  } = {
  	"cpu-usage": [
  		{ id: "i-009e8f028bc970982" },
  		{ id: "i-129e8f028bc970973" },
  	] as any,
  	"ram-usage": [
  		{ id: "i-009e8f028bc970982" },
  		{ id: "i-129e8f028bc970973" },
  	] as any,
  	"disk-usage": [
  		{ id: "i-009e8f028bc970982" },
  		{ id: "i-129e8f028bc970973" },
  	] as any,
  	"bandwidth-usage": [
  		{ id: "i-009e8f028bc970982" },
  		{ id: "i-129e8f028bc970973" },
  	] as any,
  };

  const categories = Object.keys(panels);

  const dt = new Date();
  const mockTs = () =>
  	Array(60)
  		.fill(0)
  		.map(() => {
  			dt.setSeconds(dt.getSeconds() - 1);
  			return {
  				timestamp: dt.toLocaleTimeString(),
  				value: getRandomNumber(60),
  			};
  		})
  		.reverse();

  const mockBandWidth = () =>
  	Array(24)
  		.fill(0)
  		.map(() => {
  			dt.setHours(dt.getHours() - 1);
  			return {
  				timestamp: dt.toLocaleString(),
  				value: getRandomNumber(30),
  			};
  		})
  		.reverse();

  const mockRam = () =>
  	Array(60)
  		.fill(0)
  		.map(() => {
  			dt.setSeconds(dt.getSeconds() - 5);
  			return {
  				timestamp: dt.toLocaleTimeString(),
  				value: getRandomNumber(60),
  			};
  		})
  		.reverse();

  panels["bandwidth-usage"].map((d) => {
  	d.pollInterval = "1h";
  	d.ds = mockBandWidth();
  	d.unit = "hour";
  });
  panels["disk-usage"].map((d) => {
  	d.pollInterval = "1h";
  	d.ds = mockBandWidth();
  	d.unit = "hour";
  });
  panels["cpu-usage"].map((d) => {
  	d.ds = mockTs();
  	d.pollInterval = "1s";
  });
  panels["ram-usage"].map((d) => {
  	d.ds = mockRam();
  	d.pollInterval = "5s";
  });

  let interval: any;
  let ramInterval: any;
  onMount(() => {
  	interval = setInterval(() => {
  		const dt = new Date();
  		const res = panels["cpu-usage"].map((d) => {
  			const timeseries = d.ds || [];
  			timeseries.push({
  				value: getRandomNumber(60),
  				timestamp: dt.toLocaleTimeString(),
  			});
  			timeseries.splice(0, 1);
  			d.ds = timeseries;
  			return d;
  		});
  		panels["cpu-usage"] = res;
  	}, 1000);

  	ramInterval = setInterval(() => {
  		const dt = new Date();
  		const res = panels["ram-usage"].map((d) => {
  			const timeseries = d.ds || [];
  			timeseries.push({
  				value: getRandomNumber(60),
  				timestamp: dt.toLocaleTimeString(),
  			});
  			timeseries.splice(0, 1);
  			d.ds = timeseries;
  			return d;
  		});
  		panels["ram-usage"] = res;
  	}, 5000);
  });

  onDestroy(() => {
  	clearInterval(interval);
  	clearInterval(ramInterval);
  });
</script>

{#each categories as key, i (i)}
  <Panel containerClass="mt-5" title={key} showDD>
    <div class="grid grid-cols-12 gap-4">
      {#each panels[key] as item, i (i)}
        <div class="col-span-12 md:col-span-6">
          <div class="border-2 border-gray-100 shadow-md mt-3">
            <div
              class="flex justify-between hover:cursor-mov hover:bg-gray-100 p-1 px-2"
            >
              <div class="w-4">
                {(item?.ds || [])[(item?.ds || []).length - 1]?.value || 0}%
              </div>
              <div>
                {item.id || "N/A"}
              </div>
              <div class="flex items-center gap-1">
                <Icon
                  src="/assets/images/refresh.svg"
                  width="24"
                  alt="refresh"
                />
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
