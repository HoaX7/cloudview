<script lang="ts">
  import type { ProjectMemberApiProps } from "$src/customTypes/ProjectMembers";
  import clsx from "clsx";
  import Button from "../common/Button/Button.svelte";
  import SettingsComponent from "../common/Settings/SettingsComponent.svelte";
  import Table from "../common/Table/Table.svelte";
  import Typography from "../common/Typography/Typography.svelte";
  import CreateMember from "./CreateMember.svelte";
  import { clone } from "$src/helpers";
  import { updateProjectMember } from "$src/api/projectMembers";
	import AlertMessage from "../common/Alerts/AlertMessage.svelte";

  export let result: ProjectMemberApiProps[];
  export let showModal = false;
  export let closeModal: () => void;
  export let projectId: string;
  export let showDD = false;

  let Alert: any;

  const columns = [
  	{
  		name: "Name",
  		key: "user.username",
  		keyName: "user",
  		subKey: "username",
  	},
  	{
  		name: "Email",
  		key: "user.email",
  		keyName: "user",
  		subKey: "email",
  	},
  	{
  		name: "Joined Team On",
  		key: "createdAt",
  		isDate: true,
  	},
  	{
  		name: "Last Edited",
  		key: "updatedAt",
  		isDate: true,
  	},
  ];

  const handleToggleStatus = async (item: ProjectMemberApiProps) => {
  	try {
  		const resp = await updateProjectMember(
  			{ id: item.id },
  			{
  				isActive: !item.isActive,
  				projectId,
  			}
  		);
  		if (resp.error) throw resp;
  		const res = clone(result);
  		const idx = res.findIndex((r) => r.id === item.id);
  		if (idx >= 0) {
  			res[idx].isActive = !item.isActive;
  			result = res;
  		}
  		Alert?.alert("Successfully saved", true);
  	} catch (err: any) {
  		console.error("Unable to save data", err);
  		Alert?.alert(err?.message || "Unable to save data");

  	}
  };

  const handleRemoveMember = async (item: ProjectMemberApiProps) => {
  	try {
  		const resp = await updateProjectMember(
  			{ id: item.id },
  			{
  				isDeleted: true,
  				projectId,
  			}
  		);
  		if (resp.error) throw resp;
  		const res = clone(result);
  		const idx = res.findIndex((r) => r.id === item.id);
  		if (idx >= 0) {
  			res.splice(idx, 1);
  			result = res;
  		}
  		Alert?.alert("Member removed successfully", true);
  	} catch (err: any) {
  		console.error("Unable to remove member", err);
  		Alert?.alert(err?.message || "Unable to remove member");
  	}
  };
</script>

<AlertMessage bind:this={Alert} />

<Typography variant="p" weight="regular" font={16} classname="mt-3">
  Add and Manage Members on your Team. Upgrade your Account to invite more
  members.
</Typography>
<div class="mt-10 border-t border-black pt-10">
  <Typography variant="h3" weight="semi-bold" font={18}>Team Members</Typography
  >
  <Table {columns} data={result}>
    <svelte:fragment slot="head">
      <td class="p-3 font-semibold"> Status </td>
      <td class="p-3 font-semibold"> Actions </td>
    </svelte:fragment>
    <svelte:fragment slot="extra-row-td" let:item>
      <td class="p-3 border-b">
        {item.isActive ? "Active" : "Deactivated"}
      </td>
      <td class="p-3 border-b">
        {#if !item.isOwner}
        <SettingsComponent width={24} showMenu={showDD}>
            <Button
              classname={clsx(
              	"!rounded-none !p-3 hover:bg-gray-100 w-full text-start",
              	"!font-medium !rounded-t-lg border-b"
              )}
              on:click={(e) => {
              	e.stopPropagation();
              	handleToggleStatus(item);
              }}
            >
              {item.isActive ? "Deactivate" : "Activate"}
            </Button>
            <Button
              classname={clsx(
              	"!p-3 !rounded-none text-red-500 w-full",
              	"hover:bg-red-500 hover:text-white text-start !font-medium",
              	"!rounded-b-lg"
              )}
              on:click={(e) => {
              	e.stopPropagation();
              	handleRemoveMember(item);
              }}
            >
              Remove
            </Button>
          </SettingsComponent>
        {:else}
            -
        {/if}
      </td>
    </svelte:fragment>
    <svelte:fragment slot="footer">
      {#if result.length <= 0}
        <tr>
          <td colspan="5" class="text-center p-3">No data available</td>
        </tr>
      {/if}
    </svelte:fragment>
  </Table>

  {#if showModal}
    <CreateMember
      {projectId}
      onClose={closeModal}
      onSave={(obj) => {
      	const res = clone(result);
      	res.push(obj);
      	result = res;
      }}
    />
  {/if}
</div>
