<script lang="ts">
	import Funnel from '~icons/ph/funnel';
	import InvisibleButton from '$lib/components/fragments/InvisibleButton.svelte';
	import PopupStore from '$lib/stores/CurrentPopup.svelte';
	import Popup from '$lib/components/fragments/Popup/Popup.svelte';

	type Props = {
		id: string;
	};

	let { id }: Props = $props();

	let popupId = $derived('filter-popup-' + id);
	let popupShown = $derived(PopupStore.popup === popupId);
</script>

<div class="popup-container">
	<InvisibleButton
		--display="inline-flex"
		--width="max-content"
		class="show-button"
		onclick={() => PopupStore.toggle(popupId)}
	>
		<Funnel width="1.2rem" height="1.2rem" />
	</InvisibleButton>
	{#if popupShown}
		<Popup id={popupId}>WIP</Popup>
	{/if}
</div>

<style>
	.popup-container {
		display: flex;
		align-items: center;
		flex-shrink: 0;
		position: relative;
	}
</style>
