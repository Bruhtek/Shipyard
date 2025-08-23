<script lang="ts">
	import { userPreferencesMetadata } from '$lib/stores/UserPreferences.svelte';
	import PreferenceEntry from './PreferenceEntry.svelte';

	type Entries<T> = { [K in keyof T]: [K, T[K]] }[keyof T][];
	function entriesOf<T extends object>(obj: T): Entries<T> {
		return Object.entries(obj) as Entries<T>;
	}
</script>

<h1>Preferences</h1>

<div class="content">
	{#each entriesOf(userPreferencesMetadata) as [key, metadata] (key)}
		<PreferenceEntry {key} {metadata} />
	{/each}
</div>

<style>
	.content {
		margin-top: 1rem;
	}
</style>
