<script lang="ts">
  import { COLOR_SCHEME } from "$src/colorConfig";
  import type { DynamoDBProps } from "$src/customTypes/services";
  import { bytesToMegaBytes, toLocaleDate } from "$src/helpers";
  import Drawer from "$src/lib/components/common/Drawer/Drawer.svelte";
  import Typography from "$src/lib/components/common/Typography/Typography.svelte";

  export let data: DynamoDBProps[0]["Table"] | null;
  export let closeModal: () => void;
  export let showModal = false;

  const columns = [ {
  	group: [ {
  		name: "Size",
  		value: bytesToMegaBytes(data?.TableSizeBytes || 0) + " MB"
  	}, {
  		name: "Item Count",
  		value: data?.ItemCount || 0
  	}, {
  		name: "Status",
  		value: data?.TableStatus || "Unknown"
  	} ]
  } ];
</script>

<Drawer
  {showModal}
  {closeModal}
  icon={{
  	src: "/assets/images/aws/dynamodb.png",
  	alt: "ec2",
  	width: 50,
  	class: "p-2",
  }}
  iconStyle={`background-color:${COLOR_SCHEME.DB};`}
>
  <svelte:fragment slot="header">
    {#if data}
      <div class="ml-3">
        <Typography variant="h3" weight="semi-bold" font={30} classname="">
          {data.TableName} ({data.TableId})
        </Typography>
        <div class="flex items-center gap-8 text-sm text-gray-400">
          <Typography variant="div" weight="medium" font={14}>
            Created {toLocaleDate(data.CreationDateTime)}
          </Typography>
        </div>
      </div>
    {/if}
  </svelte:fragment>
  {#if data}
    <Typography variant="h3" weight="medium" font={20}>Database</Typography>
    <div class="mt-3 grid grid-cols-2 gap-4">
      <div class="col-span-6 grid grid-cols-5 gap-4 mt-5">
        <div class="col-span-1 text-gray-500">Table Arn:</div>
        <div class="col-span-4">
          {data.TableArn}
        </div>
      </div>
      {#each columns as { group }, index (index)}
        <div class="col-span-12 md:col-span-6 lg:col-span-4">
          {#each group as col, idx (idx)}
            <div class="col-span-6 grid grid-cols-5 gap-4 mt-5">
              <div class="col-span-2 text-gray-500">
                {col.name}:
              </div>
              <div class="col-span-3 break-words">
                {col.value}
              </div>
            </div>
          {/each}
        </div>
      {/each}
    </div>
  {/if}
</Drawer>
