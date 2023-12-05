<script lang="ts">
	import { toLocaleDate } from "$src/helpers";
	import clsx from "clsx";

    export let columns: {
        name: string;
        key: string;
        isDate?: boolean;
        extraKey?: string;
        class?: string;
        subKey?: string;
        keyName?: string;
    }[];
    export let data: any[];
</script>

<table class="table-fixed border-collapse w-full">
    <thead class="border-b text-sm">
        <tr class="">
            <slot  name="head-before-each"/>
            {#each columns as col (col.key)}
                <td class="font-semibold p-3 break-words">
                    {col.name}
                </td>
            {/each}
            <slot name="head" />
        </tr>
    </thead>
    <tbody>
        <slot name="body" />
        {#each data as item, index (index)}
            <tr class={clsx("hover:bg-gray-100 border-b", 
            	typeof item.isActive === "boolean" && !item.isActive ? "text-gray-600" : "")}>
                <slot name="extra-row-td-before-each" {item} />
                {#each columns as col, index (index)}
                    <td class="p-3 break-words">
                        {#if col.isDate}
                            {toLocaleDate(item[col.key])}
                        {:else if col.subKey && col.keyName}
                            <span class={col.class}>{(item[col.keyName] || {})[col.subKey] || "-"}</span>
                        {:else}
                            <span class={col.class}>{item[col.key] || "-"}</span>
                        {/if}
                    </td>
                {/each}
                <slot name="extra-row-td" {item} />
            </tr>
        {/each}
    </tbody>
    <tfoot>
        <slot name="footer" />
    </tfoot>
</table>
