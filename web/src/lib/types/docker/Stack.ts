import { z } from 'zod';
import { ContainerUpToDate, TContainer } from '$lib/types/docker/Container';
import { TNetwork } from '$lib/types/docker/Network';

export const TStack = z.object({
	ID: z.string().optional(),
	Name: z.string(),
	Status: z.string(),
	ConfigFiles: z.string(),

	Containers: z.array(TContainer),
	Networks: z.array(TNetwork),

	UpToDate: z.nativeEnum(ContainerUpToDate).default(ContainerUpToDate.Pending)
});

export type Stack = z.infer<typeof TStack>;
export type StackWithID = Stack & { ID: string };

export const TStackResponse = z.object({
	Stacks: z.record(z.string(), TStack)
});
