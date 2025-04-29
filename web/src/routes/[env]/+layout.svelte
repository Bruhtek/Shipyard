<script lang="ts">
	import type { LayoutProps } from './$types';
	import EnvStore from '$lib/stores/EnvStore.svelte';
	import { untrack } from 'svelte';

	let { data, children }: LayoutProps = $props();

	$effect(() => {
		if (data.env) {
			untrack(() => {
				EnvStore.name = data.env;
			});
		}
		return () => {
			EnvStore.clear();
		};
	});
</script>

{@render children()}
