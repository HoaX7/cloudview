/**
 * konva stores acts as a medium to be able to draw correct
 * bounding rectangle backgrounds for services.
 */

import type { KonvaStoreProps } from "$src/customTypes/store";
import type { Writable } from "svelte/store";
import storable from ".";

const INITIAL_STATE = {
	externalBoundingRect: {
		x: 0,
		y: 0,
		width: 0,
		height: 0
	},
	internalBoundingRect: {
		x: 0,
		y: 0,
		width: 0,
		height: 0
	},
	rowCount: {
		internal: 0,
		external: 0
	}
} as KonvaStoreProps;

const KonvaStore = {
	store: {} as Writable<KonvaStoreProps>,
	init() {
		if (typeof window === "undefined") return;
		const store = storable(INITIAL_STATE, "konva-store", localStorage);
		if (store) {
			this.store = store;
		}
		return store; 
	},
	update(data: KonvaStoreProps) {
		this.store.set(data);
	},
	clear() {
		this.store.set(INITIAL_STATE);
	},
	getStore() {
		return this.store;
	}
};

export default KonvaStore;

