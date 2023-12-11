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
</script>

<Accordion
title="Legend"
buttonClass="flex items-center bg-gray-100 py-2 shadow px-5 hover:bg-gray-200"
disableDefaultBtnClass={true}
reverseDD={true}
showDD={false}
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
