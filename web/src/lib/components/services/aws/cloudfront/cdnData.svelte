<script lang="ts">
  import { COLOR_SCHEME } from "$src/colorConfig";
  import type { CloudFrontItemProps } from "$src/customTypes/aws/cloudfront";
  import { getConnectorMappings } from "$src/helpers/konva/index";
  import Drawer from "$src/lib/components/common/Drawer/Drawer.svelte";
  import Table from "$src/lib/components/common/Table/Table.svelte";
  import Typography from "$src/lib/components/common/Typography/Typography.svelte";
  import Datastore from "$src/store/data";
  import RenderConnectionsTable from "../../RenderConnectionsTable.svelte";

  export let data: CloudFrontItemProps;
  export let showModal = false;
  export let closeModal: () => void;
  const datastore = Datastore.getDatastore();

  let connections: any[] = [];
  $: if (showModal) {
  	connections = getConnectorMappings($datastore, data.Id);
  }

  const groups = [
  	{
  		group: [
  			{
  				name: "Domain Name",
  				key: "DomainName",
  			},
  			{
  				name: "Arn",
  				key: "ARN",
  			},
  			{
  				name: "Enabled",
  				key: "Enabled",
  				boolean: true,
  			},
  		],
  	},
  	{
  		group: [
  			{
  				name: "Status",
  				key: "Status",
  			},
  			{
  				name: "Http Version",
  				key: "HttpVersion"
  			}
  		],
  	},
  ];

  const renderCol = (col: { name: string; key: string; boolean?: boolean; }) => {
  	const key = col.key as keyof CloudFrontItemProps;
  	if (col.boolean) {
  		return (data || {})[key] ? "Yes" : "No";
  	}
  	return (data || {})[key] || "-";
  };

  const originaCols = [ {
  	name: "Domain Name",
  	key: "DomainName"
  }, {
  	name: "Origin Shield Region",
  	keyName: "OriginShield",
  	subKey: "OriginShieldRegion",
  	key: ""
  } ];
</script>

<Drawer
  {closeModal}
  {showModal}
  icon={{
  	src: "/assets/images/aws/cloudfront.png",
  	class: "p-2",
  	alt: "cdn",
  	width: 50,
  }}
  iconStyle={`background-color:${COLOR_SCHEME.CDN};`}
>
  <svelte:fragment slot="header">
    <Typography classname="ml-3" variant="h3" weight="semi-bold" font={30}>
      {data?.Id || "-"}
    </Typography>
  </svelte:fragment>
  {#if data}
  <Typography variant="h3" weight="medium" font={20}>
    Cloudfront
  </Typography>
    <div class="grid grid-cols-12 gap-4">
        {#each groups as { group }, index (index)}
        <div class="col-span-12 md:col-span-6 lg:col-span-4">
            {#each group as col, idx (idx)}
            <div class="col-span-6 grid grid-cols-5 gap-4 mt-5">
              <div class="col-span-2 text-gray-500">
                {col.name}:
              </div>
              <div class="col-span-3 break-words">
                {renderCol(col)}
              </div>
            </div>
          {/each}
        </div>
        {/each}
    </div>
  {/if}
  {#if (data?.Origins?.Items || []).length > 0}
    <Typography classname="mt-5" variant="h3" weight="semi-bold" font={20}>
      Origins
    </Typography>
    <div class="mt-3">
        <Table columns={originaCols} data={data.Origins.Items} />
    </div>
  {/if}
  <RenderConnectionsTable {connections} />
</Drawer>
