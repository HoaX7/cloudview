<script lang="ts">
  import { getProjectById } from "$src/api/projects.js";
  import ServiceIndex from "$src/lib/components/services/ServiceIndex.svelte";
  import { onMount } from "svelte";

  /**
   * This route takes url params and query params of service id
   * to fetch details.
   * url ex: /projects/{id}?serviceId={id}
   */
  export let data;

  let project = history.state;
  const state = {
  	loading: false,
  	error: "",
  };

  onMount(async () => {
  	if (project.id !== data.projectId) {
  		try {
  			state.loading = true;
  			const result = await getProjectById({ id: data.projectId });
  			if (result.error || !result.data) throw result;
  			project = result.data;
  		} catch (err) {
  			console.error("Unable to fetch project", err);
  			state.error = "Unable to fetch project";
  		}
  		state.loading = false;
  	}
  });
</script>

<ServiceIndex {project} services={data.services} />
