<script lang="ts">
  import type { ServiceProps } from "$src/customTypes/Services";
  import Button from "../common/Button/Button.svelte";
  import Input from "../common/Input/Input.svelte";
  import Textarea from "../common/Input/Textarea.svelte";
  import Modal from "../common/Modal/index.svelte";
  import Accordion from "../common/Accordion/Accordion.svelte";
  import Select from "../common/Select/Select.svelte";
  import { createService, editService } from "$src/api/services";
  import { isEmptyObject } from "$src/helpers";
  import Spinner from "../common/Loaders/Spinner.svelte";
  import AlertMessage from "../common/Alerts/AlertMessage.svelte";
  import clsx from "clsx";
  import FormButtons from "../common/Modal/formButtons.svelte";

  export let selectedService: ServiceProps | null;
  export let isCreate = false;
  export let projectId: string;
  export let onSave: (data: ServiceProps) => void;
  export let onClose: () => void;

  const data = {
  	name: selectedService?.name || "",
  	description: selectedService?.description || "",
  	provider: selectedService?.provider || "AWS",
  	accessKeyId: "",
  	accessKeySecret: "",
  } as ServiceProps;

  let state = {
  	saving: false,
  	error: "",
  };

  let Alert: any;

  const edit = async () => {
  	if (!selectedService) return;
  	const body = {} as ServiceProps;
  	Object.keys(data).forEach((key) => {
  		const k = key as keyof ServiceProps;
  		if (selectedService && data[k] !== selectedService[k] && data[k] !== "") {
  			Object.assign(body, { [k]: data[k] });
  		}
  	});
  	if (isEmptyObject(body)) {
  		console.log("Nothing to edit");
  		onClose();
  		return;
  	}
  	state.saving = true;
  	await editService(selectedService.id, {
  		...body,
  		projectId: projectId,
  	});
  	onSave({
  		...data,
  		id: selectedService.id,
  	});
  };

  const handleSubmit = async (e: CustomEvent<any>) => {
  	try {
  		e.preventDefault();
  		if (isCreate) {
  			if (!data.accessKeyId || !data.accessKeySecret) {
  				alert("Please fill all the required fields.");
  				return;
  			}
  			const res = await createService({
  				projectId: projectId,
  				provider: "AWS",
  				name: data.name,
  				accessKeyId: data.accessKeyId,
  				accessKeySecret: data.accessKeySecret,
  			});
  			if (res.error || !res.data) throw res;
  			onSave(res.data);
  		} else {
  			await edit();
  		}
  		Alert?.alert("Successfully saved", true);
  		onClose();
  	} catch (err: any) {
  		console.error("Unable to save service", err);
  		Alert?.alert(err?.message || "Unable to save changes");
  	}
  	state.saving = false;
  };
</script>

<AlertMessage bind:this={Alert} />
<Modal
  title={isCreate ? "Create Service" : "Edit Service"}
  closeModal={onClose}
  isForm={true}
  showButtons={true}
  on:submit={handleSubmit}
  description="Your Access keys cannot be retrieved once they are saved."
>
  <div class="mt-3 text-start">
    <label for="name"><span class="text-red-600">*</span>Name:</label>
    <Input
      name="name"
      classname="mt-3"
      value={data.name}
      placeholder="Start Typing..."
      required
      on:input={(e) => {
      	data.name = e.detail;
      }}
    />
  </div>
  <div class="mt-3 text-start">
    <label for="description">Description:</label>
    <Textarea
      name="description"
      classname="mt-3 h-24"
      value={data.description || ""}
      placeholder="Short Description"
      maxlength="160"
      on:input={(e) => {
      	data.description = e.detail;
      }}
    />
  </div>
  <Accordion title="Advanced" containerClass="mt-3" buttonClass="px-3 !rounded" showDD={isCreate}>
    <div class="mt-3 text-start">
      <label for="accessKeyId"
        ><span class={clsx("text-red-600", isCreate ? "" : "hidden")}>*</span
        >Access Key ID:</label
      >
      <Input
        name="accessKeyId"
        classname="mt-3"
        placeholder="Start Typing..."
        on:input={(e) => {
        	data.accessKeyId = e.detail;
        }}
        required={isCreate}
      />
    </div>
    <div class="mt-3 text-start">
      <label for="accessKeySecret"
        ><span class={clsx("text-red-600", isCreate ? "" : "hidden")}>*</span
        >Access Key Secret:</label
      >
      <Input
        type="password"
        name="accessKeySecret"
        classname="mt-3"
        placeholder="Start Typing..."
        on:input={(e) => {
        	data.accessKeySecret = e.detail;
        }}
        required={isCreate}
      />
    </div>
    <div class="mt-3 text-start">
      <label for="provider">Provider:</label>
      <Select classname="mt-3" name="provider">
        <option>AWS</option>
      </Select>
    </div>
  </Accordion>
  <svelte:fragment slot="buttons">
    <FormButtons saving={state.saving} {onClose} />
  </svelte:fragment>
</Modal>
