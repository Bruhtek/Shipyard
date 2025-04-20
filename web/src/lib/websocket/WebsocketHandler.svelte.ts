import type { Action } from 'svelte/action';
import { dev } from '$app/environment';
import WSDataStore, { ConnectionStatus } from '$lib/websocket/DataStore.svelte';

const websocketHandler: Action<HTMLElement> = () => {
	$effect(() => {
		const url = dev ? 'localhost:4000' : window.location.host;

		const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
		let socket: WebSocket;
		let shouldReconnect = true;

		const initialize = () => {
			socket = new WebSocket(`${protocol}://${url}/ws`);

			const sendMessage = (data: object) => {
				const message = JSON.stringify(data)

				if (socket.readyState === WebSocket.OPEN) {
					socket.send(message);
				} else {
					console.error('[WS] Socket not open, cannot send message');
				}
			};

			socket.onopen = () => {
				console.log('[WS] Connected');
				WSDataStore._setSendMessage(sendMessage);
				WSDataStore.setConnectedStatus(ConnectionStatus.CONNECTED);
			};
			socket.onclose = () => {
				console.log('[WS] Disconnected');
				WSDataStore._setSendMessage(null);
				WSDataStore.setConnectedStatus(ConnectionStatus.RECONNECTING);
				if (shouldReconnect) {
					initialize(); // try to Reconnect
				}
			};
			socket.onmessage = (e: MessageEvent) => {
				WSDataStore.addMessage(e.data.toString());
			};
		};

		initialize();

		return () => {
			shouldReconnect = false;
			socket.close();
		};
	});
};

export default websocketHandler;