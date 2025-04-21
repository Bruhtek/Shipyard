import TerminalHandles from '$lib/websocket/handles/terminal';

export type WSMessageHandle = {
	test: (message: object) => boolean;
	handle: (message: object) => void;
}

const Handles: WSMessageHandle[] = [
	...TerminalHandles
]

const handleMessage = (message: object) => {
	for (const handle of Handles) {
		if (handle.test(message)) {
			handle.handle(message);
			return;
		}
	}
}

export default handleMessage;