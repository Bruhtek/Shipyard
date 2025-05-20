import { z } from 'zod';
import { TContainer } from '$lib/types/docker/Container';

export const TNetwork = z.object({
	CreatedAt: z.string().datetime({ offset: true }),
	ID: z.string(),
	Driver: z.enum(['bridge', 'host', 'overlay', 'ipvlan', 'macvlan', 'none', 'null']),
	Internal: z.boolean(),
	IPv6: z.boolean(),
	Name: z.string(),
	Scope: z.enum(['local', 'swarm', 'global']),
	Containers: z.array(TContainer)
});

export type Network = z.infer<typeof TNetwork>;

export const TNetworkResponse = z.object({
	Networks: z.record(z.string(), TNetwork)
});
