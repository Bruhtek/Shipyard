<script lang="ts">
	import PopupStore from '$lib/stores/CurrentPopup.svelte';
	import PrettyButton from '$lib/components/fragments/PrettyButton.svelte';

	import EnvStore from '$lib/stores/EnvStore.svelte';
	import Popup from '$lib/components/fragments/Popup/Popup.svelte';
	import PopupShowButton from '$lib/components/fragments/Popup/PopupShowButton.svelte';
	import type { ActionData } from '$lib/components/table/actionButtons/ActionData';

	type Props = {
		actions: ActionData[];
		id: string;
		name: string;
	};

	let { actions, id, name }: Props = $props();

	function handleClick(callback: ActionData['onClick']) {
		if (!EnvStore.name) {
			return;
		}
		callback(id, name);
		PopupStore.clear();
	}

	let popupShown = $derived(PopupStore.popup === id);
</script>

<div class="container" class:shown={popupShown}>
	<PopupShowButton {id} />
	{#if popupShown}
		<Popup {id}>
			{#each actions as action (action.text)}
				<PrettyButton
					hoverBackground={action.hoverBackground}
					hoverColor={action.hoverColor}
					onclick={() => handleClick(action.onClick)}
				>
					<div class="icon-holder">
						<action.icon width="1.2rem" height="1.2rem" />
					</div>
					{action.text}
				</PrettyButton>
			{/each}
		</Popup>
	{/if}
</div>

<style>
	.container {
		position: relative;
		width: 100%;
		height: 100%;
	}

	.container :global(svg) {
		transition: transform 0.2s ease-in-out;
	}
	.container.shown :global(.show-button svg) {
		transform: rotate(90deg);
	}
	.icon-holder {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 1.3rem;
		height: 1.3rem;
	}
</style>
