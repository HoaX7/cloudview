<script lang="ts">
  import Datastore from "$src/store/data";
  import KonvaStore from "$src/store/konva";
  import Group from "../../common/KonvaCanvas/Group.svelte";
  import IPhone from "../../common/KonvaCanvas/IPhone.svelte";
  import MacbookIc from "../../common/KonvaCanvas/MacbookIc.svelte";
  import Rect from "../../common/KonvaCanvas/Rect.svelte";
  import { clone } from "$src/helpers";
  import Text from "../../common/KonvaCanvas/Text.svelte";

  const konvastore = KonvaStore.getStore();
  const datastore = Datastore.getDatastore();
  let midY = 0;
  const groupH = 254;
  const groupW = 128;
  const borderWidth = groupW + 40;
  const borderHeight = groupH + 40;
  const position = {
  	x: 0,
  	y: 0,
  };
  $: {
  	midY =
      $konvastore.externalBoundingRect?.height / 2 -
      groupH / 2 +
      ($konvastore.externalBoundingRect?.y || 0) / 2;
  	position.x = $konvastore.externalBoundingRect.x - (groupW + 200); // 200 is spacing.
  	position.y = midY;
  	const res = clone($datastore);
  	const idx = res.konvaTargetFromNodes.findIndex(
  		(nd) => nd.id === "user-connection"
  	);
  	if (idx >= 0) {
  		res.konvaTargetFromNodes[idx].x = position.x + groupW / 2;
  		res.konvaTargetFromNodes[idx].y = midY + groupH / 2 - 10;
  	} else {
  		res.konvaTargetFromNodes.push({
  			x: position.x - 20,
  			y: midY - 10,
  			id: "user-connection",
  			label: "",
  			from: "user-connection",
  			to: "external-services",
  			lineStyle: {
  				dash: [ 5, 5 ],
  				dashOffset: 5,
  				stroke: "gray",
  			},
  		});
  	}
  	const idex = res.konvaConnectableNodes.findIndex(
  		(nd) => nd.id === "external-services"
  	);
  	if (idex >= 0) {
  		res.konvaConnectableNodes[idex].x =
        $konvastore.externalBoundingRect.x - 20;
  		res.konvaConnectableNodes[idex].y = midY + groupH / 2 - 10;
  	} else {
  		res.konvaConnectableNodes.push({
  			id: "external-services",
  			label: "",
  			x: $konvastore.externalBoundingRect.x,
  			y: midY,
  		});
  	}
  	$datastore = res;
  }
</script>

<Rect
  config={{
  	stroke: "black",
  	// dash: [ 5, 5 ],
  	// dashEnabled: true,
  	// dashOffset: 2,
  	cornerRadius: 5,
  	width: borderWidth,
  	height: borderHeight,
  	x: position.x - 20,
  	y: position.y - 10,
  	shadowColor: "black",
  	fill: "white",
  	shadowBlur: 5,
  	shadowOffset: {
  		x: 2,
  		y: 2,
  	},
  }}
/>
<Text
  config={{
  	text: "WWW",
  	x: position.x,
  	y: position.y - 30,
  	fontSize: 14,
  	fontStyle: "bold",
  }}
/>
<Group
  config={{
  	x: position.x,
  	y: position.y,
  }}
>
  <IPhone />
  <MacbookIc />
</Group>
