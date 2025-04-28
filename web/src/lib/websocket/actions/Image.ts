import WSDataStore, { type WSPayload } from '$lib/websocket/MessageHandler.svelte';

const ALLOWED_ACTIONS = ['rm', 'pull'];

const ImageAction = (environment: string, action: string, objectId: string) => {
	if (!environment || !action || !objectId) {
		throw new Error('Invalid parameters');
	}
	if (!ALLOWED_ACTIONS.includes(action)) {
		throw new Error('Invalid action');
	}

	const payload: WSPayload = {
		Action: action,
		ObjectId: objectId,
		Environment: environment,
		Object: 'image'
	};

	WSDataStore.sendMessage(payload);
};

export default ImageAction;
