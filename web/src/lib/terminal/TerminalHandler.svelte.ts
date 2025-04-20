import type { Action } from 'svelte/action';
import { dev } from '$app/environment';

const terminalHandler: Action<HTMLElement> = () => {
	$effect(() => {
		async function main() {
			const url = dev ? 'localhost:4000' : window.location.host;

			const res = await fetch(window.location.protocol + "//" + url + "/api/actions")

			if (!res.ok) {
				console.error('Failed to fetch actions');
				return;
			}
		}

		main();
	});
}

export default terminalHandler;