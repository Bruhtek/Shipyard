import { nanoid } from 'nanoid';
import handleMessage from '$lib/websocket/handles';

export enum ConnectionStatus {
	DISCONNECTED = 0,
	CONNECTED = 1,
	RECONNECTING = 2
}

export type WSPayload = {
	Environment: string;
	Action: string;
	ObjectId: string;
	Object: string;
};

class MessageHandler {
	connectionStatus: ConnectionStatus = $state<ConnectionStatus>(0);
	_sendMessage: ((data: object) => void) | null = null;

	_setSendMessage(sendMessage: ((data: object) => void) | null) {
		this._sendMessage = sendMessage;
	}

	constructor() {
		this.connectionStatus = 0;

		// if, after 5 seconds, we are still disconnected, change the status to show it to the user
		setTimeout(() => {
			if (this.connectionStatus === ConnectionStatus.DISCONNECTED) {
				this.connectionStatus = ConnectionStatus.RECONNECTING;
			}
		}, 5000);
	}

	addMessage(message: string) {
		try {
			const json = JSON.parse(message);
			if (!json.Type) {
				return;
			}

			try {
				handleMessage(json);
			} catch (e) {
				console.error('Failed to handle message', e);
			}
			// eslint-disable-next-line @typescript-eslint/no-unused-vars
		} catch (_) {
			// silently ignore
		}
	}

	setConnectedStatus(status: ConnectionStatus) {
		this.connectionStatus = status;
	}

	async sendMessage(payload: WSPayload) {
		if (!this._sendMessage) {
			throw new Error('WebSocket not connected');
		}

		this._sendMessage(payload);
	}
}

const WSDataStore = new MessageHandler();
export default WSDataStore;
