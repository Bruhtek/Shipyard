<script lang="ts">
	import { Terminal as XTerm } from '@xterm/xterm';
	import type { Terminal } from '$lib/terminal/TerminalStore.svelte';
	import { onMount } from 'svelte';
	import TerminalStore from '$lib/terminal/TerminalStore.svelte';

	type Props = {
		content: string;
		id: string;
		terminal: Terminal;
	}

	const xterm = new XTerm();
	let termObject: HTMLDivElement;

	let { content, id, terminal }: Props = $props();

	function onData(data: string) {
		xterm.write(data);
	}

	onMount(() => {
		xterm.open(termObject);
		xterm.resize(60, 12);
		xterm.write(content);

		TerminalStore.subscribe(id, onData);

		return () => {
			TerminalStore.unsubscribe(id, onData);
			xterm.dispose();
		};
	})
</script>

<div class="container">
	<div class="metadata">
		<p>Id: {terminal.id}</p>
		<p>Environment: {terminal.Environment}</p>
		<p>Action: {terminal.Action}</p>
		<p>Object: {terminal.Object}</p>
	</div>
	<div id="terminal" bind:this={termObject}></div>
</div>

<style>
	.container {
		width: max-content;
	}
</style>