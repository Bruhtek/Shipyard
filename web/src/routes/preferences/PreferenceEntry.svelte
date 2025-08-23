<script lang="ts">
	import UserPreferencesStore, {
		type UserPreferences,
		userPreferencesMetadata
	} from '$lib/stores/UserPreferences.svelte';

	type Props = {
		key: keyof UserPreferences;
		metadata: (typeof userPreferencesMetadata)[keyof UserPreferences];
	};

	const { key, metadata }: Props = $props();

	function handleInput(event: Event & { currentTarget: HTMLInputElement }) {
		const target = event.target as HTMLInputElement;
		let value: number | boolean;

		if (typeof metadata.default === 'boolean') {
			value = target.checked;
		} else {
			// it's a number
			value = parseFloat(target.value);
			if (isNaN(value)) return;
		}

		UserPreferencesStore.setPreference(key, value);
	}

	function handleChange(event: Event & { currentTarget: HTMLInputElement }) {
		const target = event.target as HTMLInputElement;
		let value: number | boolean;

		if (typeof metadata.default === 'boolean') {
			value = target.checked;
		} else {
			// it's a number
			value = parseFloat(target.value);
			if (isNaN(value)) return;
		}

		UserPreferencesStore.setPreference(key, value);
	}
</script>

<div>
	<h2>{metadata.label}</h2>
	{#if metadata.tip}
		<p>{metadata.tip}</p>
	{/if}
	{#if typeof metadata.default === 'boolean'}
		<!-- This double !! makes svelte stop complaining about a possibility of a number, despite the default being boolean -->
		<input type="checkbox" checked={!!UserPreferencesStore.p[key]} onchange={handleChange} />
	{:else if typeof metadata.default === 'number'}
		<input
			type="number"
			value={UserPreferencesStore.p[key]}
			min={metadata.limits?.min || undefined}
			max={metadata.limits?.max || undefined}
			oninput={handleInput}
		/>
	{/if}
</div>

<style>
</style>
