<script lang="ts">
	import Hourglass from '~icons/ph/hourglass';
	import Spinner from '~icons/ph/spinner';
	import Check from '~icons/ph/check';
	import WarningOctagon from '~icons/ph/warning-octagon';
	import QuestionMark from '~icons/ph/question-mark';

	import type { Component } from 'svelte';
	import { ActionStatus } from '$lib/types/Action';

	type Props = {
		status: ActionStatus;
	}

	const statusToIcon: { [key in ActionStatus]: Component} = {
		[ActionStatus.Pending]: Hourglass,
		[ActionStatus.Running]: Spinner,
		[ActionStatus.Success]: Check,
		[ActionStatus.Failed]: WarningOctagon,
		[ActionStatus.Unknown]: QuestionMark
	};

	const { status }: Props = $props();

	const Icon = $derived(statusToIcon[status])
</script>

<Icon height="1.5rem" width="1.5rem" />