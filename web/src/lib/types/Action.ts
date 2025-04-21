import { z } from 'zod';

export type ActionStatus = "Pending" | "Running" | "Success" | "Failed" | "Unknown";

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
});

const actionMessageBase = z.object({
	ActionId: z.string(),
	Type: z.string().startsWith("Action")
})

export const actionMessage = z.object({
	Message: z.string()
}).merge(actionMessageBase);

export const actionMessageMetadata = z.object({
	Metadata: actionMetadata
}).merge(actionMessageBase);

export const actionRemovedMessage = z.object({
	Removed: z.boolean()
}).merge(actionMessageBase);

export type ActionMetadata = z.infer<typeof actionMetadata>;

export const actionsList = z.object({
	Actions: z.record(z.string(), actionMetadata)
})
