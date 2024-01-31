import isEmpty from "lodash/isEmpty";
import isEqual from "lodash/isEqual";
import uniqueWith from "lodash/uniqWith";
import cloneDeep from "lodash/cloneDeep";
import groupBy from "lodash/groupBy";
import type { ResourceDataReturnType } from "$src/customTypes/services";

export const isBrowser = () => typeof window !== "undefined";

export const toLocaleDate = (date: Date | string) =>
	new Date(date).toLocaleTimeString("en-us", {
		year: "2-digit",
		month: "short",
		day: "numeric",
	});

export const isEmptyObject = (val = {}) => isEmpty(val);

export const getRandomNumber = (max: number) => Math.floor(Math.random() * max);

export const uniqueArray = (array: any[]) => uniqueWith(array, isEqual);

export const clone = <T>(value: T) => {
	return cloneDeep(value);
};

export const delay = (ms: number) =>
	new Promise((resolve) => setTimeout(resolve, ms));

export const bytesToMegaBytes = (bytes: number) =>
	Number((bytes / 1024 ** 2).toFixed(2));

export function debounce<Params extends any[]>(
	func: (...args: Params) => any,
	timeout = 300
): (...args: Params) => void {
	let timer: NodeJS.Timeout;
	return (...args: Params) => {
		clearTimeout(timer);
		timer = setTimeout(() => {
			func(...args);
		}, timeout);
	};
}

export const reorderArray = <T>(array: T[], priorityOrder: any[]) => {
	if (!Array.isArray(array)) return array;
	if (!Array.isArray(priorityOrder)) return priorityOrder;
	const temp = [] as any[];
	priorityOrder.map((opt) => {
		const idx = array.findIndex((it) => it === opt);
		if (idx >= 0) temp.push(array[idx]);
	});
	return temp;
};

export const reorderAwsServices = (
	array: ResourceDataReturnType,
	order: string[]
) => {
	const temp = [] as ResourceDataReturnType;
	order.forEach((name) => {
		const idx = array.findIndex((it) => it.name === name);
		if (idx >= 0) temp.push(array[idx]);
	});
	return temp;
};

export const groupByKey = <T>(array: T[], key: keyof T) => {
	return groupBy(array, key);
};

export const isMobile = () => {
	return Math.min(window.screen.width, window.screen.height) < 768;
};
