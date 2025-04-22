import { z } from 'zod';

export enum ActionStatus {
	Unknown = -1,
	Pending = 0,
	Running = 1,
	Success = 2,
	Failed = 3,
}


export const actionMetadata = z.object({
	// command details
	ActionId: z.string().nonempty(),
	Environment: z.string().nonempty(),
	Object: z.string().nonempty(),
	Action: z.string().nonempty(),
	ObjectId: z.string(),

	StartedAt: z.string().datetime({ offset: true }),
	FinishedAt: z.string().datetime({ offset: true }),
	Status: z.number().int().min(ActionStatus.Unknown).max(ActionStatus.Failed),

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
