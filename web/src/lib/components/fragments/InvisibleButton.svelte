<script lang="ts">

  	import type { Snippet } from 'svelte';
		import type { HTMLButtonAttributes } from 'svelte/elements';

	type Props = {
		children: Snippet<[]>;
	} & HTMLButtonAttributes;

	let { children, ...rest }: Props = $props();

</script>


<button
	{...rest}
	class:invisible-button={true}
>
	{@render children()}
</button>

<style>
	.invisible-button {
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
	}
</style>