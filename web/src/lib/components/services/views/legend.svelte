<script lang="ts">
  import type { LegendProps } from "$src/customTypes/Konva";
  import { groupByKey } from "$src/helpers";
  import clsx from "clsx";
  import Accordion from "../../common/Accordion/Accordion.svelte";
  import Typography from "../../common/Typography/Typography.svelte";
  import { createEventDispatcher } from "svelte";

  export let legend: LegendProps[];

  const legendGroup = groupByKey(legend, "name");

  const dispatch = createEventDispatcher();

  //   export let internalServiceBoundingRect: Vector2d & {
  //     width: number;
  //     height: number;
  //   } = {
  //   	x: 0,
  //   	y: 0,
  //   	width: 0,
  //   	height: 0,
  //   };
  //   let group: Konva.Group | null = null;
  //   const padding = 100;
  //   const rectXpadding = 20;
  //   const rectYpadding = 10;

  //   let config = {
  //   	x: 0,
  //   	y: 0,
  //   	draggable: true,
  //   	visible: false,
  //   };

  //   let rectConfig = {
  //   	fill: "#f3f4f6",
  //   	width: 0,
  //   	height: 0,
  //   	x: 0,
  //   	y: 0,
  //   };

  //   const timeout = setTimeout(() => {
  //   	config.x =
  //       internalServiceBoundingRect.x +
  //       internalServiceBoundingRect.width +
  //       padding;
  //   	config.y = internalServiceBoundingRect.y;
  //   	config.visible = true;
  //   	const proportions = group?.getClientRect();
  //   	clearTimeout(timeout);
  //   	if (proportions) {
  //   		rectConfig.width = proportions.width + rectXpadding;
  //   		rectConfig.height = proportions.height + rectXpadding;
  //   		rectConfig.x = config.x - rectYpadding;
  //   		rectConfig.y = config.y - rectYpadding;
  //   	}
  //   }, 200);

  //   const getPreviousKey = (item: Record<string, any>, index: number) => {
  //   	const keys = Object.keys(item);
  //   	if (index <= 0) return keys[0];
  //   	return keys[index - 1];
  //   };
</script>

<!-- <Rect bind:config={rectConfig} />
<Group
  bind:config
  getHandler={(handle) => {
  	group = handle;
  }}
  on:dragmove={() => {
  	rectConfig.x = config.x - rectYpadding;
  	rectConfig.y = config.y - rectYpadding;
  }}
>
  <Text config={{ text: "Legend" }} /> -->
<!-- {#each Object.keys(legendGroup) as key, index (index)}
    <Text
      config={{
      	text: key,
      	y:
          index * 50 +
          (0 * legendGroup[getPreviousKey(legendGroup, index)].length),
      	fontStyle: "bold",
      	fontSize: 13,
      }}
    />
    {#each legendGroup[key] as item, idx (idx)}
      <Group config={{ y: (index * 50) + 25 + (idx * 30), }}>
        {#each item.colors as color, cidex (cidex)}
          <Circle
            config={{
            	radius: 50,
            	fill: color,
            	width: 10,
            	height: 10,
            	x: 5,
            }}
          />
        {/each}
        <Text
          config={{
          	text: item.id,
          	fontSize: 12,
          	// x: 15 * item.colors.length,
          }}
        />
      </Group>
    {/each}
  {/each} -->
<!-- </Group> -->

<div>
  <Accordion
    title="Legend"
    buttonClass="flex items-center bg-gray-100 py-2 shadow px-5 hover:bg-gray-200"
    disableDefaultBtnClass={true}
    reverseDD={true}
  >
    <div class="shadow bg-gray-100 p-3 mb-2">
      {#each Object.keys(legendGroup) as key, index (index)}
        <div class="mb-4">
          <Typography classname="" variant="div" weight="medium" font={14}>
            {key}
          </Typography>
          {#each legendGroup[key] as item, idx (idx)}
            <div class="flex items-center mt-3 text-sm">
              {#each item.colors as color, cidx (cidx)}
                <span
                  class={clsx("mr-1 w-3 h-3 rounded-full inline-block")}
                  style={`background-color:${color};`}
                />
              {/each}
              <span
                class="hover:cursor-pointer hover:underline"
                on:mouseover={() => {
                	dispatch("highlight", item.highlight);
                }}
                on:mouseleave={() => {
                	dispatch("reset-highlight");
                }}
                role="contentinfo"
                on:focus={() => {
                	return;
                }}
              >
                {item.id}
              </span>
            </div>
          {/each}
        </div>
      {/each}
    </div>
  </Accordion>
</div>
