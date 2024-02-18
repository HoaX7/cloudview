<script lang="ts">
  import clsx from "clsx";
  import Button from "../Button/Button.svelte";

  export let title: string;
  export let containerClass = "";
  export let buttonClass = "";
  export let showDD = false;
  export let disableDefaultBtnClass = false;
  export let reverseDD = false;
</script>

<div id="accordion-collapse" data-accordion="collapse" class={clsx(containerClass, "relative w-auto")}>
  <div id="accordion-collapse-heading-1" class="">
    {#if reverseDD}
      <div
        id="accordion-collapse-body-1"
        class={showDD ? "absolute bottom-10 right-0 md:left-0 w-[250px]" : "hidden"}
        aria-labelledby="accordion-collapse-heading-1"
      >
      <button class="absolute top-3 right-5" on:click={() => showDD = false}>
        &times;
      </button>
        <slot />
      </div>
    {/if}
    <Button
      type="button"
      classname={clsx(
      	!disableDefaultBtnClass &&
          "flex items-center justify-between w-full p-3",
      	!disableDefaultBtnClass && "bg-gray-100 hover:bg-gray-200 rounded",
      	buttonClass
      )}
      data-accordion-target="#accordion-collapse-body-1"
      aria-expanded="false"
      aria-controls="accordion-collapse-body-1"
      on:click={() => {
      	showDD = !showDD;
      }}
    >
      <span>{title}</span>
      <img src="/assets/images/right-arrow.svg" width="6" alt="dd"
        class={clsx(showDD ? reverseDD ? "-rotate-90" : "rotate-90" : "", "ml-2")}
      />
    </Button>
  </div>
  {#if !reverseDD}
    <div
      id="accordion-collapse-body-1"
      class={showDD ? "" : "hidden"}
      aria-labelledby="accordion-collapse-heading-1"
    >
      <slot />
    </div>
  {/if}
</div>
