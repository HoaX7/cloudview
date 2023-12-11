<script lang="ts">
  import type { Vector2d } from "konva/lib/types";
  import Group from "../../common/KonvaCanvas/Group.svelte";
  import Rect from "../../common/KonvaCanvas/Rect.svelte";
  import Text from "../../common/KonvaCanvas/Text.svelte";

  export let color: string;
  export let proportions: Vector2d;
  type Data = {
    name: string;
    value: string;
  };
  export let data: Data[] = [];

  const alignText = (index: number) => {
  	return {
  		x: 10,
  		y: 15 + (index * 20)
  	};
  };
</script>

<Group config={proportions} getHandler={(handle) => {
	handle.moveToTop();
}}>
    <Rect config={{
    	fill: color,
    	shadowColor: color,
    	shadowBlur: 5,
    	shadowOffset: {
    		x: 5,
    		y: 5
    	},
    	width: 200,
    	height: 100,
    	cornerRadius: 5
    }} />
    {#each data.slice(0, 4) as item, index (index)}
    <Text config={{
    	text: `${item.name}: ${item.value || "N/A"}`,
    	fill: "white",
    	fontStyle: "bold",
    	...alignText(index)
    }} />
    {/each}
</Group>
