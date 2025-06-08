<script lang="ts">
	import type { Snippet } from 'svelte';
	import InvisibleButton from '$lib/components/fragments/InvisibleButton.svelte';
	import CaretDown from '~icons/ph/caret-down';
	import CarotUp from '~icons/ph/caret-up';

	type Props = {
		children: Snippet<[]>;
		current: string;
		currentDirection: 'asc' | 'desc';
		sortByKey: string;
	};

	let {
		children,
		current = $bindable(),
		currentDirection = $bindable('asc'),
		sortByKey
	}: Props = $props();

	function handleClick() {
		if (current === sortByKey) {
			currentDirection = currentDirection === 'asc' ? 'desc' : 'asc';
		} else {
			current = sortByKey;
			currentDirection = 'asc';
		}
	}
</script>

<div class="sorting-button" class:active={current === sortByKey}>
	<InvisibleButton onclick={handleClick} spaceBetween={true}>
		{@render children()}
		{#if current === sortByKey}
			{#if currentDirection === 'asc'}
				<CaretDown />
			{:else}
				<CarotUp />
			{/if}
		{:else}
			<CaretDown />
		{/if}
	</InvisibleButton>
</div>

<style>
	.sorting-button {
		transition: all 0.2s ease;
		width: 100%;
		flex-shrink: 1;
	}
	.sorting-button:not(.active):not(:hover) {
		opacity: 0.7;
	}
</style>
