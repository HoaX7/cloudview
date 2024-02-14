<script lang="ts">
  import Datastore from "$src/store/data";
  import clsx from "clsx";
  import Button from "../common/Button/Button.svelte";
  import { clone } from "$src/helpers";
  import KonvaStore from "$src/store/konva";
  import WindowEvents from "../common/Hooks/WindowEvents.svelte";

  const datastore = Datastore.getDatastore();

  const awsRegions: any = {
  	// "US East": [
  	// 	{
  	// 		name: "N. Virginia",
  	// 		key: "us-east-1",
  	// 	},
  	// 	{
  	// 		name: "Ohio",
  	// 		key: "us-east-2",
  	// 	},
  	// ],
  	"US West": [
  		// {
  		// 	name: "N. California",
  		// 	key: "us-west-1",
  		// },
  		{
  			name: "Oregon",
  			key: "us-west-2",
  		},
  	],
  	// "Asia Pacific": [
  	// 	{
  	// 		name: "Mumbai",
  	// 		key: "ap-south-1",
  	// 	},
  	// 	{
  	// 		name: "Osaka",
  	// 		key: "ap-northeast-3",
  	// 	},
  	// 	{
  	// 		name: "Seoul",
  	// 		key: "ap-northeast-2",
  	// 	},
  	// 	{
  	// 		name: "Singapore",
  	// 		key: "ap-southeast-1",
  	// 	},
  	// 	{
  	// 		name: "Sydney",
  	// 		key: "ap-southeast-2",
  	// 	},
  	// 	{
  	// 		name: "Tokyo",
  	// 		key: "ap-northeast-1",
  	// 	},
  	// ],
  	// Canada: [
  	// 	{
  	// 		name: "Central",
  	// 		key: "ca-central-1",
  	// 	},
  	// ],
  	// Europe: [
  	// 	{
  	// 		name: "Frankfurt",
  	// 		key: "eu-central-1",
  	// 	},
  	// 	{
  	// 		name: "Ireland",
  	// 		key: "eu-west-1",
  	// 	},
  	// 	{
  	// 		name: "London",
  	// 		key: "eu-west-2",
  	// 	},
  	// 	{
  	// 		name: "Paris",
  	// 		key: "eu-west-3",
  	// 	},
  	// 	{
  	// 		name: "Stockholm",
  	// 		key: "eu-north-1",
  	// 	},
  	// ],
  	// "South America": [
  	// 	{
  	// 		name: "São Paulo",
  	// 		key: "sa-east-1",
  	// 	},
  	// ],
  };

  let showDD = false;
</script>

<div
  class="md:relative absolute top-10 left-[50%] md:top-0 md:left-0 -translate-x-1/2 md:translate-0"
>
  <WindowEvents
    callback={() => {
    	if (showDD) showDD = false;
    }}
  />
  <Button
    classname={clsx(
    	"bg-white shadow py-2 px-4 hover:bg-gray-200",
    	showDD ? "rounded-t-3xl rounded-b-none" : "",
    )}
    on:click={(e) => {
    	e.detail.stopPropagation();
    	showDD = !showDD;
    }}
  >
    AWS | {$datastore.selectedRegion}
</Button>
  {#if showDD}
    <div
      class="absolute w-full bg-white shadow z-30 text-sm max-h-60 overflow-auto rounded-b-3xl"
    >
      {#each Object.keys(awsRegions) as region, index (index)}
        <div class="p-2 font-bold">
          {region}
          <hr />
        </div>
        {#each awsRegions[region] as item, idx (idx)}
          <button
            class="p-2 hover:bg-gray-100 w-full text-start"
            on:click={(e) => {
            	e.stopPropagation();
            	if ($datastore.selectedRegion === item.key) {
            		showDD = false;
            		return;
            	}
            	const res = clone($datastore);
            	res.selectedRegion = item.key;
            	res.konvaConnectableNodes = [];
            	res.konvaTargetFromNodes = [];
            	res.fetchData = true;
            	$datastore = res;
            	KonvaStore.clear();
            	showDD = false;
            }}
          >
            {item.name}
            {item.key}
          </button>
        {/each}
      {/each}
    </div>
  {/if}
</div>
