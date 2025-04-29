<script lang="ts">
	import type { LayoutProps } from './$types';
	import EnvStore from '$lib/stores/EnvStore.svelte';

	let { data, children }: LayoutProps = $props();

	$effect(() => {
		if (data.env) {
			EnvStore.name = data.env;

			return () => {
				EnvStore.clear();
			};
		}
	});
</script>

{@render children()}
