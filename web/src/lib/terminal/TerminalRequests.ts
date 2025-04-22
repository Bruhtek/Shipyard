import type { Terminal } from '$lib/terminal/TerminalStore.svelte';
import { URLPrefix } from '$lib';

// either stops and dismisses if running, else simply dismisses
export const DismissTerminal = async (t: Terminal): Promise<boolean> => {
	// it has already been deleted on the backend
	if (t.MarkedForDeletion) {
		return true;
	}

	t.MarkedForDeletion = true;

	const res = await fetch(URLPrefix + `/api/actions/${t.id}`, {
		method: 'DELETE',
	})

	if (res.ok || res.status === 404) {
		return true;
	}

	return false;
}

const TerminalRequests = {
	DismissTerminal,
}

export default TerminalRequests;