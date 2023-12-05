<script lang="ts">
    import css from "./Typography.module.css";
	import type { TypographyProps } from "$src/customTypes/Typography";
	import clsx from "clsx";

    export let variant: TypographyProps["variant"];
    export let weight: TypographyProps["weight"];
    export let font: TypographyProps["font"];
    export let classname: TypographyProps["class"] = "";

</script>


<svelte:element this={variant} class={clsx(classname, css[weight], css[`font-${font}`])}>
    <slot />
</svelte:element>
