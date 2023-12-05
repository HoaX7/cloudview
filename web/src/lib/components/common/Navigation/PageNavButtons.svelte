<script lang="ts">
	import clsx from "clsx";
  /**
   * To navigate between pages with back/forward buttons.
   */
  import Icon from "../Image/index.svelte";

  export let back = false;
  export let className = "";
</script>

{#if back}
  <button
    on:click={(e) => {
    	e.stopPropagation();
    	history.back();
    }}
    class={clsx("hover:bg-gray-200 rounded-full p-2 w-8 h-8 text-center",
    	className)}
  >
    <Icon src="/assets/images/right-arrow.svg" width="10" class={"rotate-180"} alt="back" />
  </button>
{:else}
  <button
    on:click={(e) => {
    	e.stopPropagation();
    	history.forward();
    }}
    class={clsx("hover:bg-gray-200 rounded-full p-2 w-8 h-8 text-center",
    	className)}
  >
    <Icon src="/assets/images/right-arrow.svg" width="10" alt="forward" />
  </button>
{/if}
