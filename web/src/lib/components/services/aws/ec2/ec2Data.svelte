<script lang="ts">
  import type {
  	Ec2VolumeProps,
  	Ec2InstanceProps,
  } from "$src/customTypes/Services";
  import Drawer from "$lib/components/common/Drawer/Drawer.svelte";
  import Typography from "$lib/components/common/Typography/Typography.svelte";
  import { bytesToMegaBytes, toLocaleDate } from "$src/helpers";
  import { getConnectorMappings } from "$src/helpers/konva/index";
  import Datastore from "$src/store/data";
  import RenderConnectionsTable from "../../RenderConnectionsTable.svelte";
  import { getEc2UsageData } from "$src/api/aws";
  import Spinner from "$src/lib/components/common/Loaders/Spinner.svelte";
  import Table from "$src/lib/components/common/Table/Table.svelte";
  import { COLOR_SCHEME } from "$src/colorConfig";

  export let instance: Ec2InstanceProps | null;
  export let volume: Ec2VolumeProps | null;
  export let projectId: string;
  export let serviceId: string;
  export let region: string;

  const datastore = Datastore.getDatastore();

  export let onClose: () => void;
  export let showModal = false;

  let loading = false;

  let connections: { from: string; to: string }[] = [];
  let usage = {
  	bandwidth: 0,
  	cpu: 0
  };
  $: {
  	if (showModal) {
  		connections = getConnectorMappings($datastore, instance?.PublicIpAddress);
  		fetchUsageData();
  	}
  }

  const fetchUsageData = async () => {
  	try {
  		if (!instance) {
  			console.error("No instance selected");
  			return;
  		}
  		loading = true;
  		const resp = await getEc2UsageData({
  			instance: "ec2",
  			instanceId: instance.InstanceId,
  			serviceId,
  			projectId,
  			region,
  		});
  		if (resp.error) throw resp;
  		if (resp.data) {
  			if (resp.data.Bandwidth) {
  				const totalBandwidthUsage = resp.data.Bandwidth.MetricDataResults.reduce((acc, r) => {
  					acc = acc + r.Values.reduce((ac, ar) => ac + ar, 0);
  					return acc;
  				}, 0);
  				usage.bandwidth = totalBandwidthUsage;
  			}
  			if (resp.data.Cpu) {
  				usage.cpu = resp.data.Cpu.MetricDataResults[0].Values[0];
  			}
  		}
  	} catch (err) {
  		console.error("Unable to fetch usage data", err);
  	}
  	loading = false;
  };

  const columns = [
  	{
  		group: [
  			{
  				name: "Region",
  				key: "Placement",
  				subKey: "AvailabilityZone",
  			},
  			{
  				name: "Public IP Address",
  				key: "PublicIpAddress",
  			},
  			{
  				name: "Public DNS Name",
  				key: "PublicDnsName",
  			},
  			{
  				name: "Private IP Address",
  				key: "PrivateIpAddress",
  			},
  			{
  				name: "Private DNS Name",
  				key: "PrivateDnsName",
  			},
  		],
  	},
  	{
  		group: [
  			{
  				name: "vCpus",
  				key: "CpuOptions",
  				subKey: "CoreCount",
  			},
  			{
  				name: "RAM",
  				key: "InstanceType",
  			},
  			{
  				name: "State",
  				key: "State",
  				subKey: "Name",
  			},
  			{
  				name: "Subnet",
  				key: "SubnetId",
  			},
  			{
  				name: "Vpc",
  				key: "VpcId"
  			}
  		],
  	},
  ];

  const renderCol = (col: { name: string; key: string; subKey?: string }) => {
  	const key = col.key as keyof Ec2InstanceProps;
  	if (col.subKey && instance) {
  		return instance[key][col.subKey] || "-";
  	}
  	return (instance || {})[key] || "-";
  };
</script>

<Drawer
  {showModal}
  closeModal={() => {
  	onClose();
  }}
  icon={{
  	src: "/assets/images/aws/ec2.png",
  	alt: "ec2",
  	width: 50,
  	class: "p-2",
  }}
  iconStyle={`background-color:${COLOR_SCHEME.VM};`}
>
  <svelte:fragment slot="header">
    {#if instance}
      <div class="ml-3">
        <Typography variant="h3" weight="semi-bold" font={30} classname="">
          {instance.KeyName} ({instance.InstanceId}) Ec2
        </Typography>
        <div class="flex items-center gap-8 text-sm text-gray-400">
          <Typography variant="div" weight="medium" font={14}>
            {instance.PublicIpAddress}
          </Typography>
          <Typography variant="div" weight="medium" font={14}>
            {instance.Placement.AvailabilityZone}
          </Typography>
          <Typography variant="div" weight="medium" font={14}>
            Created {toLocaleDate(instance.LaunchTime)}
          </Typography>
        </div>
      </div>
    {/if}
  </svelte:fragment>
  {#if instance}
    <Typography variant="h3" weight="medium" font={20}>Usage</Typography>
    <Typography variant="p" weight="medium" font={14} classname="text-gray-400">
      Data is refreshed every 15 minutes.
    </Typography>
    <div class="mt-3 grid grid-cols-2 gap-4">
      <div class="col-span-1 flex items-center">
        Bandwidth: {#if loading}
        <Spinner size="xxs" className="ml-2" />
      {:else}
        {bytesToMegaBytes(usage.bandwidth || 0) || "-"} MB This Month
      {/if}
      </div>
      <div class="col-span-1 flex items-center">
        Cpu: {#if loading}
        <Spinner size="xxs" className="ml-2" />
      {:else}
        {(usage.cpu || 0).toFixed(2) || "-"}%
      {/if}
      </div>
    </div>
    <Typography variant="h3" weight="medium" font={20} classname="mt-5">
      Instance
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

      {#if volume}
        <div class="col-span-12 md:col-span-6 lg:col-span-4">
          <div class="col-span-6 grid grid-cols-5 gap-4 mt-5">
            <div class="col-span-2 text-gray-500">Device Name:</div>
            <div class="col-span-3">
              {volume.Attachments[0].Device}
            </div>
          </div>
          <div class="col-span-6 grid grid-cols-5 gap-4 mt-5">
            <div class="col-span-2 text-gray-500">Size:</div>
            <div class="col-span-3">
              {volume.Size || "-"} GB ({volume.VolumeType}{volume.Encrypted
              	? " Encrypted"
              	: ""})
            </div>
          </div>
        </div>
      {/if}
    </div>
  {/if}
  {#if (instance?.Tags || []).length > 0}
        <Typography classname="mt-3" variant="h3" font={20} weight="medium">
            Tags
        </Typography>
        {#each (instance?.Tags || []) as tag, index (index)}
            <div class="mt-3">
                {tag.Key}: {tag.Value}
            </div>
        {/each}
        {/if}
  {#if (instance?.SecurityGroups || []).length > 0}
  <Typography variant="h3" weight="medium" font={20} classname="mt-5">
    Security Groups
  </Typography>
    <Table columns={[ {
    	name: "Id",
    	key: "GroupId"
    }, {
    	name: "Name",
    	key: "GroupName"
    } ]} data={instance?.SecurityGroups || []} />
  {/if}
  <RenderConnectionsTable {connections} />
</Drawer>
