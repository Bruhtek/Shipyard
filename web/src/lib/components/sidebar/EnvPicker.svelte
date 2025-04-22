<script lang="ts">
	import EnvStore from '$lib/stores/EnvStore.svelte';
	import EnvShow from '$lib/components/sidebar/EnvShow.svelte';

	let availableEnvs = $derived(EnvStore.availableEnvsList)
</script>

<div class="container">
	<h3 class="title">Environments</h3>
	{#each availableEnvs as env (env.Name)}
		<EnvShow envType={env.EnvType} name={env.Name} />
	{/each}
</div>

<style>
	.title {
		padding: 0.5rem;
	}
	.container {
		margin: 0.5rem 0;
		background-color: var(--surface-a0);
		border-radius: var(--border-radius);
	}
</style>