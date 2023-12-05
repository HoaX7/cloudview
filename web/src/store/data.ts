/**
 * Data-store acts as a small caching system to store projects
 * and selected services information.
 */

import type { DatastoreProps } from "$src/customTypes/Store";
import type { Writable } from "svelte/store";
import storable from ".";

const INITIAL_STATE = {
	selectedRegion: "",
	fetchData: false,
	konvaConnectableNodes: [],
	konvaTargetFromNodes: [],
	dragNodeId: null,
} as DatastoreProps;

/**
 * This store is being used by canvas to draw connecting lines
 * between nodes.
 * 
 * If you want to add additional different type of data
 * it is recommended to create a different store.
 */
const Datastore = {
	store: {} as Writable<DatastoreProps>,
	init() {
		if (typeof window === "undefined") return;
		const store = storable(INITIAL_STATE, "data-store", localStorage);
		if (store) {
			this.store = store;
		}
		return store; 
	},
	update(data: DatastoreProps) {
		this.store.set(data);
	},
	clear(region?: string) {
		if (region) {
			INITIAL_STATE.selectedRegion = region;
		}
		this.store.set(INITIAL_STATE);
	},
	getDatastore() {
		return this.store;
	}
};

export default Datastore;