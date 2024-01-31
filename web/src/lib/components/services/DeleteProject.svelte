<script lang="ts">
  import type { ProjectProps } from "$src/customTypes/projects";
  import Modal from "$lib/components/common/Modal/index.svelte";
  import Typography from "../common/Typography/Typography.svelte";
	import Input from "../common/Input/Input.svelte";
	import Button from "../common/Button/Button.svelte";
	import Spinner from "../common/Loaders/Spinner.svelte";
	import AlertMessage from "../common/Alerts/AlertMessage.svelte";
	import { updateProject } from "$src/api/projects";

  export let project: ProjectProps;
  export let onClose: () => void;
  export let onDelete: (id: string) => void;

  let saving = false;

  let name = "";
  let Alert: any;

  const handleDelete = async (e: CustomEvent<any>) => {
  	try {
  		e.preventDefault();
  	if (name !== project.name) return;
  		const resp = await updateProject(project.id, { isDeleted: true });
  		if (resp.error) throw resp;
  		Alert?.alert(`${project.name} has been deleted`, true);
  	    onDelete(project.id);
  	} catch (err: any) {
  		console.error("Unable to delete project", err);
  		Alert?.alert(err?.message || "Unabel to delete project");
  	}
  };
</script>

<AlertMessage bind:this={Alert} />
<Modal
  title="Are you ABSOLUTELY sure?"
  description="Unexpected bad things will happen if you don't read this!"
  descriptionclass="bg-yellow-200 font-medium p-3"
  closeModal={onClose}
  showButtons={true}
  isForm={true}
  on:submit={handleDelete}
>
  <Typography variant="p" weight="regular" font={16} classname="mt-3">
    This action <span class="font-bold">CANNOT</span> be undone. This will
    permanently delete the <span class="font-bold">{project.name}</span> project,
    services and all collaborators associated.
  </Typography>
  <Typography variant="p" weight="regular" font={18} classname="mt-3">
    Please type in the name of the project to confirm.
  </Typography>
  <Input 
    required
    classname="mt-3"
    on:input={(e) => {
    	name = e.detail;
    }}
  />
  <svelte:fragment slot="buttons">
    <Button
      type="submit"
      classname="px-3 py-2 !rounded text-red-500 hover:bg-gray-100"
      disabled={saving || (name != project.name)}
    >
      {#if saving}
        <Spinner size="xxs" />
      {:else}
        I understand the consequences, delete this project
      {/if}
    </Button>
  </svelte:fragment>
</Modal>
