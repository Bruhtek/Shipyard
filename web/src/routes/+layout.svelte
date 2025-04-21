<script lang="ts">
	import '../app.css';
	import '@xterm/xterm/css/xterm.css';

	import WebSocketStatus from '$lib/components/WebSocketStatus.svelte';
	import websocketHandler from '$lib/websocket/WebsocketHandler.svelte.js';
	import terminalHandler from '$lib/terminal/TerminalHandler.svelte.js';
	import TerminalPopup from '$lib/components/terminal/TerminalPopup.svelte';

	let { children } = $props();
</script>

<div use:websocketHandler></div>
<div use:terminalHandler></div>


<div>
	<main class="container mx-auto">
		<WebSocketStatus />

		<TerminalPopup />

		{@render children()}
	</main>
</div>