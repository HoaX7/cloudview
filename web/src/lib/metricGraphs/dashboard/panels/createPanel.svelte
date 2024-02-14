<script lang="ts">
  import Input from "$src/lib/components/common/Input/Input.svelte";
  import Typography from "$src/lib/components/common/Typography/Typography.svelte";
  import Icon from "$lib/components/common/Image/index.svelte";
  import clsx from "clsx";
</script>

<div class="mt-3">
  <label for="name"><span class="text-red-600">*</span>Name:</label>
  <Input
    required
    id="name"
    on:input={(e) => {}}
    placeholder="Start typing..."
  />
</div>
<div class="mt-3">
  <label for="resource"><span class="text-red-600">*</span>Resource:</label>
  <select id="resource">
    <option> Ec2 </option>
  </select>
</div>
<div class="mt-3">
  <Typography
    variant="p"
    font={16}
    weight="semi-bold"
    classname={clsx("flex items-center gap-1 bg-yellow-200 p-3 rounded-md")}
  >
    <Icon src="/assets/images/info.svg" alt="info" width="24" /> To enable real-time
    monitoring, follow the instructions below:
  </Typography>
</div>
