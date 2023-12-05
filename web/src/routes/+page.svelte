<script lang="ts">
  import AlertMessage from "$src/lib/components/common/Alerts/AlertMessage.svelte";
  import Button from "$src/lib/components/common/Button/Button.svelte";
  import Typography from "$src/lib/components/common/Typography/Typography.svelte";
  import CreateProject from "$src/lib/components/projects/CreateProject.svelte";
  import ProjectIndex from "$src/lib/components/projects/ProjectIndex.svelte";
  import { onMount } from "svelte";

  export let data;

  let showModal = false;

  let Alert: any;

  onMount(() => {
  	if (data.error) {
  		Alert?.alert(data.error);
  	}
  });
</script>

<div class="grid grid-cols-4 gap-8 mt-10 mb-8 container mx-auto">
  <div class="col-span-4 md:col-span-2">
    <Typography weight="bold" font={24} variant="h3" classname=""
      >Manage Your Projects</Typography
    >
    <Typography weight="regular" font={16} variant="p" classname="mt-5">
      Connect your favorite cloud provider to view and monitor your deployed
      services.
      <a href="/learnmore" class="hover:underline text-blue-600"
        >Learn more <span aria-hidden="true">→</span></a
      >
    </Typography>
    <Button
      classname="mt-3 bg-gradient text-white px-4 py-2"
      type="button"
      on:click={() => {
      	showModal = true;
      }}
    >
      Add Project
    </Button>
  </div>
  <div class="col-span-4 md:col-span-2 md:px-4">
    <Typography variant="div" weight="semi-bold" font={16}
      >Boost Your Account</Typography
    >
    <Typography classname="mt-3" variant="p" weight="regular" font={16}>
      Upgrade your account to unlock more projects, cloud providers from
      <span class="bg-yellow-500 text-white font-bold">AWS, GCP, AZURE</span>
      and collaborators.
      <a href="/learnmore" class="hover:underline text-blue-600"
        >Learn more <span aria-hidden="true">→</span></a
      >
    </Typography>
  </div>
</div>
<AlertMessage bind:this={Alert} />
{#if !data.error}
  <ProjectIndex projects={data.result || []} />
{/if}
{#if showModal}
  <CreateProject
    allowClose={true}
    onClose={() => {
    	showModal = false;
    }}
    onSave={(obj) => {
    	data.result.push(obj);
    }}
  />
{/if}
