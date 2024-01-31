/**
 * konva stores acts as a medium to be able to draw correct
 * bounding rectangle backgrounds for services.
 */

import type { SettingStoreProps } from "$src/customTypes/store";
import type { Writable } from "svelte/store";
import storable from ".";

const INITIAL_STATE = { animate: true } as SettingStoreProps;

const SettingStore = {
	store: {} as Writable<SettingStoreProps>,
	init() {
		if (typeof window === "undefined") return;
		const store = storable(INITIAL_STATE, "settings-store", localStorage);
		if (store) {
			this.store = store;
		}
		return store; 
	},
	update(data: SettingStoreProps) {
		this.store.set(data);
	},
	clear() {
		this.store.set(INITIAL_STATE);
	},
	getStore() {
		return this.store;
	}
};

export default SettingStore;

