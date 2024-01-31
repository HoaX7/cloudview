<script lang="ts">
  import { inviteProjectMember } from "$src/api/projectMembers";
	import type { ProjectMemberApiProps } from "$src/customTypes/projectMembers";
  import AlertMessage from "../common/Alerts/AlertMessage.svelte";
  import Button from "../common/Button/Button.svelte";
  import Input from "../common/Input/Input.svelte";
  import Spinner from "../common/Loaders/Spinner.svelte";
  import FormButtons from "../common/Modal/formButtons.svelte";
  import Modal from "../common/Modal/index.svelte";

  export let projectId: string;
  export let onClose: () => void;
  export let onSave: (data: ProjectMemberApiProps) => void;

  let email = "";
  let error = "";
  let saving = false;

  const handleSubmit = async (e: CustomEvent<any>) => {
  	try {
  		e.preventDefault();
  		if (!email || !projectId) {
  			alert("Please select a valid project and enter a valid email.");
  			return;
  		}
  		saving = true;
  		const res = await inviteProjectMember({
  			email,
  			projectId,
  		});
  		if (res.error || !res.data) throw res;
  		onSave(res.data);
  		onClose();
  	} catch (err: any) {
  		console.error("Unable to invite members", err);
  		error = err?.message || "Unable to invite members";
  		alert(error);
  	}
  	saving = false;
  };
</script>

<Modal
  isForm={true}
  title="Send Invite"
  on:submit={handleSubmit}
  closeModal={onClose}
  showButtons={true}
>
  <div class="mt-3">
    <label for="email"><span class="text-red-600">*</span>Email:</label>
    <Input
      type="email"
      required
      placeholder="Start Typing..."
      classname="mt-3"
      on:input={(e) => {
      	email = e.detail;
      }}
    />
  </div>
  <svelte:fragment slot="buttons">
    <FormButtons {saving} {onClose} />
  </svelte:fragment>
</Modal>
