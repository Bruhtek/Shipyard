import { z } from 'zod';

export const actionMetadata = z.object({
	// command details
	ActionId: z.string().nonempty(),
	Environment: z.string().nonempty(),
	Object: z.string().nonempty(),
	Action: z.string().nonempty(),
	ObjectId: z.string(),

	StartedAt: z.string().datetime({ offset: true }),
	FinishedAt: z.string().datetime({ offset: true }),
	Status: z.number().int().min(0).max(3),

	Output: z.string().optional(),
})

export const actionMessage = z.object({
	ActionId: z.string(),
	Message: z.string()
});
export const actionMessageMetadata = z.object({
	ActionId: z.string(),
	Metadata: actionMetadata
});
export type ActionMetadata = z.infer<typeof actionMetadata>;

export const actionsList = z.object({
	Actions: z.record(z.string(), actionMetadata)
})
