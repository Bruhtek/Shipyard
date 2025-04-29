<script lang="ts">
	import type { ContainerState } from '$lib/types/docker/Container';
	import { toProperCase } from '$lib';
	import MapPinSimple from '~icons/ph/map-pin-simple';
	import ArrowsClockwise from '~icons/ph/arrows-clockwise';
	import CheckCircle from '~icons/ph/check-circle';
	import Pause from '~icons/ph/pause';
	import DoorOpen from '~icons/ph/door-open';
	import Skull from '~icons/ph/skull';
	import type { Component } from 'svelte';
	import Badge from '$lib/components/fragments/Badge.svelte';

	type Props = {
		state: ContainerState;
	};

	let { state }: Props = $props();

	const StateToIcon: Record<ContainerState, Component> = {
		created: MapPinSimple,
		restarting: ArrowsClockwise,
		running: CheckCircle,
		paused: Pause,
		exited: DoorOpen,
		dead: Skull
	};

	const Icon = $derived(StateToIcon[state]);
</script>

{#snippet ContainerStatus(bg: string, color: string | undefined)}
	<Badge background={bg} {color}>
		<span class="icon">
			<Icon width="1.4rem" height="1.4rem" />
		</span>
		{toProperCase(state)}
	</Badge>
{/snippet}

{#if state === 'created' || state === 'exited'}
	{@render ContainerStatus('var(--surface-tonal-a20)', undefined)}
{:else if state === 'restarting' || state === 'paused'}
	{@render ContainerStatus('var(--yellow-a20)', 'var(--dark-a0)')}
{:else if state === 'running'}
	{@render ContainerStatus('var(--green-a20)', 'var(--dark-a0)')}
{:else if state === 'dead'}
	{@render ContainerStatus('var(--red-a20)', 'var(--dark-a0)')}
{/if}

<style>
	.icon {
		display: inline-block;
		width: 1.4rem;
		height: 1.4rem;
	}
</style>
