<script lang="ts">
  import Login from "$src/lib/components/layout/Login.svelte";
  import { Login as apiLogin, getUserSession } from "$src/api/oauth";
  import "../css/app.css";
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import FullPageLoader from "$src/lib/components/common/Loaders/FullPageLoader.svelte";
  import Profile from "$lib/components/profile/index.svelte";
  import auth from "$src/store/auth";
  import Datastore from "$src/store/data";
  import SelectProviderAndRegion from "$src/lib/components/layout/SelectProviderAndRegion.svelte";
  import KonvaStore from "$src/store/konva";
  import SettingStore from "$src/store/settings";
  import Brandlogo from "$src/lib/components/layout/Brandlogo.svelte";
  import { page } from "$app/stores";
  import MainLayout from "$src/lib/components/layout/MainLayout.svelte";
  import { SEO } from "$src/helpers/seo";

  export let data;
  // Initialize auth store
  let user = auth.initializeStore();
  Datastore.init();
  KonvaStore.init();
  SettingStore.init();
  let isLoggedIn = false;
  //
  // Added a Loading state to avoid unnecessary layout shifts/component changes.
  //
  let isLoading = true;
  if ($user?.id && data.cookie) {
  	isLoggedIn = true;
  }

  onMount(async () => {
  	/**
     * If session storage data is lost but cookie
     * exists, try to fetch only the session data.
     */
  	if (!isLoggedIn && data.cookie) {
  		console.log(
  			"cookie was found but no session data. re-fetching session data..."
  		);
  		try {
  			const resp = await getUserSession();
  			if (resp.error || !resp.data) throw resp;
  			$user = resp.data;
  			isLoggedIn = true;
  		} catch (err) {
  			console.log("Unable to fetch session:", err);
  			isLoggedIn = false;
  		}
  	}
  	isLoading = false;
  	if (!isLoggedIn && data.code) {
  		const provider = data.params.provider;
  		if (provider !== "github" && provider !== "google") {
  			console.error("Invalid provider, please re-login");
  			return;
  		}
  		try {
  			isLoading = true;
  			const result = await apiLogin({
  				provider: provider === "google" ? "google" : "github",
  				code: data.code,
  			});
  			if (result.error) throw result;
  			$user = result.data;
  			isLoggedIn = true;
  			goto("/");
  		} catch (err) {
  			console.error("Unable to login", err);
  			isLoggedIn = false;
  		}
  		isLoading = false;
  	}
  });
</script>

<MainLayout
  title={SEO.home.title}
  description={SEO.home.description}
  keywords={SEO.home.keywords}
  url={$page.url.pathname}
>
  {#if isLoading}
    <FullPageLoader />
  {:else if isLoggedIn && $user}
    <nav class="p-3 bg-gray-100 relative z-20 shadow">
      <div class="container mx-auto flex items-center justify-between">
        <a href="/"><Brandlogo /></a>
        {#if $page.data.showRegionDD}
          <SelectProviderAndRegion />
        {/if}
        <Profile profile={$user} />
      </div>
    </nav>
    {#if data.params.provider}
      <slot />
      <!-- to show error message popups -->
      <div id="modal-root" />
    {:else}
      <div class="container mx-auto h-full p-3 md:p-4 md:max-w-[1200px]">
        <slot />
      </div>
    {/if}
  {:else}
    <Login />
  {/if}
</MainLayout>
