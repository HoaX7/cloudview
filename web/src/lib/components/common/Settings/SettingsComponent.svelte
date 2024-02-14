<script lang="ts">
  import { onDestroy, onMount } from "svelte";
    import Icon from "../Image/index.svelte";
  import WindowEvents from "../Hooks/WindowEvents.svelte";
  import clsx from "clsx";

    export let showMenu = false;

    export let width = 32;

    const closeModal = () => {
    	showMenu = false;
    };
    
</script>

<WindowEvents callback={closeModal} />
<div class="relative">
    <button class={clsx("hover:bg-gray-100 rounded-full help-text",
    	"hover:before:content-['settings'] hover:before:-bottom-8",
    	"hover:before:w-[60px]")}
    on:click={(e) => {
    	e.stopPropagation();
    	showMenu = !showMenu;
    }}
>
    <Icon
        src="/assets/images/settings.svg"
        alt="settings"
        width={width}
    />
</button>
{#if showMenu}
<div class="absolute shadow-2xl w-60 right-0 z-10 bg-white border rounded-lg">
    <slot />
</div>
{/if}
</div>
