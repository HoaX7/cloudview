<script lang="ts">
  import { COLOR_SCHEME } from "$src/colorConfig";
  import type { LambdaFunctionProps } from "$src/customTypes/ervices;
  import { toLocaleDate } from "$src/helpers";
  import { getConnectorMappings } from "$src/helpers/konva/index";
  import Drawer from "$src/lib/components/common/Drawer/Drawer.svelte";
  import Typography from "$src/lib/components/common/Typography/Typography.svelte";
  import Datastore from "$src/store/data";
  import RenderConnectionsTable from "../../RenderConnectionsTable.svelte";

  export let data: LambdaFunctionProps["Functions"][0] | null;
  export let showModal = false;
  export let closeModal: () => void;

  const datastore = Datastore.getDatastore();

  let connections: any[];
  $: if (showModal) {
  	connections = getConnectorMappings($datastore, data?.FunctionArn);
  }

  const groups = [
  	{
  		group: [
  			{
  				name: "Function Name",
  				key: "FunctionName",
  			},
  			{
  				name: "ARN",
  				key: "FunctionArn",
  			},
  			{
  				name: "Handler",
  				key: "Handler",
  			},
  			{
  				name: "Run Time",
  				key: "Runtime",
  			},
  		],
  	},
  	{
  		group: [
  			{
  				name: "Memory Size",
  				key: "MemorySize",
  				isSize: true,
  			},
  			{
  				name: "Ephemeral Storage",
  				key: "EphemeralStorage",
  				subKey: "Size",
  				isSize: true,
  			},
  			{
  				name: "State",
  				key: "State",
  			},
  			{
  				name: "State Reason",
  				key: "StateReason",
  			},
  		],
  	},
  ];

  const renderCol = (col: {
    name: string;
    key: string;
    isSize?: boolean;
    subKey?: string;
  }) => {
  	const key = col.key as keyof LambdaFunctionProps["Functions"][0];
  	if (key === "EphemeralStorage") {
  		return (data?.EphemeralStorage.Size || "0") + " MB";
  	} else if (key === "MemorySize") {
  		return (data?.MemorySize || "0") + " MB";
  	}
  	return (data || {})[key] || "-";
  };
</script>

<Drawer
  {showModal}
  {closeModal}
  icon={{
  	src: "/assets/images/aws/lambda.png",
  	width: 50,
  	class: "p-2",
  	alt: "lambda",
  }}
  iconStyle={`background-color:${COLOR_SCHEME.SERVERLESS};`}
>
  <svelte:fragment slot="header">
    {#if data}
      <div class="ml-3">
        <Typography variant="h3" weight="semi-bold" font={30}>
          {data.FunctionName}
        </Typography>
        <Typography font={14} weight="medium" variant="div">
          Last Modifier {toLocaleDate(data.LastModified)}
        </Typography>
      </div>
    {/if}
  </svelte:fragment>
  {#if data}
    <Typography variant="h3" weight="medium" font={20}>
      Lambda Function
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
  <RenderConnectionsTable {connections} />
</Drawer>
