<script lang="ts">
  import type { ResourceDataReturnType } from "$src/customTypes/services";
  import Icon from "$lib/components/common/Image/index.svelte";


    export let resourceList: ResourceDataReturnType = [];
    export let data: any;

    console.log({
    	data,
    	resourceList 
    });
</script>

<div class="mt-5">
    <button class="rounded-full py-2 px-3 bg-gradient text-white text-sm font-bold">
        Add Panel
    </button>
</div>
