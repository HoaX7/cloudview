<script lang="ts">
	import Button from "../common/Button/Button.svelte";
	import Typography from "../common/Typography/Typography.svelte";
	import FullPageLoader from "../common/Loaders/FullPageLoader.svelte";
  import Brandlogo from "./Brandlogo.svelte";
  import Input from "../common/Input/Input.svelte";
  import { Login } from "$src/api/oauth";
  import AlertMessage from "../common/Alerts/AlertMessage.svelte";

	let isLoading = false;

    let state = {
    	username: "",
    	password: ""
    };

    let alertRef: any;

    const handleSubmit = async (e: SubmitEvent) => {
    	try {
    		e.preventDefault();
    		isLoading = true;
    		const res = await Login({
    			provider: "local",
    			code: "dummy-code",
    			username: state.username,
    			password: state.password
    		});
    		if (res.error || !res.data) throw res;
    		state = {
    			username: "",
    			password: ""
    		};
    	} catch (err: any) {
    		console.error("Unable to login", err);
    		alertRef.alert(err?.message || "Unable to login", false);
    	}
    	isLoading = false;
    };
</script>

<main class="container mx-auto h-full">
    <AlertMessage bind:this={alertRef} />
	{#if isLoading}
		<FullPageLoader />
	{/if}
	<div class="flex flex-col items-center justify-center h-full">
		<div class="w-80 flex flex-col justify-center text-center items-center">
			<Brandlogo font={32} />
			<Typography variant="p" weight="medium" font={16} classname="mt-3">
				Visualize and Monitor your micro services deployed on AWS, GCP, Azure and more in one place
			</Typography>
            <form class="mt-3 w-full" on:submit={handleSubmit}>
                <Input id="username" placeholder="Username" classname=""
                    on:input={(e) => {
                    	state.username = e.detail;
                    }}
                />
                <Input id="password" placeholder="Password" type="password" classname="mt-3" 
                    on:input={(e) => {
                    	state.password = e.detail;
                    }}
                />
                <Button type="submit" classname="bg-gradient text-white mt-3 text-sm"
                    disabled={!state.username || !state.password}
                >
                    See Demo
                </Button>
            </form>
		</div>
	</div>
</main>
