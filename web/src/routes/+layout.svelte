<script lang="ts">
	import '../app.css';
	import '@xterm/xterm/css/xterm.css';

	import WebSocketStatus from '$lib/components/WebSocketStatus.svelte';
	import websocketHandler from '$lib/websocket/WebsocketHandler.svelte.js';
	import terminalHandler from '$lib/terminal/TerminalHandler.svelte.js';
	import TerminalPopup from '$lib/components/terminal/TerminalPopup.svelte';
		import Sidebar from '$lib/components/sidebar/Sidebar.svelte';

	let { children } = $props();
</script>

<div use:websocketHandler></div>
<div use:terminalHandler></div>

<div class="page">
	<Sidebar />

	<WebSocketStatus />
	<TerminalPopup />

	<main class="main">
		<div class="content">
			{@render children()}
		</div>
	</main>
</div>

<style>
	.page {
		display: flex;
		min-height: 100vh;
	}

	.main {
		flex: 1;
		margin-left: 0.5rem;
	}
	.content {
		margin: 0 auto;
		padding: 1rem;
	}
</style>