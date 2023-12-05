<script lang="ts">
  import { COLOR_SCHEME } from "$src/colorConfig";
  import type { S3Props } from "$src/customTypes/Services";
  import { bytesToMegaBytes, toLocaleDate } from "$src/helpers";
  import { getConnectorMappings } from "$src/helpers/konva/index";
  import Drawer from "$src/lib/components/common/Drawer/Drawer.svelte";
  import Table from "$src/lib/components/common/Table/Table.svelte";
  import Typography from "$src/lib/components/common/Typography/Typography.svelte";
  import Datastore from "$src/store/data";
  import RenderConnectionsTable from "../../RenderConnectionsTable.svelte";

  export let bucket: S3Props["Data"]["Buckets"][0] | null;
  export let metric: S3Props["Metrics"][0] | null;
  export let acl: S3Props["ACLList"][0]["ACL"] | null;
  export let showModal = false;
  export let onClose: () => void;

  const datastore = Datastore.getDatastore();

  let connections: any[] = [];
  let grantData: any[] = [];
  $: {
  	if (showModal) {
  		connections = getConnectorMappings($datastore, bucket?.Name);
  	}
  	const arr = acl?.Grants.map((gt) => {
  		return {
  			...gt.Grantee,
  			Permission: gt.Permission,
  		};
  	});
  	if (arr) {
  		grantData = arr;
  	}
  }

  // show warning if s3 bucket has public write access
  const isPublicWrite = (acl?.Grants || []).some((grant) => {
  	return (
  		grant.Grantee.URI &&
      grant.Grantee.URI === "http://acs.amazonaws.com/groups/global/AllUsers" &&
      grant.Permission === "WRITE"
  	);
  });

  const columns = [
  	{
  		name: "Type",
  		key: "Type",
  	},
  	{
  		name: "Name",
  		key: "DisplayName",
  	},
  	{
  		name: "Email",
  		key: "EmailAddress",
  	},
  	{
  		name: "URI",
  		key: "URI",
  	},
  	{
  		name: "Permission",
  		key: "Permission",
  	},
  ];
</script>

<Drawer
  {showModal}
  closeModal={onClose}
  icon={{
  	src: "/assets/images/aws/s3.png",
  	alt: "s3",
  	width: 50,
  	class: "p-2",
  }}
  iconStyle={`background-color:${COLOR_SCHEME.OBJECT_STORAGE};`}
>
  <svelte:fragment slot="header">
    {#if bucket}
      <div class="ml-3">
        <Typography variant="h3" weight="semi-bold" font={30}>
          {bucket.Name}
        </Typography>
        <div>
          <Typography
            variant="div"
            weight="medium"
            font={14}
            classname="text-sm text-gray-400"
          >
            Created {toLocaleDate(bucket.CreationDate)}
          </Typography>
        </div>
      </div>
    {/if}
  </svelte:fragment>
  <Typography classname="" variant="h3" weight="semi-bold" font={20}>
    S3 Bucket
  </Typography>
  {#if metric}
    <div class="mt-3">
      Storage Used (approx): {metric.Statistics.Datapoints.length > 0
      	? bytesToMegaBytes((metric.Statistics.Datapoints || [])[0]?.Sum || 0)
      	: "N/A"} MB
    </div>
  {/if}
  {#if acl}
    <Typography classname="mt-5" variant="h3" weight="semi-bold" font={20}>
      Access Control List
    </Typography>
    {#if isPublicWrite}
      <div class="bg-yellow-200 text-sm font-medium p-3 mt-3 rounded">
        <strong>{bucket?.Name || "-"}</strong> S3 Bucket has
        <strong>Public Write</strong> Access. If this was not intended change this
        setting immediately from AWS console.
      </div>
    {/if}
    <div class="mt-3">
      {#if grantData.length > 0}
        <Table data={grantData} {columns} />
      {/if}
    </div>
  {/if}
  <RenderConnectionsTable {connections} />
</Drawer>
