<script lang="ts">
  import Modal from "../Modal/index.svelte";
  let state = {
  	message: "",
  	status: "",
  	showModal: false,
  };

  export const alert = (message: string, success: boolean) => {
  	state.showModal = true;
  	state.message = message;
  	state.status = success ? "Success" : "Error";
  };
</script>

{#if state.showModal}
  <Modal
    className="!z-30"
    description={state.message}
    title={state.status}
    icon={state.status.toLowerCase()}
    closeModal={() => {
    	state.showModal = false;
    	state.status = "";
    	state.message = "";
    }}
  />
{/if}
