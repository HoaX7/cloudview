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
      <svg
        data-accordion-icon
        class={clsx(
        	showDD ? reverseDD ? "rotate-0" : "rotate-180" : "rotate-90",
        	"w-3 h-3 shrink-0 ml-2"
        )}
        aria-hidden="true"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 10 6"
      >
        <path
          stroke="currentColor"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M9 5 5 1 1 5"
        />
      </svg>
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
