<script lang="ts">
  import PageNavButtons from "$src/lib/components/common/Navigation/PageNavButtons.svelte";
  import Typography from "$src/lib/components/common/Typography/Typography.svelte";
  import Icon from "$lib/components/common/Image/index.svelte";
  import Button from "$src/lib/components/common/Button/Button.svelte";
  import MembersIndex from "$src/lib/components/members/MembersIndex.svelte";
  import SettingsComponent from "$src/lib/components/common/Settings/SettingsComponent.svelte";
  import clsx from "clsx";

  export let data;

  let project = history.state;
  const state = {
  	loading: false,
  	error: "",
  	showModal: false,
  	showMenu: false
  };
</script>

{#if !data?.projectId || data.error}
  <Typography
    weight="semi-bold"
    font={18}
    variant="div"
    classname="flex flex-col items-center justify-center h-full"
  >
    {data.error || "Please select a valid Project"}
    <Button
      classname="mt-3 hover:bg-gray-100 !rounded bg-gray-200 !py-2 !px-3 font-medium text-sm"
      type="button"
      on:click={() => {
      	history.back();
      }}
    >
      Go back
    </Button>
  </Typography>
{:else}
  <div class="mt-10 container mx-auto mb-8">
    <div class="flex items-center justify-between">
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
      <SettingsComponent showMenu={state.showMenu}>
        <Button
          classname={clsx(
          	"!p-3 rounded-lg",
          	"hover:bg-gray-100 w-full text-start !font-medium"
          )}
          on:click={() => {
          	state.showModal = true;
          	state.showMenu = false;
          }}
        >
          <Icon
            src="/assets/images/plus.svg"
            alt="add"
            width="24"
            class="inline-block"
          /> Add members
        </Button>
      </SettingsComponent>
    </div>
    <MembersIndex
      projectId={data.projectId}
      result={data.result}
      showModal={state.showModal}
      closeModal={() => (state.showModal = false)}
    />
  </div>
{/if}
