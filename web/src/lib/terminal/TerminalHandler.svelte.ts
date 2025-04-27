import type { Action } from 'svelte/action';
import { actionsList } from '$lib/types/Action';
import TerminalStore from '$lib/terminal/TerminalStore.svelte';
import { URLPrefix } from '$lib';

export const refreshTerminals = async () => {
	const res = await fetch(URLPrefix + '/api/actions');

	if (!res.ok) {
		console.error('Failed to fetch actions');
		return;
	}

	const data = actionsList.parse(await res.json());
	const keys = Object.keys(data.Actions);

	TerminalStore.filterTerminalsById(keys);

	const existingIds = TerminalStore.getTerminalIds();
	const newIds = keys.filter((id) => !existingIds.includes(id));

	for (const id of newIds) {
		const action = data.Actions[id];

		TerminalStore.addTerminal(action);
	}
};

const terminalHandler: Action<HTMLElement> = () => {
	let intervalRef: number | null = null;

	refreshTerminals().catch((err) => {
		console.error('Error while refreshing actions:', err);
	});

	intervalRef = setInterval(() => {
		refreshTerminals().catch((err) => {
			console.error('Error while refreshing actions:', err);
		});
	}, 60 * 1000);

	return {
		destroy() {
			if (intervalRef) {
				clearInterval(intervalRef);
			}
		}
	};
};

export default terminalHandler;
