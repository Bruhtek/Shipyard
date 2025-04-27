import { z } from 'zod';

export const TContainerState = z.enum([
	'created',
	'restarting',
	'running',
	'paused',
	'exited',
	'dead'
]);

export type ContainerState = z.infer<typeof TContainerState>;

export const TContainer = z.object({
	ID: z.string(),
	Image: z.string(),
	Labels: z.record(z.string(), z.string()),
	Name: z.string(),
	Names: z.array(z.string()),
	Ports: z.array(z.string()),
	Networks: z.array(z.string()),
	State: TContainerState,
	Status: z.string(),
	CreatedAt: z.string().datetime({ offset: true }),
	Command: z.string()
});

export type Container = z.infer<typeof TContainer>;

export const TContainerResponse = z.object({
	Containers: z.record(z.string(), TContainer)
});
