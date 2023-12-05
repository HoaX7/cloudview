<script lang="ts">
	import clsx from "clsx";
    import css from "./Loader.module.css";
    export let className = "";
    export let size: "xxs" | "xs" | "sm" | "md" | "lg" | "xl" = "sm";


const classMap = {
	xxs: {
		h: "h-4",
		w: "w-4"
	},
	xs: {
		h: "h-8",
		w: "w-8"
	},
	sm: {
		h: "h-12",
		w: "w-12"
	},
	md: {
		h: "h-16",
		w: "w-16"
	},
	lg: {
		h: "h-24",
		w: "w-24"
	},
	xl: {
		h: "h-34",
		w: "w-34"
	}
};
</script>

<div class={clsx(
	css["loader"],
	"ease-linear rounded-full border-2 border-t-2",
	classMap[size].h, classMap[size].w,
	className
)}
></div>
