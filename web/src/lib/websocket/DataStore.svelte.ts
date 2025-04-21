import { nanoid } from 'nanoid';
import TerminalStore from '$lib/terminal/TerminalStore.svelte';
import { actionMessage, actionMessageMetadata } from '$lib/types/Action';

export enum ConnectionStatus {
	DISCONNECTED = 0,
	CONNECTED = 1,
	RECONNECTING = 2
}

class DataStore {
	messages: string[] = $state<string[]>([]);
	connectionStatus: ConnectionStatus = $state<ConnectionStatus>(0);
	_sendMessage: ((data: object) => void) | null = null;

	_setSendMessage(sendMessage: ((data: object) => void) | null) {
		this._sendMessage = sendMessage;
	}

	constructor() {
		this.messages = [];
		this.connectionStatus = 0;

		// if, after 5 seconds, we are still disconnected, change the status to show it to the user
		setTimeout(() => {
			if (this.connectionStatus === ConnectionStatus.DISCONNECTED) {
				this.connectionStatus = ConnectionStatus.RECONNECTING;
			}
		}, 5000);
	}

	addMessage(message: string) {
		// console.debug('Adding message', message);

		try {
			const json = JSON.parse(message);
			if (!json.ActionId) {
				return;
			}

			try {
				if (json.Message) {
					const parsedMessage = actionMessage.parse(json);

					TerminalStore.addMessage(parsedMessage.ActionId, parsedMessage.Message + '\r');
				} else if (json.Metadata) {
					const metadata = actionMessageMetadata.parse(json);

					console.log(metadata);
					TerminalStore.addTerminal(metadata);
				}
			} catch (e) {
				console.error('Failed to add message', e);
			}
			// eslint-disable-next-line @typescript-eslint/no-unused-vars
		} catch (_) {
			// silently ignore
		}
	}

	setConnectedStatus(status: ConnectionStatus) {
		this.connectionStatus = status;
	}

	async fetch(env: string, object: string, action: string) {
		if (!this._sendMessage) {
			throw new Error('Disconnected');
		}

		const actionId = nanoid();

		const data = {
			Environment: env,
			Object: object,
			Action: action,
			ActionId: actionId
		};

		this._sendMessage(data);
	}
}

const WSDataStore = new DataStore();
export default WSDataStore;