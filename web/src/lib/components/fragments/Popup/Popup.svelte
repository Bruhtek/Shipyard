<script lang="ts">
	import { onMount, type Snippet } from 'svelte';
	import PopupStore from '$lib/stores/CurrentPopup.svelte.js';

	type Props = {
		id: string;
		children: Snippet<[]>;
	};

	let { id, children }: Props = $props();

	onMount(() => {
		const clearPopup = (e: MouseEvent) => {
			if (PopupStore.popup === id) {
				const target = e.target as HTMLElement;
				const popup = document.querySelector(`.popup[data-id="${id}"]`);
				if (popup && !popup.contains(target) && !target.closest('.show-button')) {
					PopupStore.clear();
				}
			}
		};
		document.addEventListener('click', clearPopup);

		return () => {
			document.removeEventListener('click', clearPopup);
		};
	});
</script>

<div class="popup" data-id={id}>
	{@render children()}
</div>

<style>
	.popup {
		display: block;
		position: absolute;
		background-color: var(--surface-tonal-a20);
		z-index: 10;
		right: 0;
		top: 100%;
		border-radius: var(--border-radius);
		overflow: hidden;
		box-shadow: 2px 2px 15px var(--dark-a0);
	}
</style>
