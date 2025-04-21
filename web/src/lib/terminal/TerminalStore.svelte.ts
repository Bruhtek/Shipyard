import dayjs from 'dayjs';

import type { ActionMetadata } from '$lib/types/Action';

export type Terminal = {
	id: string;
	content: string;

	// command details
	Action: string;
	Environment: string;
	Object: string;

	// metadata
	FinishedAt: dayjs.Dayjs;
	StartedAt: dayjs.Dayjs;

	Status: number;
}

export const terminalStatus = (t: Terminal) => {
	switch (t.Status) {
		case 0:
			return 'Pending';
		case 1:
			return 'Running';
		case 2:
			return 'Success';
		case 3:
			return 'Failed';
		default:
			return 'Unknown';
	}
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
			this.terminals.push({
				id,
				content: message,
				Action: '',
				Environment: '',
				Object: '',

				FinishedAt: dayjs().set('year', 1900),
				StartedAt: dayjs(),
				Status: -1,
			});
		}

		// Notify subscribers
		if (this.subscriptions[id]) {
			this.subscriptions[id].forEach(callback => callback(message));
		}
	}

	addTerminal(metadata: ActionMetadata) {
		const id = metadata.ActionId;

		const t = this.terminals.find(t => t.id === id);
		if(t) {
			t.Action = metadata.Action;
			t.Environment = metadata.Environment;
			t.Object = metadata.Object;
			t.Status = metadata.Status;

			t.StartedAt = dayjs(metadata.StartedAt);
			t.FinishedAt = dayjs(metadata.FinishedAt);

			if(!t.content) {
				t.content = metadata.Output || '';
			}
		} else {
			this.terminals.push({
				id,
				content: metadata.Output || '',
				Action: metadata.Action,
				Environment: metadata.Environment,
				Object: metadata.Object,

				FinishedAt: dayjs(metadata.FinishedAt),
				StartedAt: dayjs(metadata.StartedAt),
				Status: metadata.Status,
			})
		}
	}

	getTerminalIds () {
		return this.terminals.map(t => t.id);
	}

	filterTerminalsById (ids: string[]) {
		this.terminals = this.terminals.filter(t => ids.includes(t.id));
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
