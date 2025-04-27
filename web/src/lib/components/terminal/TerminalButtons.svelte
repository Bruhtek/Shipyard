<script lang="ts">
	import { type Terminal } from '$lib/terminal/TerminalStore.svelte';
	import InvisibleButton from '$lib/components/fragments/InvisibleButton.svelte';

	import StopCircle from '~icons/ph/stop-circle';
	import X from '~icons/ph/x';
	import ArrowsClockwise from '~icons/ph/arrows-clockwise';
	import { ActionStatus } from '$lib/types/Action';
	import TerminalRequests from '$lib/terminal/TerminalRequests';

	type Props = {
		terminal: Terminal;
		show: boolean;
	};

	let { terminal, show = $bindable() }: Props = $props();

	function Dismiss() {
		TerminalRequests.DismissTerminal(terminal);
		show = false;
	}
	function Stop() {
		TerminalRequests.DismissTerminal(terminal);
	}
	function Retry() {
		TerminalRequests.RetryTerminal(terminal);
	}
</script>

<div class="control-buttons" class:show={show && !terminal.MarkedForDeletion}>
	{#if terminal.Status === ActionStatus.Running}
		<InvisibleButton center={true} class="btn btn-stop" onclick={Stop}>
			<StopCircle width="1.5rem" height="1.5rem" />
		</InvisibleButton>
	{/if}
	{#if terminal.Status === ActionStatus.Failed || terminal.Status === ActionStatus.Success}
		<InvisibleButton center={true} class="btn btn-dismiss" onclick={Dismiss}>
			<X width="1.5rem" height="1.5rem" />
		</InvisibleButton>
	{/if}
	{#if terminal.Status === ActionStatus.Failed}
		<InvisibleButton center={true} class="btn btn-retry" onclick={Retry}>
			<ArrowsClockwise width="1.5rem" height="1.5rem" />
		</InvisibleButton>
	{/if}
</div>

<style>
	.control-buttons {
		background-color: var(--surface-tonal-a10);
		position: absolute;
		width: 2rem;
		top: 0;
		left: 0;
		transition: all 0.2s ease-in-out 0.2s;
		display: none;
		opacity: 0;
		flex-direction: column;
		border-top-left-radius: 0.5rem;
		border-bottom-left-radius: 0.5rem;
		overflow: clip;
	}
	.control-buttons.show {
		display: flex;
		top: 10px;
		opacity: 1;
		left: -2rem;
	}

	.control-buttons :global(.btn) {
		padding: 0.25rem 0;
		transition: all 0.1s ease-in-out;
	}
	.control-buttons :global(.btn-stop:hover) {
		background-color: var(--red-a10);
	}
	.control-buttons :global(.btn-retry:hover) {
		background-color: var(--yellow-a10);
	}
	.control-buttons :global(.btn-dismiss:hover) {
		background-color: var(--surface-a40);
		color: var(--dark-a0);
	}
</style>
