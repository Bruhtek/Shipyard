<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes } from 'svelte/elements';

	type Props = {
		children: Snippet<[]>;
		center?: boolean;
		spaceBetween?: boolean;
	} & HTMLButtonAttributes;

	let { children, center, spaceBetween, ...rest }: Props = $props();
</script>

<button {...rest} class:invisible-button={true} class:center class:spaceBetween>
	{@render children()}
</button>

<style>
	.invisible-button {
		width: var(--width, 100%);
		display: var(--display, flex);
		gap: 0.5rem;
		background: none;
		border: none;
		padding: 0;
		margin: 0;
		cursor: pointer;
		outline: none;
		color: inherit;
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
		align-content: center;
	}
	.center {
		justify-content: center;
	}
	.spaceBetween {
		justify-content: space-between;
	}
</style>
