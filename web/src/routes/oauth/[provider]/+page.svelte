<!-- Dummy page -->
<!-- This Page is a dummy page which is used by login layout. The route must exist in order for 
Login component to work.

Contents on this page have no effect on the app. -->

<script lang="ts">
	import { goto } from "$app/navigation";
	import auth from "$src/store/auth";
	import { onMount } from "svelte";

    const user = auth.getUser();
    onMount(() => {
    	if ($user?.id) {
    		goto("/");
    		return;
    	}
    });
</script>
