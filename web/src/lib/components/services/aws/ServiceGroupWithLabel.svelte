<!-- @component
This component must be placed inside konva layer or group.

This component is used to group services horizontally and add a label.
 -->
<script lang="ts">
  import Group from "../../common/KonvaCanvas/Group.svelte";
  import Label from "../../common/KonvaCanvas/Label.svelte";
  import Rect from "../../common/KonvaCanvas/Rect.svelte";
  import Tag from "../../common/KonvaCanvas/Tag.svelte";
  import Text from "../../common/KonvaCanvas/Text.svelte";
  import KonvaStore from "$src/store/konva";
  import { getBoundingRect } from "./shapeCache";
  import { getProportions } from "$src/helpers/konva";
  import type { RectConfig } from "konva/lib/shapes/Rect";

  export let label: {
    text: string;
    fill: string;
  };
  export let borderColor: string;
  export let externalService = false;
  export let idx: number; // We are using this property to figure out x - y;

  const konvastore = KonvaStore.getStore();
  let position = {
  	x: 0,
  	y: 0
  };
  if (externalService) {
  	position = getProportions(
  		idx,
  		0,
  		"external",
  		$konvastore.externalBoundingRect?.x,
  		$konvastore.externalBoundingRect?.y
  	);
  } else {
  	position = getProportions(
  		idx,
  		0,
  		"internal",
  		$konvastore.internalBoundingRect?.x,
  		$konvastore.internalBoundingRect?.y
  	);
  }

  let borderConfig = {
  	draggable: false,
  	zIndex: 0,
  	opacity: 0.3,
  	fill: borderColor,
  	x: position.x - 10,
  	y: position.y - 40,
  	width: 0,
  	height: 145,
  	cornerRadius: 5,
  	shadowOffset: {
  		x: 5,
  		y: 5,
  	},
  	shadowOpacity: 0.5,
  	shadowColor: borderColor,
  	shadowEnabled: true,
  	shadowBlur: 10,
  } as RectConfig;
  const rect = getBoundingRect(borderConfig);
  let labelConfig = {
  	x: position.x,
  	y: position.y - 50,
  };
  $: if (!borderConfig.width || borderConfig.width <= 0) {
  	if (externalService) {
  		borderConfig.width = $konvastore.externalBoundingRect?.width || 0;
  	} else {
  		borderConfig.width = $konvastore.internalBoundingRect?.width || 0;
  	}
  }
</script>

<Rect bind:config={borderConfig} handle={rect} />
<Label bind:config={labelConfig}>
  <Tag
    config={{
    	shadowEnabled: true,
    	cornerRadius: 5,
    	fill: label.fill,
    }}
  />
  <Text
    config={{
    	text: label.text,
    	fontSize: 12,
    	fontStyle: "bold",
    	fill: "white",
    	padding: 6,
    }}
  />
</Label>
<Group>
  <slot />
</Group>
