<script lang="ts">
	import { Terminal as XTerm } from '@xterm/xterm';
	import { type Terminal, terminalStatus } from '$lib/terminal/TerminalStore.svelte';
	import { onMount } from 'svelte';
	import TerminalStore from '$lib/terminal/TerminalStore.svelte';
	import CaretUp from '~icons/ph/caret-up';
	import InvisibleButton from '$lib/components/fragments/InvisibleButton.svelte';
	import { fly } from 'svelte/transition';

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
	function onMetadata(data: Terminal) {
		if(!terminal) {
			terminal = data;
		} else {
			Object.assign(terminal, data);
		}
	}

	let cmdStatus = $derived(terminalStatus(terminal?.Status ?? -1));

	onMount(() => {
		xterm.open(termObject);
		xterm.resize(60, 12);
		xterm.write(content);

		TerminalStore.subscribe(id, {onData, onMetadata});

		return () => {
			TerminalStore.unsubscribe(id, {onData, onMetadata});
			xterm.dispose();
		};
	})

	let show = $state(false);
	let termHeight = $state(200);
</script>

<div class="container"
	 transition:fly={{ delay: 200, duration: 200, x: 500 }}
	 class:success={cmdStatus === 'Success'}
	 class:pending={cmdStatus === 'Pending'}
	 class:running={cmdStatus === 'Running'}
	 class:failed={cmdStatus === 'Failed'}

	 style="--termHeight: {termHeight}px;"
>
	<div class="terminal-content"
		class:show={show}
	>
		<div id="terminal"
			 bind:this={termObject}
			 bind:clientHeight={termHeight}
		></div>
	</div>
	<div class="info">
		<div class="left">
			<InvisibleButton onclick={() => show = !show}>
				<span
					class="icon-container"
					class:rotated={show}
				>
					<CaretUp />
				</span>
			</InvisibleButton>
			<p>{cmdStatus}</p>
		</div>
		{#if terminal}
			<div class="right">
				<p>
					{terminal.Environment} | {terminal.Action} | {terminal.Object}
				</p>
			</div>
		{/if}
	</div>
</div>

<style>
	.container {
		width: max-content;
		background-color: var(--surface-tonal-a10);
        border-bottom-left-radius: 0.5rem;
		border-top-left-radius: 0.5rem;
		overflow: clip;
		border: 2px solid var(--surface-tonal-a10);
		transition: all 0.3s ease-in-out;
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
	.terminal-content {
		max-height: 0;
		overflow: hidden;
		transition: max-height 0.2s ease-in-out;
	}
	.terminal-content.show {
		max-height: var(--termHeight);
	}
	.info {
		height: 2rem;
		width: 100%;
		display: flex;
		justify-content: space-between;
		padding: 0.3rem;
	}
	.left, .right {
		display: flex;
		gap: 0.5rem;
	}

	.icon-container {
		transition: transform 0.2s ease-in-out;
		display: block;
	}
	.rotated {
		transform: rotate(180deg);
	}
</style>