<script lang="ts">
  import { createProjectWithService } from "$src/api/projects";
  import type { ProjectProps } from "$src/customTypes/Projects";
  import Datastore from "$src/store/data";
  import Button from "../common/Button/Button.svelte";
  import Input from "../common/Input/Input.svelte";
  import Textarea from "../common/Input/Textarea.svelte";
  import FullPageLoader from "../common/Loaders/FullPageLoader.svelte";
  import Spinner from "../common/Loaders/Spinner.svelte";
  import Modal from "../common/Modal/index.svelte";
  import Select from "../common/Select/Select.svelte";

  export let onSave: (data: ProjectProps) => void;
  export let allowClose = true;
  export let onClose: () => void;

  const data = {
  	name: "",
  	description: "",
  	provider: "AWS",
  	accessKeyId: "",
  	accessKeySecret: "",
  	type: "PRIVATE", // public,private
  };

  let saving = false;
  const datastore = Datastore.getDatastore();

  const handleSubmit = async (e: CustomEvent<any>) => {
  	e.preventDefault();
  	try {
  		saving = true;
  		const result = await createProjectWithService(data);
  		if (result.error) throw result;
  		if (result.data) {
  			onSave(result.data.project);
  		}
  	} catch (err: any) {
  		alert(err?.message);
  	}
  	saving = false;
  };
</script>

{#if saving}
  <FullPageLoader />
{/if}
<Modal
  title="Create your Project"
  isForm={true}
  showButtons={true}
  closeModal={() => {
  	onClose();
  }}
  showCloseBtn={allowClose}
  on:submit={handleSubmit}
>
  <div class="text-start">
    <label for="name"><span class="text-red-600">*</span>Name:</label>
    <Input
      name="name"
      classname="mt-3"
      placeholder="Project Name"
      on:input={(e) => {
      	data.name = e.detail;
      }}
      required
      disabled={saving}
    />
  </div>
  <div class="text-start mt-3">
    <label for="description">Description:</label>
    <Textarea
      name="description"
      classname="mt-3 h-24"
      placeholder="Short Description"
      maxlength="160"
      on:input={(e) => {
      	data.description = e.detail;
      }}
      disabled={saving}
    />
  </div>
  <div class="mt-3 text-start">
    <label for="provider"
      ><span class="text-red-600">*</span> Cloud Provider:</label
    >
    <Select name="provider" classname="mt-3">
      <option>AWS</option>
    </Select>
  </div>
  <div class="text-start mt-3">
    <label for="accessId"><span class="text-red-600">*</span>Access ID:</label>
    <Input
      name="accessId"
      classname="mt-3"
      on:input={(e) => {
      	data.accessKeyId = e.detail;
      }}
      required
      disabled={saving}
    />
  </div>
  <div class="text-start mt-3">
    <label for="accessSecretKey"
      ><span class="text-red-600">*</span>Access Secret Key:</label
    >
    <Input
      name="accessSecretKey"
      classname="mt-3"
      on:input={(e) => {
      	data.accessKeySecret = e.detail;
      }}
      required
      disabled={saving}
    />
  </div>
  <svelte:fragment slot="buttons">
    <Button
      type="submit"
      classname="px-3 py-2 !rounded text-blue-600 hover:bg-gray-100"
      disabled={saving}
    >
      {#if saving}
        <Spinner size="xxs" />
      {:else}
        Save
      {/if}
    </Button>
    {#if allowClose}
      <Button
        type="button"
        on:click={() => onClose()}
        classname="px-3 py-2 !rounded text-blue-600 hover:bg-gray-100"
        disabled={saving}
      >
        Cancel
      </Button>
    {/if}
  </svelte:fragment>
</Modal>
