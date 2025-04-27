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

<span
	class="container-state"
	class:created={state === 'created'}
	class:restarting={state === 'restarting'}
	class:running={state === 'running'}
	class:paused={state === 'paused'}
	class:exited={state === 'exited'}
	class:dead={state === 'dead'}
>
	<Icon width="1.4rem" height="1.4rem" />
	{toProperCase(state)}
</span>

<style>
	.container-state {
		display: inline-flex;
		align-items: center;
		gap: 0.2rem;
		background-color: var(--surface-tonal-a20);
		border-radius: 1rem;
		padding: 0.1rem 0.4rem 0.1rem 0.3rem;
	}
	.created,
	.exited {
		background-color: var(--surface-tonal-a20);
	}
	.restarting,
	.paused {
		background-color: var(--yellow-a20);
		color: var(--dark-a0);
	}
	.running {
		background-color: var(--green-a20);
		color: var(--dark-a0);
	}
	.dead {
		background-color: var(--red-a20);
		color: var(--dark-a0);
	}
</style>
