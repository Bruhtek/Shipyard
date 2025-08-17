import WSDataStore, { type WSPayload } from '$lib/websocket/MessageHandler.svelte';

const ALLOWED_ACTIONS = ['up', 'stop', 'down', 'restart', 'pull'];

const StackAction = (environment: string, action: string, objectId: string) => {
	if (!environment || !action || !objectId) {
		throw new Error('Invalid parameters');
	}
	if (!ALLOWED_ACTIONS.includes(action)) {
		throw new Error('Invalid action');
	}

	const payload: WSPayload = {
		Environment: environment,
		Action: action,
		ObjectId: objectId,
		Object: 'stack'
	};

	WSDataStore.sendMessage(payload);
};

export default StackAction;
