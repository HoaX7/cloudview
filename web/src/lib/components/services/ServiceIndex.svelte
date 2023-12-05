<script lang="ts">
  import Button from "../common/Button/Button.svelte";
  import Table from "../common/Table/Table.svelte";
  import Typography from "../common/Typography/Typography.svelte";
  import Icon from "../common/Image/index.svelte";
  import PageNavButtons from "../common/Navigation/PageNavButtons.svelte";
  import type { ServiceProps } from "$src/customTypes/Services";
  import EditService from "./CreateEditService.svelte";
  import type { ProjectProps } from "$src/customTypes/Projects";
  import auth from "$src/store/auth";
  import Datastore from "$src/store/data";
  import { DEFAULT_REGION } from "$src/helpers/constants";
  import clsx from "clsx";
  import { goto } from "$app/navigation";
  import SettingsComponent from "../common/Settings/SettingsComponent.svelte";
  import CreateMember from "../members/CreateMember.svelte";
  import DeleteProject from "./DeleteProject.svelte";
  import { clone } from "$src/helpers";

  export let project: ProjectProps;
  export let services: ServiceProps[];

  const user = auth.getUser();
  const datastore = Datastore.getDatastore();
  if (!$datastore.selectedRegion) {
  	$datastore.selectedRegion = DEFAULT_REGION;
  }

  let state = {
  	showModal: false,
  	selectedService: null,
  	showSettingsMenu: false,
  	isCreate: false,
  	showMemberModal: false,
  	showDeleteProjectModal: false,
  };

  const columns = [
  	{
  		name: "Name",
  		key: "name",
  	},
  	{
  		name: "Description",
  		key: "description",
  	},
  	{
  		name: "Cloud Provider",
  		key: "provider",
  	},
  	{
  		name: "Last edited",
  		key: "updatedAt",
  		isDate: true,
  	},
  ];
</script>

<div class="mt-10 container mx-auto mb-8">
  {#if project}
    <div class="flex items-center justify-between w-full">
      <div class="flex items-center">
        <PageNavButtons back={true} className="mr-1" />
        <Typography
          weight="semi-bold"
          font={24}
          variant="h3"
          classname="flex items-center"
        >
          <Icon
            src={`/assets/images/${(
            	project.type || "PRIVATE"
            ).toLowerCase()}.svg`}
            alt={project.type}
            class="mr-1"
            width="16"
          />
          {project.name || "-"}
        </Typography>
        <Typography
          variant="div"
          font={16}
          weight="regular"
          classname="border-l-2 pl-3 ml-3 border-gray-300"
        >
          {project.email || "-"}
        </Typography>
      </div>
      {#if project.ownerId === $user?.id}
        <SettingsComponent showMenu={state.showSettingsMenu}>
          <div class="border-b">
            <Button
              classname={clsx(
              	"!p-3 !rounded-b-none !rounded-t-md",
              	"hover:bg-gray-100 w-full text-start !font-medium"
              )}
              on:click={(e) => {
              	state.showSettingsMenu = false;
              	state.showMemberModal = true;
              }}
            >
              <Icon
                src="/assets/images/plus.svg"
                alt="add"
                width="24"
                class="inline-block"
              /> Add members
            </Button>
            <Button
              classname={clsx(
              	"!p-3 !rounded-none",
              	"hover:bg-gray-100 w-full text-start !font-medium"
              )}
              on:click={() => {
              	goto(`/members?projectId=${project.id}`, { state: project });
              }}
            >
              <Icon
                src="/assets/images/view.svg"
                alt="add"
                width="18"
                class="inline-block mx-1"
              /> Manage Members
            </Button>
          </div>
          <Button
            classname={clsx(
            	"!p-3 !rounded-b-md text-red-500 hover:bg-red-500",
            	"hover:text-white w-full text-start",
            	"!rounded-t-none !font-medium"
            )}
            on:click={(e) => {
            	e.stopPropagation();
            	state.showDeleteProjectModal = true;
            }}
          >
            Delete Project
          </Button>
        </SettingsComponent>
      {/if}
    </div>
    <Typography variant="p" weight="regular" font={16} classname="mt-3">
      Connect and monitor your services running on your favorite cloud provider.
    </Typography>
    {#if project.ownerId === $user?.id}
      <Button
        classname="bg-gradient mt-3 text-white px-3 py-2"
        on:click={() => {
        	state.isCreate = true;
        	state.showModal = true;
        }}>Connect</Button
      >
    {:else}
      <Typography variant="div" weight="regular" font={16} classname="mt-3">
        <span class="bg-yellow-200"
          >Contact your Admin to connect to your cloud provider services.</span
        >
      </Typography>
    {/if}
    <div class="mt-10 border-t border-black pt-10">
      <Typography variant="h3" weight="semi-bold" font={18}>
        Services
      </Typography>
      <Table data={services} {columns}>
        <svelte:fragment slot="head">
          {#if $user?.id === project.ownerId}
            <td class="font-semibold p-3"> Actions </td>
          {/if}
        </svelte:fragment>
        <svelte:fragment slot="extra-row-td" let:item>
          <td class="p-3 flex items-center">
            {#if $user?.id === project.ownerId}
              <button
                on:click={(e) => {
                	state.selectedService = item;
                	state.showModal = true;
                }}
              >
                <Icon src="/assets/images/edit.svg" width="18" alt="edit" />
              </button>
            {/if}
            <a
              class="ml-3"
              href={`/cloud/${(item.provider || "").toLowerCase()}?serviceId=${
              	item.id
              }&projectId=${project.id}&region=${$datastore.selectedRegion}`}
              on:click={() => {
              	// Reset datastore nodes
              	// to avoid unnecessary side-effects of arrows
              	// being drawn to unknown positions
              	const res = clone($datastore);
              	res.konvaConnectableNodes = [];
              	res.konvaTargetFromNodes = [];
              	$datastore = res;
              }}
            >
              <Icon src="/assets/images/view.svg" width="18" alt="view" />
            </a>
          </td>
        </svelte:fragment>
        <svelte:fragment slot="footer">
          {#if services.length <= 0}
            <tr>
              <td colspan="5" class="text-center p-3">No data available</td>
            </tr>
          {/if}
        </svelte:fragment>
      </Table>
    </div>
    {#if state.showModal}
      <EditService
        isCreate={state.isCreate}
        selectedService={state.selectedService}
        onSave={(data) => {
        	const idx = services.findIndex((s) => s.id === data.id);
        	if (idx >= 0) {
        		services[idx].name = data.name;
        		services[idx].description = data.description;
        	} else {
        		services.push(data);
        	}
        }}
        onClose={() => {
        	state.showModal = false;
        	state.selectedService = null;
        	state.isCreate = false;
        }}
        projectId={project.id}
      />
    {/if}
    {#if state.showMemberModal}
      <CreateMember
        onClose={() => (state.showMemberModal = false)}
        projectId={project.id}
        onSave={() => {
        	return;
        }}
      />
    {/if}
    {#if state.showDeleteProjectModal && project}
      <DeleteProject
        {project}
        onClose={() => (state.showDeleteProjectModal = false)}
        onDelete={() => {
        	goto("/");
        }}
      />
    {/if}
  {:else}
    <Typography font={24} weight="semi-bold" variant="div">
      Oops, We could find the project you were looking for.
    </Typography>
  {/if}
</div>
