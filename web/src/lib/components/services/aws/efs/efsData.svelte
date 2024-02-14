<script lang="ts">
  import { COLOR_SCHEME } from "$src/colorConfig";
  import type { EFSProps } from "$src/customTypes/services";
  import { bytesToMegaBytes, toLocaleDate } from "$src/helpers";
  import { getConnectorMappings } from "$src/helpers/konva/index";
  import Drawer from "$src/lib/components/common/Drawer/Drawer.svelte";
  import Typography from "$src/lib/components/common/Typography/Typography.svelte";
  import Datastore from "$src/store/data";
  import RenderConnectionsTable from "../../RenderConnectionsTable.svelte";

  export let data: EFSProps["FileSystems"][0] | null;
  export let showModal = false;
  export let closeModal: () => void;

  const datastore = Datastore.getDatastore();
  let connections: any[] = [];
  $: if (showModal) {
  	connections = getConnectorMappings($datastore, data?.FileSystemId);
  }
</script>

<Drawer
  {showModal}
  {closeModal}
  icon={{
  	alt: "efs",
  	src: "/assets/images/aws/efs.png",
  	width: 50,
  	class: "p-2",
  }}
  iconStyle={`background-color:${COLOR_SCHEME.FILE_SYSTEMS};`}
>
  <svelte:fragment slot="header">
    {#if data}
      <div class="ml-3">
        <Typography variant="h3" weight="semi-bold" font={30}>
          {data.Name || "-"} ({data.FileSystemId})
        </Typography>
        <div>
          <Typography
            variant="div"
            weight="medium"
            font={14}
            classname="text-sm text-gray-400"
          >
            Created {toLocaleDate(data.CreationTime)}
          </Typography>
        </div>
      </div>
    {/if}
  </svelte:fragment>
  {#if data}
    <Typography variant="h3" weight="medium" font={20}>File Storage</Typography>
    <div class="mt-3">
      Name: {data.Name}
    </div>
    <div class="mt-3">
      Size: {bytesToMegaBytes(data?.SizeInBytes?.Value || 0)} MB
    </div>
    <div class="mt-3">
      Creation Token: {data.CreationToken}
    </div>
    <div class="mt-3">
      File System Arn: {data.FileSystemArn || "-"}
    </div>
    <div class="mt-3">
      Performance Mode: {data.PerformanceMode || "-"}
    </div>
    {#if data.Tags.length > 0}
      <Typography classname="mt-3" variant="h3" font={20} weight="medium">
        Tags
      </Typography>
      {#each data.Tags as tag, index (index)}
        <div class="mt-3">
          {tag.Key}: {tag.Value}
        </div>
      {/each}
    {/if}
  {/if}
  <RenderConnectionsTable {connections} />
</Drawer>
