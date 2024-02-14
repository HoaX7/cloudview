<script lang="ts">
  import Line from "$lib/charts/Line.svelte";
  import type { ChartDataset, ScriptableScaleContext } from "chart.js";

  export let options = { tickSuffix: "%" };
  export let datasets: ChartDataset<"line">[];
  export let labels: string[] = [];
  const getGridColors = (ctx: ScriptableScaleContext, axis: "x" | "y") => {
  	if (ctx.index === 0 || (axis === "x" && ctx.index === 5)) {
  		return "";
  	}
  	return "rgba(102, 102, 102, .2)";
  };
</script>

<Line
  data={{
  	datasets,
  	labels,
  }}
  options={{
  	scales: {
  		y: {
  			grid: { color: (ctx) => getGridColors(ctx, "y") },
  			border: { display: false },
  			beginAtZero: true,
  			ticks: {
  				callback(tickValue, index, ticks) {
  					return tickValue + options.tickSuffix;
  				},
  			},
  			stacked: true,
  			min: 0,
  			max: 10,
  		},
  		x: {
  			grid: { color: (ctx) => getGridColors(ctx, "x") },
  			border: { display: false },
  			beginAtZero: true,
  		},
  	},
  	plugins: {
  		tooltip: {
  			mode: "index",
  			intersect: false,
  		},
  		legend: { position: "bottom" },
  	},
  	parsing: false,
  	normalized: true,
  	animation: false,
  }}
/>
