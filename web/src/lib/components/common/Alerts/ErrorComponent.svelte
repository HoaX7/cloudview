<script lang="ts">
	import Typography from "../Typography/Typography.svelte";

	export let error: string;
</script>

<Typography variant="div" weight="regular" font={16} classname="text-red-500 italic">
	{error}
</Typography>
