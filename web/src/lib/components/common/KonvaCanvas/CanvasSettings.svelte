<script lang="ts">
  import Icon from "$lib/components/common/Image/index.svelte";
  import SettingStore from "$src/store/settings";
  import WindowEvents from "../Hooks/WindowEvents.svelte";

  let showModal = false;

  const store = SettingStore.getStore();
  let checked = Boolean($store.animate || false);

  $: $store.animate = checked;
</script>

<WindowEvents callback={() => {
	showModal = false;
}} />
<div class="relative flex items-center">
  {#if showModal}
    <div class="absolute bottom-8 rounded p-3 right-0 bg-white shadow">
      <div class="flex items-center">
        <label for="animate" class="mr-2"> Animate </label>
        <input name="animate" type="checkbox" bind:checked={checked} on:click={(e) => {
        	e.stopPropagation();
        }} />
      </div>
    </div>
  {/if}
  <button on:click={(e) => {
  	e.stopPropagation();
  	showModal = !showModal;
  }} class="help-text hover:before:content-['settings'] hover:before:w-[60px]" >
    <Icon src="/assets/images/settings.svg" alt="settings" width="28" />
  </button>
</div>
