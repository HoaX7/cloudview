<script lang="ts">
	import clsx from "clsx";

    export let classname = "";
</script>

<select class={clsx("border-2 rounded w-full p-2", classname)} {...$$restProps}>
    <slot />
</select>