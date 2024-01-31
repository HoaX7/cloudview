<script lang="ts">
  import type { ProjectProps } from "$src/customTypes/projects";
  import Modal from "$lib/components/common/Modal/index.svelte";
  import Spinner from "../common/Loaders/Spinner.svelte";
  import Button from "../common/Button/Button.svelte";
  import Input from "../common/Input/Input.svelte";
  import Textarea from "../common/Input/Textarea.svelte";
	import AlertMessage from "../common/Alerts/AlertMessage.svelte";
	import { updateProject } from "$src/api/projects";
  import FormButtons from "../common/Modal/formButtons.svelte";

  export let project: ProjectProps;
  export let onClose: () => void;
  export let saving = false;
  export let onSave: (data: ProjectProps) => void;

  const data = {
  	name: project.name || "",
  	description: project.description || "",
  	email: project.email || ""
  };

  let Alert: any;

  const handleSubmit = async (e: CustomEvent<any>) => {
  	e.preventDefault();
  	try {
  		if (!data.name) {
  			alert("Please fill in all required fields");
  			return;
  		}
  		if (data.name === project.name && data.description === project.description && data.email === project.email) {
  			// nothing to save
  			return;
  		}
  		saving = true;
  		const body = { description: data.description } as any;
  		if (data.name != project.name) {
  			body.name = data.name;
  		}
  		if (data.email != project.email) {
  			body.email = data.email;
  		}
  		const resp = await updateProject(project.id, body);
  		if (resp.error) throw resp;
  		project.name = data.name;
  		project.email = data.email;
  		if (data.description) {
  			project.description = data.description;
  		}
  		Alert?.alert("Successfully updated", true);
  		onSave(project);
  		onClose();
  	} catch (err: any) {
  		console.error("Unable to save", err);
  		Alert?.alert(err?.message || "Unable to save data");
  	}
  	saving = false;
  };
</script>

<AlertMessage bind:this={Alert} />
<Modal title="Edit Project" closeModal={onClose} showButtons={true}
    isForm={true}
    on:submit={handleSubmit}
>
  <div class="text-start">
    <label for="name"><span class="text-red-600">*</span>Name:</label>
    <Input
      name="name"
      classname="mt-3"
      placeholder="Project Name"
      value={data.name}
      on:input={(e) => {
      	data.name = e.detail;
      }}
      required
      disabled={saving}
    />
  </div>
  <div class="text-start mt-3">
    <label for="email"><span class="text-red-600">*</span>Email:</label>
    <Input
      name="email"
      classname="mt-3"
      placeholder="Project Email"
      value={data.email}
      on:input={(e) => {
      	data.email = e.detail;
      }}
      required
      disabled={saving}
    />
  </div>
  <div class="text-start mt-3">
    <label for="description">Description:</label>
    <Textarea
      value={data.description}
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
  <svelte:fragment slot="buttons">
    <FormButtons {onClose} {saving} />
  </svelte:fragment>
</Modal>
