<script lang="ts">
	import Search from '$lib/components/fragments/Search.svelte';

	type Props = {
		title: string;
		query: string;
	};

	let { title, query = $bindable('') }: Props = $props();
</script>

<div class="table-header">
	<h3>{title}</h3>

	<Search bind:query />
</div>

<style>
	.table-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0.5rem 1rem;
		position: sticky;
		height: 3.5rem;
		top: 0;
		z-index: 20;
		background-color: var(--surface-tonal-a0);
	}
</style>
