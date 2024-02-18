<script lang="ts">
  import Button from "../common/Button/Button.svelte";
  import Typography from "../common/Typography/Typography.svelte";
  import FullPageLoader from "../common/Loaders/FullPageLoader.svelte";
  import Brandlogo from "./Brandlogo.svelte";
  import Input from "../common/Input/Input.svelte";
  import { Login } from "$src/api/oauth";
  import AlertMessage from "../common/Alerts/AlertMessage.svelte";
  import { PUBLIC_WEBAPP_SERVICE } from "$env/static/public";

  let isLoading = false;

  let state = {
  	username: "",
  	password: "",
  };

  let alertRef: any;
  // eslint-disable-next-line max-len
  const demoUrl = `${PUBLIC_WEBAPP_SERVICE}/cloud/metrics?providerAccountId=9f16f1f2-27fd-44bb-a33d-5c1b489b9e37&projectId=14f2c455-b4f4-4e43-93c8-30d7146088d4&region=us-west-2`;

  const handleSubmit = async (e: SubmitEvent) => {
  	try {
  		e.preventDefault();
  		isLoading = true;
  		const res = await Login({
  			provider: "local",
  			code: "dummy-code",
  			username: state.username,
  			password: state.password,
  		});
  		if (res.error || !res.data) throw res;
  		state = {
  			username: "",
  			password: "",
  		};
  		window.location.assign(demoUrl);
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
        Visualize and Monitor your resources deployed on AWS, GCP, Azure
        and more in one place
      </Typography>
      <form class="mt-3 w-full" on:submit={handleSubmit}>
        <Input
          id="username"
          placeholder="Username"
          classname=""
          on:input={(e) => {
          	state.username = e.detail;
          }}
        />
        <Input
          id="password"
          placeholder="Password"
          type="password"
          classname="mt-3"
          on:input={(e) => {
          	state.password = e.detail;
          }}
        />
        <Button
          type="submit"
          classname="bg-gradient text-white mt-3 text-sm"
          disabled={!state.username || !state.password}
        >
          See Demo
        </Button>
      </form>
    </div>
  </div>
</main>
