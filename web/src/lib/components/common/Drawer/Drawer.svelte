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
  	"z-10 bg-white shadow-lg w-full md:w-11/12 fixed top-0 right-0",
  	"z-20 h-full p-5 pt-0 overflow-y-auto transition-transform duration-500",
  	showModal ? "translate-x-0" : "translate-x-full"
  )}
>
  <div class="flex items-center border-b border-black sticky pb-5 top-0 pt-5 bg-white z-10">
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
  <div class="p-5 mt-3">
    <slot />
  </div>
</div>
