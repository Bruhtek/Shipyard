import { z } from 'zod';
import { URLPrefix } from '$lib';

type EnvStoreType = {
	name: string;
};

const AvailableEnvs = z.object({
	Environments: z.array(
		z.object({
			Name: z.string(),
			EnvType: z.string()
		})
	)
});

type AvailableEnvType = z.infer<typeof AvailableEnvs.shape.Environments.element>;

class EnvStoreClass {
	envObj = $state<EnvStoreType | null>(null);
	availableEnvs = $state<AvailableEnvType[]>([]);
	interval: number | null = null;

	constructor() {
		this.fetchEnvs();
	}

	async fetchEnvs() {
		const res = await fetch(URLPrefix + '/api/env');
		if (!res.ok) {
			throw new Error('Failed to fetch environments');
		}

		const data = await res.json();
		const parsed = AvailableEnvs.parse(data);

		this.availableEnvs = parsed.Environments;
	}

	get name(): string {
		return this.envObj?.name ?? '';
	}
	set name(name: string) {
		this.envObj = { name };
		this.startRequests();
	}
	get availableEnvsList(): AvailableEnvType[] {
		return this.availableEnvs;
	}

	clear() {
		this.envObj = null;
		if (this.interval) {
			clearInterval(this.interval);
		}
	}
	startRequests() {
		if (!this.envObj) return;

		const env = this.availableEnvs.find((env) => env.Name === this.envObj!.name);
		if (!env || env.EnvType !== 'remote') {
			return;
		}

		if (this.interval) {
			clearInterval(this.interval);
		}
		this.request();
		this.interval = setInterval(() => {
			this.request();
		}, 1000 * 60);
	}

	async request() {
		const res = await fetch(URLPrefix + `/api/env/${this.envObj?.name}/request`);
		if (!res.ok) {
			throw new Error('Failed to request environment');
		}
	}
}

const EnvStore = new EnvStoreClass();

export default EnvStore;
