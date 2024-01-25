<script lang="ts">
  import type { MetricDataReturnType } from "$src/customTypes/Services";
  import Icon from "$lib/components/common/Image/index.svelte";


    export let resourceList: MetricDataReturnType = [];
    export let data: any;
</script>

<div class="mt-5">
    <button class="rounded-full py-2 px-3 bg-gradient text-white text-sm font-bold">
        Add Panel
    </button>
</div>
