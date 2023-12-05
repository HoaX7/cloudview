<script lang="ts">
    import Icon from "../Image/index.svelte";

    export let showMenu = false;

    export let width = 32;
</script>

<div class="relative">
    <button class="hover:bg-gray-100 rounded-full"
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
