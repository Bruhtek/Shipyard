import type { WSMessageHandle } from '$lib/websocket/handles/index';
import { actionMessage, actionMessageMetadata, actionRemovedMessage } from '$lib/types/Action';
import TerminalStore from '$lib/terminal/TerminalStore.svelte';

const actionOutput: WSMessageHandle = {
	test: (message: object) => {
		return "Type" in message && message.Type === 'ActionOutput'
	},
	handle: (message: object) => {
		const {success, data} = actionMessage.safeParse(message);
		if (!success) {
			return;
		}

		TerminalStore.addMessage(data.ActionId, data.Message + '\r');
	}
}
const actionMetadata: WSMessageHandle = {
	test: (message: object) => {
		return "Type" in message && message.Type === 'ActionMetadata';
	},
	handle: (message: object) => {
		const {success, data} = actionMessageMetadata.safeParse(message);
		if (!success) {
			return;
		}

		TerminalStore.addTerminal(data.Metadata);
	}
}

const actionRemoved: WSMessageHandle = {
	test: (message: object) => {
		return "Type" in message && message.Type === 'ActionRemoved';
	},
	handle: (message: object) => {
		const {success, data} = actionRemovedMessage.safeParse(message);
		if (!success) {
			return;
		}

		TerminalStore.removeTerminal(data.ActionId)
	}
}



const TerminalHandles = [
	actionOutput,
	actionMetadata,
	actionRemoved
]

export default TerminalHandles