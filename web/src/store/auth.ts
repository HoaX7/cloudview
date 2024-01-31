import type { SessionProps } from "$src/customTypes/user";
import storable from "./index";
import type { Writable } from "svelte/store";

/**
 * Reponsible for handling/storing user data in store
 * to persist until session ends or user logs out.
 * 
 * It is important to call Auth.logout() and make the API
 * to logout to remove session and cookies.
 */
const INITIAL_STATE = {} as SessionProps;

export default {
	store: {} as Writable<SessionProps>,
	initializeStore() {
		if (typeof window === "undefined") return;
		const store = storable(INITIAL_STATE, "store", sessionStorage);
		if (store) {
			this.store = store;
		}
		return store;
	},
	logout() {
		this.store.set(INITIAL_STATE);
	},
	getUser() {
		if (typeof window === "undefined") return;
		return this.store;
	}
};
