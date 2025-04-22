import type { LayoutLoad } from './$types';

export const load: LayoutLoad = ({ params }) => {
	return {
		env: params.env,
	}
}