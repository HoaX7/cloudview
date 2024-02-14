<script lang="ts">
  import { getProjectById } from "$src/api/projects.js";
  import ProviderAccountIndex from "$src/lib/components/services/ProviderAccountsIndex.svelte";
  import { onMount } from "svelte";

  /**
   * This route takes url params and query params of service id
   * to fetch details.
   * url ex: /projects/{id}?providerAccountId={id}
   */
  export let data;

  let project = history.state;
  const state = {
  	loading: false,
  	error: "",
  };

  /**
   * NOTE: If the `project` object was already fetched in the previous component
   * we can simply store the object in a `history` state.
  */
  onMount(async () => {
  	if (project.id !== data.projectId) {
  		try {
  			state.loading = true;
  			const result = await getProjectById({ id: data.projectId });
  			if (result.error || !result.data) throw result;
  			project = result.data;
  			history.replaceState(project, "");
  		} catch (err) {
  			console.error("Unable to fetch project", err);
  			state.error = "Unable to fetch project";
  		}
  		state.loading = false;
  	}
  });
</script>

<ProviderAccountIndex {project} providerAccounts={data.accounts} />
