import type { Terminal } from '$lib/terminal/TerminalStore.svelte';
import { URLPrefix } from '$lib';
import { ActionStatus } from '$lib/types/Action';

// either stops if running, else simply deletes
export const DismissTerminal = async (t: Terminal): Promise<boolean> => {
	// it has already been deleted on the backend
	if (t.MarkedForDeletion) {
		return true;
	}

	if (t.Status !== ActionStatus.Running) {
		t.MarkedForDeletion = true;
	}

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