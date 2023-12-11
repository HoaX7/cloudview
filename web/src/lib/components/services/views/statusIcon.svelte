<script lang="ts">
  import { STATUS_COLORS } from "$src/helpers/constants";
  import Circle from "../../common/KonvaCanvas/Circle.svelte";

  export let status: "RUNNING" | "STOPPED" | "UNKNOWN" = "UNKNOWN";
</script>

<Circle
        config={{
        	fill: STATUS_COLORS[status],
        	width: 10,
        	height: 10,
        	x: 10,
        	y: 10,
        }}
      />