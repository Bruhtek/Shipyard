<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes } from 'svelte/elements';

	type Props = {
		children: Snippet<[]>;
		center?: boolean;
		spaceBetween?: boolean;
		background?: string;
		hoverBackground?: string;
		color?: string;
		hoverColor?: string;
	} & HTMLButtonAttributes;

	let {
		children,
		center,
		spaceBetween,
		background,
		color,
		hoverColor,
		hoverBackground,
		...rest
	}: Props = $props();

	let style = $derived.by(() => {
		let styles = '';
		if (background) {
			styles += `--bg: ${background};`;
		}
		if (color) {
			styles += `--col: ${color};`;
		}
		if (hoverBackground) {
			styles += `--hover-bg: ${hoverBackground};`;
		}
		if (hoverColor) {
			styles += `--hover-col: ${hoverColor};`;
		}
		return styles;
	});
</script>

<button {...rest} class:pretty-button={true} class:center class:spaceBetween {style}>
	{@render children()}
</button>

<style>
	.pretty-button {
		display: flex;
		width: 100%;
		gap: 0.5rem;
		padding: 0.5rem 1rem;
		align-items: center;
		background: var(--bg, none);
		color: var(--col, inherit);
		border: none;
		margin: 0;
		cursor: pointer;
		outline: none;
		font: inherit;
		text-align: inherit;
		text-decoration: none;
		-webkit-appearance: none;
		-moz-appearance: none;
		appearance: none;
		-webkit-tap-highlight-color: transparent;
		-webkit-user-select: none;
		-moz-user-select: none;
		user-select: none;
		transition:
			all 0.2s ease-in-out,
			opacity 0s linear;
	}
	.pretty-button:hover {
		background: var(--hover-bg, var(--bg, none));
		color: var(--hover-col, var(--col, inherit));
	}
	.pretty-button:active {
		opacity: 0.8;
	}
	.center {
		justify-content: center;
		align-content: center;
	}
	.spaceBetween {
		justify-content: space-between;
	}
</style>
