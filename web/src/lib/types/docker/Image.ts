import { z } from 'zod';

export const TImage = z.object({
	ID: z.string(),
	Repository: z.string(),
	Tag: z.string(),
	Size: z.number(),
	CreatedAt: z.string().datetime({ offset: true }),

	Used: z.boolean(),
	RepoDigests: z.array(z.string())
});

export type Image = z.infer<typeof TImage>;

export const TImageResponse = z.object({
	Images: z.record(z.string(), TImage)
});
