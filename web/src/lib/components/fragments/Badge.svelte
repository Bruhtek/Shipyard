<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		children: Snippet<[]>;
		background?: string;
		color?: string;
	};

	let { children, background, color }: Props = $props();

	let styles = $derived.by(() => {
		let styles = '';
		if (background) {
			styles += `--background: ${background};`;
		}
		if (color) {
			styles += `--color: ${color};`;
		}
		return styles;
	});
</script>

<div class="badge" style={styles}>
	{@render children()}
</div>

<style>
	.badge {
		display: inline-flex;
		align-items: center;
		gap: 0.2rem;
		color: var(--color, unset);
		background-color: var(--background, var(--surface-tonal-a20));
		border-radius: 1.5rem;
		padding: 0.1rem 0.5rem 0.1rem 0.4rem;
	}
</style>
