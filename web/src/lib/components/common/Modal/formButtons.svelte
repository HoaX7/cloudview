<script lang="ts">
  import Button from "../Button/Button.svelte";
  import Spinner from "../Loaders/Spinner.svelte";

  export let saving = false;
  export let onClose: () => void;
  export let allowClose = true;
</script>

<Button
  type="submit"
  classname="px-3 py-2 !rounded text-blue-500 hover:bg-gray-100"
  disabled={saving}
>
  {#if saving}
    <Spinner size="xxs" />
  {:else}
    Save
  {/if}
</Button>
{#if allowClose}
  <Button
    type="button"
    on:click={() => onClose()}
    classname="px-3 py-2 !rounded text-blue-500 hover:bg-gray-100"
    disabled={saving}
  >
    Cancel
  </Button>
{/if}
