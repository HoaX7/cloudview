<script lang="ts">
	import clsx from "clsx";
	import Typography from "../Typography/Typography.svelte";
	import Success from "./success.icon.svelte";
	import Error from "./error.icon.svelte";
	import { createEventDispatcher } from "svelte";

	export let isForm = false;
	export let disabled = false;
	export let title: string;
	export let icon: string = ""; // success or error
	export let descriptionclass: string = "";
	export let description: string = "";
	export let modalClass = "";
	export let closeModal: () => void;
	export let showButtons = false;
	export let showCloseBtn = true;
	let variant = "div";
	if (isForm) {
		variant = "form";
	}

	export let className = "";

	const dispatch = createEventDispatcher();
</script>

<div class={clsx("relative z-20", className)} aria-labelledby="modal-title" role="dialog" aria-modal="true">
	<div class="fixed inset-0 bg-gray-600 bg-opacity-75 transition-opacity" />
	<div class="fixed z-10 inset-0 overflow-y-auto">
		<div class="flex items-end sm:items-center justify-center min-h-full text-center sm:p-0">
			<div
				class={clsx(
					"relative w-full rounded-top-corners md:rounded-lg",
					"text-left overflow-hidden shadow-xl",
					"transform transition-all sm:my-8 sm:max-w-lg sm:w-full"
				)}
			>
				<svelte:element
					this={variant}
					class={clsx("custom-shadow-md bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4", modalClass)}
					on:submit={(e) => {
						e.preventDefault();
						console.log("Default prevented:", e.defaultPrevented);
						dispatch("submit", e);
					}}
				>
					{#if showCloseBtn}
						<button
							class={clsx(
								"py-1 px-3 right-5 rounded text-right",
								"hover:bg-gray-100 absolute cursor-pointer",
								disabled ? "pointer-events-none opacity-75" : ""
							)}
							on:click={() => closeModal()}
						>
							&times;
						</button>
					{/if}
					<div class={clsx("sm:flex sm:items-start w-full")}>
						<div class="mt-3 text-center sm:mt-0 sm:text-left w-full">
							{#if icon === "success"}
								<Success />
							{:else if icon === "error"}
								<Error />
							{/if}
							<Typography variant="h3" weight="semi-bold" font={24} classname="border-b border-gray-200">
								{title}
							</Typography>
							<div class="mt-2">
								<Typography 
								variant="p" 
								weight="regular" 
								font={16} 
								classname={clsx("text-gray-600", descriptionclass)}>
									{description}
								</Typography>
								<slot />
							</div>
						</div>
					</div>
					{#if showButtons}
						<div class="py-3 sm:flex sm:flex-row-reverse">
							<slot name="buttons" />
						</div>
					{/if}
				</svelte:element>
			</div>
		</div>
	</div>
</div>
