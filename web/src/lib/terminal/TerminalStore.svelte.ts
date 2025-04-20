type Terminal = {
	id: string;
	content: string;
}

class terminalStore {
	// essentially just map Terminal(action) id to output
	terminals: Terminal[] = $state<Terminal[]>([]);
	subscriptions: { [key: string]: ((message: string) => void)[] } = {};

	addMessage(id: string, message: string) {
		// console.debug('Adding message', message);

		const terminal = this.terminals.find(t => t.id === id);
		if (terminal) {
			terminal.content += message;
		} else {
			this.terminals.push({ id, content: message });
		}

		// Notify subscribers
		if (this.subscriptions[id]) {
			this.subscriptions[id].forEach(callback => callback(message));
		}
	}

	subscribe (id: string, callback: (message: string) => void) {
		if (!this.subscriptions[id]) {
			this.subscriptions[id] = [];
		}
		this.subscriptions[id].push(callback);
	}

	unsubscribe (id: string, callback: (message: string) => void) {
		if (this.subscriptions[id]) {
			this.subscriptions[id] = this.subscriptions[id].filter(cb => cb !== callback);
		}
	}

	get terms () {
		return this.terminals;
	}
}

const TerminalStore = new terminalStore();
export default TerminalStore;
