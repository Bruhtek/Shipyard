<script lang="ts">
	import { ContainerUpToDate } from '$lib/types/docker/Container';
	import CheckCircle from '~icons/ph/check-circle';
	import Question from '~icons/ph/question';
	import Spinner from '~icons/ph/spinner';
	import Cloud from '~icons/ph/cloud';
	import PushPin from '~icons/ph/push-pin';
	import type { Component } from 'svelte';
	import Badge from '$lib/components/fragments/Badge.svelte';

	type Props = {
		state: ContainerUpToDate;
	};

	let { state }: Props = $props();

	const StateToIcon: Record<ContainerUpToDate, Component> = {
		[ContainerUpToDate.UpToDate]: CheckCircle,
		[ContainerUpToDate.Unknown]: Question,
		[ContainerUpToDate.Pending]: Spinner,
		[ContainerUpToDate.UpdateAvailable]: Cloud,
		[ContainerUpToDate.Pinned]: PushPin
	};
	const StateToText: Record<ContainerUpToDate, string> = {
		[ContainerUpToDate.UpToDate]: 'Up to date',
		[ContainerUpToDate.Unknown]: 'Unknown',
		[ContainerUpToDate.Pending]: 'Checking...',
		[ContainerUpToDate.UpdateAvailable]: 'Update available',
		[ContainerUpToDate.Pinned]: 'Pinned'
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

{#if state === ContainerUpToDate.Pending || state === ContainerUpToDate.Unknown}
	{@render UpToDateStatus('var(--surface-tonal-a20)', undefined)}
{:else if state === ContainerUpToDate.UpdateAvailable}
	{@render UpToDateStatus('var(--yellow-a20)', 'var(--dark-a0)')}
{:else if state === ContainerUpToDate.UpToDate}
	{@render UpToDateStatus('var(--green-a20)', 'var(--dark-a0)')}
{:else if state === ContainerUpToDate.Pinned}
	{@render UpToDateStatus('var(--primary-a20)', 'var(--dark-a0)')}
{:else}
	{@render UpToDateStatus('var(--surface-tonal-a20)', undefined)}
{/if}

<style>
	.icon {
		display: inline-block;
		width: 1.4rem;
		height: 1.4rem;
		flex-shrink: 0;
	}
</style>
