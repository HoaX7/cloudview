<script lang="ts">
  import { DEFAULT_REGION } from "$src/helpers/constants.js";
  import { clone } from "$src/helpers/index.js";
  import Button from "$src/lib/components/common/Button/Button.svelte";
  import Typography from "$src/lib/components/common/Typography/Typography.svelte";
  import DisplayCloudData from "$src/lib/components/services/DisplayCloudData.svelte";
  import Datastore from "$src/store/data.js";
  import { onMount } from "svelte";

  export let data;

  const datastore = Datastore.getDatastore();

  onMount(() => {
  	if ($datastore && $datastore.selectedRegion !== data.region) {
  		const res = clone($datastore);
  		res.selectedRegion = data.region || DEFAULT_REGION;
  		res.konvaConnectableNodes = [];
  		res.konvaConnectableNodes = [];
  		$datastore = res;
  	}
  });
</script>

{#if data.missingParams}
  <Typography
    weight="semi-bold"
    font={18}
    variant="div"
    classname="flex flex-col items-center justify-center h-full"
  >
    Please select a valid Project and Service
    <Button
      classname="mt-3 hover:bg-gray-100 !rounded bg-gray-200 !py-2 !px-3 font-medium text-sm"
      type="button"
      on:click={() => {
      	history.back();
      }}
    >
      Go back
    </Button>
  </Typography>
{:else}
  <DisplayCloudData
    result={data.metricData || []}
    projectId={data.projectId || ""}
    serviceId={data.serviceId || ""}
  />
{/if}
