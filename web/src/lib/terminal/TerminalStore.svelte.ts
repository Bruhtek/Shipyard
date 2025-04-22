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

	DoNotDelete?: boolean;
	MarkedForDeletion?: boolean;
}

type Subscription = {
	onData: (message: string) => void;
}

class terminalStore {
	// essentially just map Terminal(action) id to output
	terminals: Terminal[] = $state<Terminal[]>([]);
	subscriptions: { [key: string]: Subscription[] } = {};

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
			this.subscriptions[id].forEach(sub => sub.onData(message));
		}
	}

	addTerminal(metadata: ActionMetadata) {
		const id = metadata.ActionId;

		let t = this.terminals.find(t => t.id === id);
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
			t = {
				id,
				content: metadata.Output || '',
				Action: metadata.Action,
				Environment: metadata.Environment,
				Object: metadata.Object,

				FinishedAt: dayjs(metadata.FinishedAt),
				StartedAt: dayjs(metadata.StartedAt),
				Status: metadata.Status,
			}
			this.terminals.push(t)
		}
	}

	getTerminalIds () {
		return this.terminals.map(t => t.id);
	}

	removeTerminal(id: string) {
		const index = this.terminals.findIndex(t => t.id === id);
		if (index !== -1) {
			if(this.terminals[index].DoNotDelete) {
				this.terminals[index].MarkedForDeletion = true;
				return;
			}

			this.terminals.splice(index, 1);
		}

		if (this.subscriptions[id]) {
			delete this.subscriptions[id];
		}
	}

	filterTerminalsById (ids: string[]) {
		this.terminals = this.terminals.filter(t => ids.includes(t.id) || t.DoNotDelete);

		const doNotDeletes = this.terminals.filter(t => !ids.includes(t.id) && t.DoNotDelete);
		for (const t of doNotDeletes) {
			t.MarkedForDeletion = true;
		}

		const subscriptionIds = Object.keys(this.subscriptions);
		const filteredSubscriptions = subscriptionIds.filter(
			id => !ids.includes(id) &&
				!this.terminals.find(t => t.id === id)
		);

		filteredSubscriptions.forEach(id => {
			this.subscriptions[id] = [];
		})
	}

	subscribe (id: string, subscription: Subscription) {
		if (!this.subscriptions[id]) {
			this.subscriptions[id] = [];
		}
		this.subscriptions[id].push(subscription);
	}

	unsubscribe (id: string, subscription: Subscription) {
		if (this.subscriptions[id]) {
			this.subscriptions[id] = this.subscriptions[id].filter(sub => sub !== subscription);
		}
	}

	get terms () {
		return this.terminals;
	}

	getTerminal(id: string) {
		return this.terminals.find(t => t.id === id);
	}
}

const TerminalStore = new terminalStore();
export default TerminalStore;
