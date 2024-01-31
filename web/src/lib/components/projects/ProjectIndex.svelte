<script lang="ts">
  import type { ProjectProps } from "$src/customTypes/projects";
  import auth from "$src/store/auth";
  import CreateProject from "./CreateProject.svelte";
  import SelectProject from "./SelectProject.svelte";

  export let projects: ProjectProps[] = [];
  const user = auth.getUser();
  let isCreate = true;
  if (projects.length > 0) {
  	isCreate = false;
  }
</script>

{#if isCreate}
  <CreateProject
    onSave={(obj) => {
    	projects.push(obj);
    }}
    allowClose={false}
	onClose={() => {
		return;
	}}
  />
{:else}
  <SelectProject projects={projects || []} userId={$user?.id || ""} />
{/if}
