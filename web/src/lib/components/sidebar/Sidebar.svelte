<script lang="ts">
	import EnvPicker from '$lib/components/sidebar/EnvPicker.svelte';
	import EnvStore from '$lib/stores/EnvStore.svelte';
</script>

<div class="sidebar">
	<a href="/">
		<h2>Shipyard</h2>
	</a>

	<EnvPicker />
</div>

<style>
	.sidebar {
		height: 100vh;
		background-color: var(--surface-tonal-a10);
		padding: 1rem;
	}
</style>