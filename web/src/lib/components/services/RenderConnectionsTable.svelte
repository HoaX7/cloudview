<script lang="ts">
  import Table from "../common/Table/Table.svelte";
  import Typography from "../common/Typography/Typography.svelte";

  export let connections: any[] = [];
</script>

{#if connections.length > 0}
  <Typography variant="h3" weight="medium" font={20} classname="mt-5">
    Connections
  </Typography>
  <Table
    data={connections}
    columns={[
    	{
    		name: "From Service",
    		key: "from",
    	},
    	{
    		name: "To Service",
    		key: "to",
    	},
    ]}
  />
{/if}
