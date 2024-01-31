<script lang="ts">
  import type { FontSize } from "$src/customTypes/typography";
    import Icon from "../common/Image/index.svelte";
  import Typography from "../common/Typography/Typography.svelte";

  export let width: number = 32;
  export let font: FontSize = 24;
</script>

<div class="flex items-center">
    <Icon src="/assets/images/cloudfriendly-logo.svg" alt="cloudfriendly-logo-ic" {width} /> 
    <Typography variant="div" {font} weight="bold" classname="ml-2 text-[#478af5]">
        Cloud Friendly
    </Typography>
</div>
