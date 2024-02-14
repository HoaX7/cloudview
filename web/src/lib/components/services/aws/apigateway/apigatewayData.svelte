<script lang="ts">
  import { COLOR_SCHEME } from "$src/colorConfig";
  import type { ApiGatewayWithIntegrationProps } from "$src/customTypes/services";
  import { getConnectorMappings } from "$src/helpers/konva/index";
  import Drawer from "$src/lib/components/common/Drawer/Drawer.svelte";
  import Table from "$src/lib/components/common/Table/Table.svelte";
  import Typography from "$src/lib/components/common/Typography/Typography.svelte";
  import Datastore from "$src/store/data";
  import RenderConnectionsTable from "../../RenderConnectionsTable.svelte";

  export let data: ApiGatewayWithIntegrationProps | null;
  const datastore = Datastore.getDatastore();

  export let showModal = false;
  export let closeModal: () => void;

    let connections: { from: string; to: string; }[] = [];
    $: if (showModal) {
    	connections = getConnectorMappings($datastore, data?.ApiId);
    }
</script>

<Drawer
  {showModal}
  {closeModal}
  icon={{
  	src: "/assets/images/aws/api-gateway.png",
  	alt: "api-gateway",
  	width: 50,
  	class: "p-2",
  }}
  iconStyle={`background-color:${COLOR_SCHEME.GATEWAY};`}
>
  <svelte:fragment slot="header">
    {#if data}
      <Typography variant="h3" weight="semi-bold" font={30} classname="ml-3">
        {data.Name} ({data.ApiId})
      </Typography>
    {/if}
  </svelte:fragment>
    {#if data}
      <Typography variant="h3" weight="medium" font={20}>
        Integrations
      </Typography>
      {#each data.integrations as integration (integration.IntegrationId)}
        <Typography variant="div" weight="regular" font={16} classname="mt-3">
          {integration.IntegrationMethod}
          {integration.IntegrationUri}
        </Typography>
      {/each}
    {/if}
    <RenderConnectionsTable {connections} />
</Drawer>
