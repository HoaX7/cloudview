<script lang="ts">
  import { onDestroy, onMount } from "svelte";

    export let callback: () => void;

    onMount(() => {
    	window.addEventListener("click", callback);
    });

    onDestroy(() => {
    	window.removeEventListener("click", callback);
    });

</script>