<script lang="ts">
	import TerminalStore from '$lib/terminal/TerminalStore.svelte';
	import TerminalContainer from '$lib/components/terminal/TerminalContainer.svelte';

	let terminals = $derived(TerminalStore.terms.toSorted((a, b) => {
		if (a.StartedAt < b.StartedAt) {
			return -1;
		} else if (a.StartedAt > b.StartedAt) {
			return 1;
		}
		return 0;
	}));
</script>

<div class="container">
	{#each terminals as t (t.id)}
		<TerminalContainer
			content={t.content}
			id={t.id}
		/>
	{/each}
</div>

<style>
	.container {
		position: absolute;
		bottom: 2rem;
		right: 0;
	}
</style>