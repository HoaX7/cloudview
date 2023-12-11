<!-- @component

This component must be placed inside of a `InfiniteCanvas` Component.

This is a common component that renders service data based on accessibility.
The accessibility grids show services that can be access publicily (external services) and internal services. -->

<script lang="ts">
  import type { GroupedData } from "$src/customTypes/Services";
  import Konva from "konva";
  import Group from "../../common/KonvaCanvas/Group.svelte";
  import Rect from "../../common/KonvaCanvas/Rect.svelte";
  import Text from "../../common/KonvaCanvas/Text.svelte";
  import KonvaStore from "$src/store/konva";

  export let groupedData: GroupedData;
  export let setInternalServiceBoundingArea: (
    proportions: any
  ) => void = () => {
  	return;
  };

  const konvastore = KonvaStore.getStore();

  const getGroupArray = (groupedData: GroupedData, key: string) =>
  	groupedData[key as keyof GroupedData];

  /**
   * To draw borders we need to fetch the group proportions for
   * each internal and external groups
   */
  let groups: any = {
  	externalGroup: null,
  	internalGroup: null,
  };
  const opacity = 0;
  const colorParts = Konva.Util.getRGB((255 * opacity).toString(16)); // + (255 * settings.opacity).toString(16);
  function getHex(val: number) {
  	return val.toString(16).length < 2
  		? "0" + val.toString(16)
  		: val.toString(16);
  }
  const fillColor =
    "#" +
    getHex(colorParts.r) +
    getHex(colorParts.g) +
    getHex(colorParts.b) +
    getHex(Math.floor(255 * opacity));
  let borderConfigs: any = {
  	externalGroup: {
  		draggable: false,
  		x: 0,
  		y: 0,
  		width: 0,
  		height: 0,
  		visible: true,
  		stroke: "red",
  		strokeWidth: 1,
  		fill: fillColor,
  		dashEnabled: true,
  		dashOffset: 2,
  		dash: [ 5, 5 ],
  		zIndex: 0,
  		cornerRadius: 5,
  		// opacity: .1
  	},
  	internalGroup: {
  		draggable: false,
  		x: 0,
  		y: 0,
  		width: 0,
  		height: 0,
  		visible: true,
  		stroke: "blue",
  		strokeWidth: 1,
  		fill: fillColor,
  		dashEnabled: true,
  		dashOffset: 2,
  		dash: [ 5, 5 ],
  		zIndex: 0,
  		cornerRadius: 5,
  		// opacity: .1
  	},
  };
  let groupingTextConfigs: any = {
  	externalGroup: {
  		x: 0,
  		y: 0,
  		text: "External Services",
  		visible: false,
  		fontStyle: "bold",
  		fontSize: 16,
  	},
  	internalGroup: {
  		x: 0,
  		y: 0,
  		text: "Internal Services",
  		visible: false,
  		fontStyle: "bold",
  		fontSize: 16,
  	},
  };

  type Proportions = {
	x: number;
	y: number;
	width: number;
	height: number;
  }
  export const updateBorder = async (ex_proportions: Proportions, in_proportions: Proportions) => {
  	if (ex_proportions) {
  		borderConfigs.externalGroup.x = ex_proportions.x - 25;
  		borderConfigs.externalGroup.y = ex_proportions.y - 70;
  		borderConfigs.externalGroup.width = ex_proportions.width + 30;
  		borderConfigs.externalGroup.height = ex_proportions.height;
  		borderConfigs.externalGroup.zIndex = 0;
  		groupingTextConfigs.externalGroup.x = ex_proportions.x - 20;
  		groupingTextConfigs.externalGroup.y = ex_proportions.y - 90;
  		groupingTextConfigs.externalGroup.visible = true;
  	}
  	if (in_proportions) {
  		groupingTextConfigs.internalGroup.x = in_proportions.x - 20;
  		groupingTextConfigs.internalGroup.y = in_proportions.y - 90;
  		groupingTextConfigs.internalGroup.visible = true;
  		borderConfigs.internalGroup.x = in_proportions.x - 25;
  		borderConfigs.internalGroup.y = in_proportions.y - 70;
  		borderConfigs.internalGroup.width = in_proportions.width + 30;
  		borderConfigs.internalGroup.height = in_proportions.height;
  		borderConfigs.internalGroup.zIndex = 0;
  		setInternalServiceBoundingArea(in_proportions);
  	}
  };
  $: updateBorder($konvastore.externalBoundingRect, $konvastore.internalBoundingRect);
</script>

{#each Object.keys(groupedData) as key, index (index)}
  <Rect bind:config={borderConfigs[key]} />
  <Text bind:config={groupingTextConfigs[key]} />
  <Group
    getHandler={(handle) => {
    	groups[key] = handle;
    }}
  >
    {#each getGroupArray(groupedData, key) as item, idx (idx + item.name)}
      <slot
        {item}
        {index}
        {idx}
        externalGroup={groupedData.externalGroup}
      />
    {/each}
  </Group>
{/each}
