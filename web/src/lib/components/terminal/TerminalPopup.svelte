<script lang="ts">
	import TerminalStore from '$lib/terminal/TerminalStore.svelte';
	import TerminalContainer from '$lib/components/terminal/TerminalContainer.svelte';
</script>

<p>hehe</p>

{#each TerminalStore.terms as t (t.id)}
	<TerminalContainer content={t.content} id={t.id} />
{/each}