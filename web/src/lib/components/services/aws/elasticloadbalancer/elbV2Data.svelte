<script lang="ts">
  import { COLOR_SCHEME } from "$src/colorConfig";
  import type { ELBV2Props } from "$src/customTypes/ervices;
  import { toLocaleDate } from "$src/helpers";
  import { getConnectorMappings } from "$src/helpers/konva/index";
  import Drawer from "$src/lib/components/common/Drawer/Drawer.svelte";
  import Table from "$src/lib/components/common/Table/Table.svelte";
  import Typography from "$src/lib/components/common/Typography/Typography.svelte";
  import Datastore from "$src/store/data";
  import RenderConnectionsTable from "../../RenderConnectionsTable.svelte";

  export let data: ELBV2Props["LoadBalancers"][0] | null;
  export let showModal = false;
  export let closeModal: () => void;
  const datastore = Datastore.getDatastore();

  let connections: any[] = [];
  $: if (showModal) {
  	connections = getConnectorMappings($datastore, data?.DNSName);
  }

  const columns = [ {
  	name: "Subnet",
  	key: "SubnetId"
  }, {
  	name: "Zone",
  	key: "ZoneName"
  } ];
</script>

<Drawer
  {showModal}
  {closeModal}
  icon={{
  	src: "/assets/images/aws/elbv2.png",
  	width: 50,
  	alt: "loadbalancer",
  	class: "p-2",
  }}
  iconStyle={`background-color:${COLOR_SCHEME.LOADBALANCER};`}
>
  <svelte:fragment slot="header">
    {#if data}
    <div class="ml-3">
        <Typography variant="h3" weight="semi-bold" font={30}>
          {data.LoadBalancerName}
          {data.CanonicalHostedZoneId}
        </Typography>
        <Typography variant="div" weight="medium" font={14} classname="text-sm text-gray-400">
          Created {toLocaleDate(data.CreatedTime)}
        </Typography>
      </div>
    {/if}
  </svelte:fragment>
  {#if data}
    <Typography weight="medium" font={20} variant="h3">
        Load Balancer
    </Typography>
    <div class="mt-3">
        DNS Name: {data.DNSName}
    </div>
    <div class="mt-3">
        ARN: {data.LoadBalancerArn}
    </div>
    <div class="mt-3">
        Scheme: {data.Scheme}
    </div>
    <div class="mt-3">
        Security Groups: {data.SecurityGroups.join(", ") || "-"}
    </div>
    <div class="mt-3">
        Vpc: {data.VpcId || "-"}
    </div>
    <div class="mt-3">
        State: {data.State.Code}
    </div>
    <div class="mt-3">
        Reason: {data.State.Reason || "-"}
    </div>
    <Typography weight="medium" font={20} variant="h3" classname="mt-5">
        Availability Zones
    </Typography>
    <Table columns={columns} data={data.AvailabilityZones || []} />
  {/if}
  <RenderConnectionsTable {connections} />
</Drawer>
