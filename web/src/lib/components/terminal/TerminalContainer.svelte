<script lang="ts">
	import { Terminal as XTerm } from '@xterm/xterm';
	import { type Terminal, terminalStatus } from '$lib/terminal/TerminalStore.svelte';
	import { onMount } from 'svelte';
	import TerminalStore from '$lib/terminal/TerminalStore.svelte';
	import InvisibleButton from '$lib/components/fragments/InvisibleButton.svelte';
	import { fly } from 'svelte/transition';
	import TerminalInfo from '$lib/components/terminal/TerminalInfo.svelte';

	type Props = {
		content: string;
		id: string;
	}

	const xterm = new XTerm();
	let termObject: HTMLDivElement;

	let { content, id }: Props = $props();

	let terminal = $state<Terminal|null>(TerminalStore.getTerminal(id) ?? null);

	function onData(data: string) {
		xterm.write(data);
	}

	let cmdStatus = $derived(terminalStatus(terminal?.Status ?? -1));

	onMount(() => {
		xterm.options.theme = {
			background: '#151515', // this is surface-a0, xterm doesn't support CSS variables
		}
		xterm.open(termObject);
		xterm.resize(60, 12);
		xterm.write(content);

		TerminalStore.subscribe(id, {onData});

		return () => {
			TerminalStore.unsubscribe(id, {onData});
			xterm.dispose();
		};
	})

	let show = $state(false);

	$effect(() => {
		if (terminal) {
			terminal.DoNotDelete = show;
		}

		if(!show && terminal?.MarkedForDeletion) {
			TerminalStore.removeTerminal(terminal.id);
		}
	})

	let termHeight = $state(200);
	let termWidth = $state(500);
	let shortFormWidth = $state(120);
</script>

<div
	class="aside-container"
	class:expanded={show}

	style="
		--termHeight: {termHeight}px;
		--termWidth: {termWidth}px;
		--shortFormWidth: {shortFormWidth}px;
	"
>
	<div class="container"
		 transition:fly={{ delay: 200, duration: 200, x: 500 }}

		 class:success={cmdStatus === 'Success'}
		 class:pending={cmdStatus === 'Pending'}
		 class:running={cmdStatus === 'Running'}
		 class:failed={cmdStatus === 'Failed'}

		 class:toDelete={terminal?.MarkedForDeletion}
	>
		<div class="terminal-content"
			class:show={show}
		>
			<div id="terminal"
				 bind:this={termObject}
				 bind:clientHeight={termHeight}
				 bind:clientWidth={termWidth}
			></div>
		</div>
		<InvisibleButton
			onclick={() => show = !show}
		>
			{#if terminal}
				<TerminalInfo
					{terminal}
					{show}
					bind:shortFormWidth={shortFormWidth}
					{cmdStatus}
				/>
			{/if}
		</InvisibleButton>
	</div>
</div>

<style>
	.aside-container {
		width: calc(var(--shortFormWidth) + 1rem);
		transition: all 0.3s ease-in-out;
	}
	.aside-container.expanded {
        /* 1rem for padding, 4px for border */
        width: calc(var(--termWidth) + 1rem + 4px);
	}
	.container {
		background-color: var(--surface-tonal-a10);
        border-bottom-left-radius: 0.5rem;
		border-top-left-radius: 0.5rem;
		overflow: clip;
		border: 2px solid var(--surface-tonal-a10);
		transition: all 0.3s ease-in-out;
        width: max-content;
    }
    .container.pending {
        border: 2px solid var(--yellow-a10);
    }
	.container.success {
		border: 2px solid var(--green-a10);
	}
	.container.failed {
		border: 2px solid var(--red-a10);
	}
	.container.toDelete {
		background-color: var(--surface-a40);
		border: 2px solid var(--surface-a40);
		color: var(--dark-a0);
	}
	.terminal-content {
		background-color: var(--surface-a0);
        padding: 0;
        max-height: 0;
		overflow: hidden;
		transition: all 0.2s ease-in-out;
	}
	.terminal-content.show {
        padding: 0.5rem;
        max-height: calc(var(--termHeight) + 1rem); /* 1rem for padding */
	}
</style>