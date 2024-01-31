<script lang="ts">
  import type { ProjectProps } from "$src/customTypes/projects";
  import Typography from "../common/Typography/Typography.svelte";
  import Icon from "../common/Image/index.svelte";
  import clsx from "clsx";
  import Table from "../common/Table/Table.svelte";
  import { goto } from "$app/navigation";
	import EditProject from "./EditProject.svelte";
	import { clone } from "$src/helpers";

  export let projects: ProjectProps[];
  export let userId: string;

  let state = {
  	selectedProject: null,
  	loading: false,
  	tab: "you-own",
  	projects: projects.filter((p) => p.ownerId === userId),
  	showModal: false,
  };

  const columns = [
  	{
  		name: "Email",
  		key: "email",
  	},
  	// {
  	// 	name: "Members",
  	// 	key: "memberLimit",
  	// },
  	{
  		name: "Created on",
  		key: "createdAt",
  		isDate: true,
  	},
  	{
  		name: "Last edited",
  		key: "updatedAt",
  		isDate: true,
  	},
  ];

  const tabs = [
  	{
  		name: "You own",
  		key: "you-own",
  		filter: () =>
  			(state.projects = projects.filter((p) => p.ownerId === userId)),
  	},
  	{
  		name: "You're a member",
  		key: "member",
  		filter: () =>
  			(state.projects = projects.filter((p) => p.ownerId !== userId)),
  	},
  ];
</script>

<div class="mx-auto container border-t border-black pt-10">
  <div class="w-full">
    <div class="flex items-center">
      <Typography variant="h3" weight="semi-bold" font={24}>Projects</Typography
      >
      <nav class="ml-3 border-gray-300 border-l-2">
        {#each tabs as tab, index (index)}
          <Typography
            variant="span"
            font={14}
            weight="medium"
            classname={clsx(
            	"ml-3",
            	state.tab === tab.key
            		? "border-b-2 border-black"
            		: "hover:border-b-2 border-black"
            )}
          >
            <button
              on:click={() => {
              	state.tab = tab.key;
              	tab.filter();
              }}>{tab.name}</button
            >
          </Typography>
        {/each}
      </nav>
    </div>
    <Table data={state.projects} {columns}>
        <svelte:fragment slot="head-before-each">
        <td class="font-semibold p-3">Name</td>
        </svelte:fragment>
      <svelte:fragment slot="head">
        <td class="font-semibold p-3"> Actions </td>
      </svelte:fragment>
      <svelte:fragment slot="extra-row-td-before-each" let:item>
        <td class="p-3 flex items-center">
          <Icon
            src={`assets/images/${(item.type || "PRIVATE").toLowerCase()}.svg`}
            alt={item.type}
            class="mr-1"
            width="16"
          />
          <button
            class={clsx(
            	"font-bold hover:bg-teal-100 border-b-4 border-teal-100 cursor-pointer"
            )}
            on:click={() => goto(`/projects/${item.id}`, { state: item })}
          >
            {item.name}
          </button>
          {#if item.description}
            <div>
              {item.description}
            </div>
          {/if}
        </td>
      </svelte:fragment>
      <svelte:fragment slot="extra-row-td" let:item>
        <td class="p-3 flex items-center">
          {#if state.tab === "you-own"}
            <button on:click={() => {
            	state.showModal = true;
            	state.selectedProject = item;
            }}>
              <Icon src="/assets/images/edit.svg" width="24" alt="edit" />
            </button>
          {/if}
          <button
            class="ml-3"
            on:click={() => {
            	goto(`/projects/${item.id}`, { state: item });
            }}
          >
            <Icon src="/assets/images/view.svg" width="24" alt="view" />
          </button>
        </td>
      </svelte:fragment>
      <svelte:fragment slot="footer">
        {#if state.projects.length <= 0}
          <tr>
            <td colspan="5" class="text-center p-3">No data available</td>
          </tr>
        {/if}
      </svelte:fragment>
    </Table>
  </div>
</div>
{#if state.showModal && state.selectedProject}
<EditProject project={state.selectedProject} onClose={() => {
	state.showModal = false;
	state.selectedProject = null;
}} onSave={(obj) => {
	const res = clone(state.projects);
	const idx = res.findIndex((r) => r.id === obj.id);
	if (idx >= 0) {
		res[idx] = obj;
		state.projects = res;
	}
}} />
{/if}
