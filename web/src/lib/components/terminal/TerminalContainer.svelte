<script lang="ts">
	import { Terminal } from '@xterm/xterm';
	import { onMount } from 'svelte';
		import TerminalStore from '$lib/terminal/TerminalStore.svelte';

	type Props = {
		content: string;
		id: string;
	}

	const term = new Terminal();
	let termObject: HTMLDivElement;

	let { content, id }: Props = $props();

	function onData(data: string) {
		term.write(data);
	}

	onMount(() => {
		term.open(termObject);
		term.resize(80, 12);
		term.write(content);

		TerminalStore.subscribe(id, onData);

		return () => {
			TerminalStore.unsubscribe(id, onData);
			term.dispose();
		};
	})
</script>

<div class="container">
	<div id="terminal" bind:this={termObject}></div>
</div>

<style>
	.container {
		width: max-content;
	}
</style>