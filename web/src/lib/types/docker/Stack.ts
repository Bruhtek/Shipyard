import { z } from 'zod';

export const TStack = z.object({
	ID: z.string().optional(),
	Name: z.string(),
	Status: z.string(),
	ConfigFiles: z.string()
});

export type Stack = z.infer<typeof TStack>;
export type StackWithID = Stack & { ID: string };

export const TStackResponse = z.object({
	Stacks: z.record(z.string(), TStack)
});
