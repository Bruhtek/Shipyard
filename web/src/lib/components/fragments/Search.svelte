<script lang="ts">
	import MagnifyingGlass from '~icons/ph/magnifying-glass';
	import { onMount, tick } from 'svelte';
	import { replaceState } from '$app/navigation';
	import { page } from '$app/state';

	type Props = {
		query: string;
	};

	let { query = $bindable('') }: Props = $props();

	onMount(() => {
		let url = new URL(window.location.href);
		if (url.searchParams.has('query')) {
			query = url.searchParams.get('query') || '';
		}
	});

	$effect(() => {
		if (query.trim() === '') {
			page.url.searchParams.delete('query');
		} else {
			page.url.searchParams.set('query', query);
		}

		async function replace() {
			await tick();
			replaceState(page.url, page.state);
		}

		replace();
	});
</script>

<div class="container">
	<input type="text" class="search" bind:value={query} />
	{#if query.trim() === ''}
		<div class="placeholder">
			<MagnifyingGlass width="1.4rem" height="1.4rem" />
			Search...
		</div>
	{/if}
</div>

<style>
	.container {
		display: flex;
		justify-content: center;
		align-items: center;
		position: relative;
		background: var(--surface-tonal-a20);
		border-radius: 0.5rem;
	}
	.search {
		display: block;
		width: 100%;
		height: 100%;
		padding: 0.5rem;
		margin: 0;
		background: transparent;
		border: none;
		color: var(--light-a0);
		border-radius: 0.5rem;
	}
	.placeholder {
		pointer-events: none;
		position: absolute;
		left: 0.5rem;
		top: 50%;
		transform: translateY(-50%);
		display: flex;
		gap: 0.5rem;
		opacity: 0.7;
	}
</style>
