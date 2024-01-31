<script lang="ts">
	import type { SessionProps } from "$src/customTypes/user";
	import Image from "$lib/components/common/Image/index.svelte";
	import Typography from "../common/Typography/Typography.svelte";
	import Button from "../common/Button/Button.svelte";
	import clsx from "clsx";
	import Spinner from "../common/Loaders/Spinner.svelte";
	import { logout } from "$src/api/oauth";
	import Auth from "$src/store/auth";
	import FullPageLoader from "../common/Loaders/FullPageLoader.svelte";
	import Datastore from "$src/store/data";
  import WindowEvents from "../common/Hooks/WindowEvents.svelte";

	export let profile: SessionProps;
	let showMenu = false;
	let saving = false;
</script>

{#if saving}
	<FullPageLoader />
{/if}
<!-- <WindowEvents callback={() => {
	if (showMenu) showMenu = false;
}} /> -->
{#if profile?.id}
	<div class="relative">
		<Button
			type="button"
			classname="shadow flex items-center bg-white py-1 px-2 hover:bg-gray-200"
			on:click={(e) => {
				e.detail.stopPropagation();
				showMenu = !showMenu;
			}}
		>
			<Image class="rounded-full h-8 w-8" src={profile.avatarUrl} alt={profile.username} />
			<Typography variant="span" weight="regular" font={16} classname="ml-3 truncate">
				{profile.username || ""}
			</Typography>
		</Button>
		{#if showMenu}
			<div class={clsx("z-20 top-12 rounded absolute bg-white shadow md:w-80 w-40 right-2", "p-3")}>
				<button
					class="absolute right-3 top-1"
					on:click={(e) => {
						e.stopPropagation();
						showMenu = false;
					}}
				>
					&times;</button
				>
				<div class="flex items-start mt-5">
					<Image
						class="rounded-full border-2 md:block hidden mr-3"
						width={80}
						height={80}
						src={profile.avatarUrl}
						alt={profile.username}
					/>
					<div>
						<Typography variant="div" weight="semi-bold" font={16}>
							{profile.username || ""}
						</Typography>
						<Typography variant="div" weight="regular" font={16}>
							{profile.email || ""}
						</Typography>
					</div>
				</div>
				<ul class="border-t mt-5">
					<li class="mt-2">
						<button
							disabled={saving}
							on:click={async (e) => {
								e.stopPropagation();
								try {
									saving = true;
									await logout();
									Auth.logout();
									Datastore.clear();
									window.location.reload();
								} catch (err) {
									console.log("Unable to logout", err);
								}
								saving = false;
								showMenu = false;
							}}
							class={clsx("hover:bg-gray-100 w-full text-start p-1", "disabled:cursor-not-allowed")}
							>Logout
							{#if saving}
								<Spinner size="xxs" className="mr-3 inline-block" />
							{/if}
						</button>
					</li>
				</ul>
			</div>
		{/if}
	</div>
{/if}
