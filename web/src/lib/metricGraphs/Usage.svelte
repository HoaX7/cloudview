<script lang="ts">
  import type { ChartDataset, ScriptableScaleContext } from "chart.js";
  import Line from "../charts/Line.svelte";

    export let options = { tickSuffix: "%" };
	// 'data' prop in datasets must in Point.
	// i.e {x: "", y: ""}
    export let datasets: any[];
	export let labels: string[];
	export let unit: "hour" | "seconds" = "seconds";
    const getGridColors = (ctx: ScriptableScaleContext, axis: "x" | "y") => {
    	if (ctx.index === 0 || (axis === "x" && ctx.index === 5)) {
    		return "";
    	}
    	return "rgba(102, 102, 102, .2)";
    };
</script>

<Line data={{
	datasets,
	labels,
}} options={{
	scales: {
		y: {
			grid: { color: (ctx) => getGridColors(ctx, "y") },
			border: { display: false },
			// beginAtZero: true,
			ticks: {
				callback(tickValue, index, ticks) {
					return tickValue + options.tickSuffix;
				},
				stepSize: 20
			},
			stacked: true,
			min: 0,
			max: 100,
		},
		x: {
			grid: { color: (ctx) => getGridColors(ctx, "x"), },
			border: { display: false },
			ticks: {
				callback(tickValue, index, ticks) {
					const value = labels[index];
					if (unit === "hour") {
						const [ dd, tt ] = value.split(",");
						const [ hh, mm ] = tt.split(":");
						return `${hh}:${mm}`;
					}
					const [ hh, mm, ss ] = value.split(":");
					return `${mm}:${ss}`;
				},
				source: "labels",
				autoSkipPadding: 10
			},
			
		},
	},
	plugins: {
		tooltip: {
			mode: "index",
			intersect: false,
			callbacks: {}
		},
		legend: { position: "bottom" },
	},
	parsing: false,
	normalized: true,
	animation: false,
}} />
