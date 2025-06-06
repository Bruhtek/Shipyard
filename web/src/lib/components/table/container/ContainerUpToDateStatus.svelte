<script lang="ts">
	import { ContainerUpToDate } from '$lib/types/docker/Container';
	import CheckCircle from '~icons/ph/check-circle';
	import Question from '~icons/ph/question';
	import Spinner from '~icons/ph/spinner';
	import Cloud from '~icons/ph/cloud';
	import type { Component } from 'svelte';
	import Badge from '$lib/components/fragments/Badge.svelte';

	type Props = {
		state: ContainerUpToDate;
	};

	let { state }: Props = $props();

	const StateToIcon: Record<ContainerUpToDate, Component> = {
		[ContainerUpToDate.UpToDate]: CheckCircle,
		[ContainerUpToDate.Error]: Question,
		[ContainerUpToDate.Unknown]: Spinner,
		[ContainerUpToDate.UpdateAvailable]: Cloud
	};
	const StateToText: Record<ContainerUpToDate, string> = {
		[ContainerUpToDate.UpToDate]: 'Up to date',
		[ContainerUpToDate.Error]: 'Unknown',
		[ContainerUpToDate.Unknown]: 'Checking...',
		[ContainerUpToDate.UpdateAvailable]: 'Update available'
	};

	const Icon = $derived(StateToIcon[state]);
</script>

{#snippet UpToDateStatus(bg: string, color: string | undefined)}
	<Badge background={bg} {color}>
		<span class="icon">
			<Icon width="1.4rem" height="1.4rem" />
		</span>
		{StateToText[state]}
	</Badge>
{/snippet}

{#if state === ContainerUpToDate.Unknown || state === ContainerUpToDate.Error}
	{@render UpToDateStatus('var(--surface-tonal-a20)', undefined)}
{:else if state === ContainerUpToDate.UpdateAvailable}
	{@render UpToDateStatus('var(--yellow-a20)', 'var(--dark-a0)')}
{:else if state === ContainerUpToDate.UpToDate}
	{@render UpToDateStatus('var(--green-a20)', 'var(--dark-a0)')}
{/if}

<style>
	.icon {
		display: inline-block;
		width: 1.4rem;
		height: 1.4rem;
	}
</style>
