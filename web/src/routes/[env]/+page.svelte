<script lang="ts">
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>{data.env} - Shipyard</title>
</svelte:head>

<h2>
	Overview for {data.env}
</h2>
