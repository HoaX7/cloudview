<script lang="ts">
  import clsx from "clsx";

  export let showDD = false;
  export let containerClass = "";
  export let title: string;
</script>

<div
  id="accordion-collapse"
  data-accordion="collapse"
  class={clsx(containerClass, "relative w-auto")}
>
  <div id="accordion-collapse-heading-1" class="">
    <button
      type="button"
      data-accordion-target="#accordion-collapse-body-1"
      aria-expanded="false"
      aria-controls="accordion-collapse-body-1"
      on:click={() => {
      	showDD = !showDD;
      }}
      class={clsx("flex items-center justify-between w-full p-1 px-2",
      	"bg-gray-100 hover:bg-gray-200 rounded-md")}
    >
      <span>{title}</span>
      <img src="/assets/images/right-arrow.svg" width="6" alt="dd" 
        class={showDD ? "rotate-90" : ""}
      />
    </button>
    <div
      id="accordion-collapse-body-1"
      class={showDD ? "" : "hidden"}
      aria-labelledby="accordion-collapse-heading-1"
    >
      <slot />
    </div>
  </div>
</div>
