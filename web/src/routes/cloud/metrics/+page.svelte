<script lang="ts">
  import PageNavButtons from "$lib/components/common/Navigation/PageNavButtons.svelte";
  import DashboardIndex from "$lib/metricGraphs/dashboard/index.svelte";
  import Typography from "$src/lib/components/common/Typography/Typography.svelte";
  import Icon from "$lib/components/common/Image/index.svelte";
  import clsx from "clsx";
  import { copyText } from "$src/utils/index.js";
  export let data;
</script>

{#if data}
  <div class="flex items-center gap-3">
    <Icon src="/assets/images/dashboard.svg" alt="dashboard" width="24" />
  <Typography variant="div" font={18} weight="medium">
    <a href={`/projects/${data.projectId}`} class="hover:underline">
      {data.providerAcc?.project.name || "-"}
    </a>
  </Typography>
  <span>/</span>
  <Typography variant="div" font={18} weight="medium">
    <a
      href={`/cloud/aws?providerAccountId=${data.providerAccountId}` +
      `&projectId=${data.projectId}&region=${data.region}`}
      class="hover:underline"
    >
      {data.providerAcc?.name || "-"}
    </a>
  </Typography>
  <button class={clsx("help-text hover:before:-bottom-6 hover:before:left-5",
  	"hover:before:content-['share'] hover:before:w-[60px] focus:hover:before:content-['copied']")}
      on:click={(e) => {
      	copyText(window.location.href);
      }}
    >
    <Icon alt="share" src="/assets/images/share.svg" width="24" />
  </button>
  </div>
  <Typography classname="mt-5" font={16} weight="regular" variant="p">
    Customize and Monitor your resource usage in one single place, in real-time.
  </Typography>
  <DashboardIndex resourceList={data.metricData} data={data.providerAcc?.metadata || {}} />
{/if}
