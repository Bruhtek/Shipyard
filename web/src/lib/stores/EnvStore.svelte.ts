import { z } from 'zod';
import { URLPrefix } from '$lib';

type EnvStoreType = {
	name: string;
}

const AvailableEnvs = z.object({
	Environments: z.array(z.object({
		Name: z.string(),
		EnvType: z.string(),
	}))
})

type AvailableEnvType = z.infer<typeof AvailableEnvs.shape.Environments.element>;

class EnvStoreClass {
	envObj = $state<EnvStoreType|null>(null);
	availableEnvs = $state<AvailableEnvType[]>([]);

	constructor() {
		this.fetchEnvs()
	}

	async fetchEnvs() {
		const res = await fetch(URLPrefix + "/api/env")
		if (!res.ok) {
			throw new Error("Failed to fetch environments");
		}

		const data = await res.json();
		const parsed = AvailableEnvs.parse(data);

		this.availableEnvs = parsed.Environments
	}

	get name(): string {
		return this.envObj?.name ?? "";
	}
	set name(name: string) {
		this.envObj = { name };
	}
	get availableEnvsList(): AvailableEnvType[] {
		return this.availableEnvs;
	}

	clear() {
		this.envObj = null;
	}

}

const EnvStore = new EnvStoreClass();

export default EnvStore;

