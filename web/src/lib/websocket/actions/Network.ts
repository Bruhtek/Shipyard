import WSDataStore, { type WSPayload } from '$lib/websocket/MessageHandler.svelte';

const ALLOWED_ACTIONS = ['remove'];

const NetworkAction = (environment: string, action: string, ...objectIds: string[]) => {
	if (!environment || !action || !objectIds) {
		throw new Error('Invalid parameters');
	}
	if (!ALLOWED_ACTIONS.includes(action)) {
		throw new Error('Invalid action');
	}

	const objectId = objectIds.join(',');

	const payload: WSPayload = {
		Action: action,
		ObjectId: objectId,
		Environment: environment,
		Object: 'network'
	};

	WSDataStore.sendMessage(payload);
};

export default NetworkAction;
