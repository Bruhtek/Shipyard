<script lang="ts">
	type Props = {
		length?: number;
	};

	let { length = 10 }: Props = $props();

	let text = $derived.by(() => {
		const char = ' ';
		return char.repeat(length);
	});
</script>

<span class="shadow">{text}</span>

<style>
	.shadow {
		display: inline-block;
		width: 100%;
		height: 1.3rem;
		margin-bottom: 0.2rem;
		background: var(--surface-a20);
		border-radius: var(--border-radius);
		opacity: 0.7;
		font-family: monospace;
		white-space: pre;
		animation: pulse 2.5s infinite;
	}

	@keyframes pulse {
		0% {
			background: var(--surface-a20);
		}
		25% {
			background: var(--surface-a10);
		}
		50% {
			background: var(--surface-a20);
		}
	}
</style>
