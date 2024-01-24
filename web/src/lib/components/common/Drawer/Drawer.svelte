<script lang="ts">
  import clsx from "clsx";
  import Icon from "../Image/index.svelte";

  export let showModal = false;
  export let closeModal: () => void;
  export let iconStyle = "";

  export let icon = {
  	src: "",
  	alt: "",
  	class: "",
  	width: 0,
  };
</script>

<div
  class={clsx(
  	"bg-white shadow-lg w-full lg:w-3/4 md:w-11/12 fixed top-0 right-0",
  	"z-20 h-full p-5 pt-0 overflow-y-auto transition-transform duration-500",
  	showModal ? "translate-x-0" : "translate-x-full"
  )}
>
<div class="border-b border-black sticky pb-5 top-0 pt-5 bg-white z-10">
  <div class="flex items-center">
    <button
      class=""
      on:click={() => {
      	showModal = false;
      	closeModal();
      }}
    >
      <Icon
        alt="back"
        src="/assets/images/right-arrow.svg"
        width={10}
        class="rotate-180"
      />
    </button>
    {#if icon.src}
    <Icon
        src={icon.src}
        alt={icon.alt}
        width={icon.width}
        class={clsx("ml-3 shadow rounded hidden md:block", icon.class)}
        style={iconStyle}
      />
    {/if}
    <slot name="header" />
  </div>
  <div class="mt-3 ml-3 flex items-center gap-4">
    <button class={clsx("rounded-full border-black border-2 p-2",
    	"disabled:opacity-25")} disabled>
      <img alt="start-instance" src="/assets/images/start.svg" width="20" />
    </button>
    <button class={clsx("rounded-full border-2 p-3 border-black",
    	"disabled:opacity-25")} disabled>
      <img alt="stop-instance" src="/assets/images/stop.svg" width="12" />
    </button>
    <button class={clsx("rounded-full border-2 border-black w-[40px] h-[40px]",
    	"flex items-center justify-center disabled:opacity-25")} disabled>
      <img alt="delete-instance" src="/assets/images/delete.svg" width="18" />
    </button>
  </div>
</div>
  <div class="p-5 mt-3">
    <slot />
  </div>
</div>
