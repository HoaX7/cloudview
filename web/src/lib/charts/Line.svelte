<script lang="ts">
  import type { DefaultDataPoint } from "chart.js";
  import {
  	Chart as ChartJS,
  	LineController,
  	Title,
  	LineElement,
  	LinearScale,
  	CategoryScale,
  	PointElement,
  	Filler,
  	Tooltip,
  	TimeScale,
  } from "chart.js";
  import Chart from "./Chart.svelte";
  import type { ChartBaseProps } from "$src/customTypes/chart";
  import { onDestroy } from "svelte";

  interface $$Props<TData = DefaultDataPoint<"line">, TLabel = unknown>
    extends Omit<ChartBaseProps<"line", TData, TLabel>, "type"> {
    chart?: ChartJS<"line", TData, TLabel> | null;
  }

  ChartJS.register(
  	Title,
  	LineElement,
  	LineController,
  	LinearScale,
  	CategoryScale,
  	PointElement,
  	Filler,
  	Tooltip,
  	TimeScale,
  );

  $: ChartJS.defaults.elements.point.radius = 0;

  export let chart: $$Props["chart"] = null;
  let props: $$Props;
  let baseChartRef: Chart;

  $: props = $$props as $$Props;

  onDestroy(() => {
  	if (chart) chart.destroy();
  	chart = null;
  });
</script>

<Chart bind:this={baseChartRef} bind:chart type="line" {...props} />
