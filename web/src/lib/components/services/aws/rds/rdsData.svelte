<script lang="ts">
  import { COLOR_SCHEME } from "$src/colorConfig";
  import type { RDSInstanceProps } from "$src/customTypes/aws/rds";
  import { toLocaleDate } from "$src/helpers";
  import { getConnectorMappings } from "$src/helpers/konva/index";
  import Drawer from "$src/lib/components/common/Drawer/Drawer.svelte";
  import Table from "$src/lib/components/common/Table/Table.svelte";
  import Typography from "$src/lib/components/common/Typography/Typography.svelte";
  import Datastore from "$src/store/data";
  import RenderConnectionsTable from "../../RenderConnectionsTable.svelte";

    export let data: RDSInstanceProps | null;
    export let showModal = false;
    export let closeModal: () => void;

    const datastore = Datastore.getDatastore();
    let connections: any[] = [];
    let subnetData: any[] = [];
    const subnetCols = [ {
    	name: "ID",
    	key: "id"
    }, {
    	name: "Availability Zone",
    	key: "az"
    }, {
    	name: "Status",
    	key: "status"
    } ];
    $: {
    	if (showModal) {
    	    connections = getConnectorMappings($datastore, data?.DBInstanceArn);
    	}
    	if (data?.DBSubnetGroup?.Subnets) {
    		subnetData = data.DBSubnetGroup.Subnets.map((dt) => {
    			return {
    				id: dt.SubnetIdentifier,
    				az: dt.SubnetAvailabilityZone.Name,
    				status: dt.SubnetStatus
    			};
    		});
    	}
    }

    const columns = [
    	{
    		group: [
    			{
    				name: "Database",
    				key: "DBInstanceIdentifier"
    			},
    			{
    				name: "Class",
    				key: "DBInstanceClass"
    			},
    			{
    				name: "Engine",
    				key: "Engine"
    			},
    			{
    				name: "Version",
    				key: "EngineVersion"
    			},
    			{
    				name: "Storage Type",
    				key: "StorageType"
    			}
    		]
    	},
    	{
    		group: [
    			{
    				name: "Storage",
    				key: "AllocatedStorage"
    			},
    			{
    				name: "Maximum Storage Threshold",
    				key: "MaxAllocatedStorage"
    			},
    			{
    				name: "Master Username",
    				key: "MasterUsername"
    			},
    			{
    				name: "Publicly Accessible",
    				key: "PubliclyAccessible",
    				boolean: true
    			},
    			{
    				name: "Vpc",
    				key: "DBSubnetGroup",
    				subKey: "VpcId"
    			}
    		]
    	}
    ];

    const renderCol = (col: { key: string; name: string; boolean?: boolean; subKey?: string; }) => {
    	const key = col.key as keyof RDSInstanceProps;
    	if (!data) return;
    	if (key === "DBSubnetGroup") {
    		return data[key]["VpcId"] || "-";
    	}
    	if (col.boolean) {
    		return data[key] ? "Yes" : "No";
    	}
    	if (key === "AllocatedStorage" || key === "MaxAllocatedStorage") {
    		return (data[key] || "0") + " GB";
    	}
    	return data[key] || "-";
    };
</script>

<Drawer {showModal} {closeModal} icon={{
	src: "/assets/images/aws/rds.png",
	width: 50,
	alt: "rds",
	class: "p-2"
}}
  iconStyle={`background-color:${COLOR_SCHEME.DB};`}
>
<svelte:fragment slot="header">
    {#if data}
    <div class="ml-3">
        <Typography variant="h3" weight="semi-bold" font={30} classname="">
            {data.DBInstanceIdentifier}
          </Typography>
          <Typography
          variant="div"
          weight="medium"
          font={14}
          classname="text-sm text-gray-400">
            Created {toLocaleDate(data.InstanceCreateTime)}
          </Typography>
        </div>
    {/if}
</svelte:fragment>
{#if data}
    <Typography classname="" variant="h3" weight="medium" font={20}>
        RDS Instance
    </Typography>
    <div class="grid grid-cols-12 gap-4">
        {#each columns as { group }, index (index)}
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
      <div class="mt-3">
        End Point: {data.Endpoint.Address}
      </div>
      <div class="mt-3">
        Port: {data.Endpoint.Port}
      </div>
      <div class="mt-3">
        Hosted Zone: {data.Endpoint.HostedZoneId}
      </div>
{/if}
{#if subnetData.length > 0}
<Typography classname="mt-5" variant="h3" weight="medium" font={20}>
    Subnets
</Typography>
<div class="mt-3">
    <Table columns={subnetCols} data={subnetData} />
</div>
{/if}
{#if (data?.VpcSecurityGroups || []).length > 0}
<Typography classname="mt-5" variant="h3" weight="medium" font={20}>
    Vpc Security Groups
</Typography>
<div class="mt-3">
    <Table columns={[ {
    	name: "ID",
    	key: "VpcSecurityGroupId"
    }, {
    	name: "Status",
    	key: "Status"
    } ]} data={data?.VpcSecurityGroups || []} />
</div>
{/if}
<RenderConnectionsTable {connections} />
</Drawer>
