<script lang="ts">
	import PopupStore from '$lib/stores/CurrentPopup.svelte';
	import PrettyButton from '$lib/components/fragments/PrettyButton.svelte';

	import Trash from '~icons/ph/trash';
	import Funnel from '~icons/ph/funnel';
	import EnvStore from '$lib/stores/EnvStore.svelte';
	import Popup from '$lib/components/fragments/Popup/Popup.svelte';
	import PopupShowButton from '$lib/components/fragments/Popup/PopupShowButton.svelte';
	import NetworkAction from '$lib/websocket/actions/Network';
	import { goto } from '$app/navigation';
	import ImageAction from '$lib/websocket/actions/Image';

	type Props = {
		id: string;
		name: string;
	};

	let { id, name }: Props = $props();

	let popupShown = $derived(PopupStore.popup === id);

	const redirectToFilteredContainers = () => {
		const url = new URL(window.location.href);
		url.pathname = `/${EnvStore.name}/containers`;
		url.searchParams.set('query', `network:${name}`);
		goto(url.toString());
	};

	const removeNetwork = () => {
		if (!EnvStore.name) {
			return;
		}

		NetworkAction(EnvStore.name, 'remove', id);
		PopupStore.clear();
	};
</script>

<div class="container" class:shown={popupShown}>
	<PopupShowButton {id} />
	{#if popupShown}
		<Popup {id}>
			<PrettyButton
				hoverBackground="var(--red-a20)"
				hoverColor="var(--dark-a0)"
				onclick={removeNetwork}
			>
				<div class="icon-holder">
					<Trash width="1.2rem" height="1.2rem" />
				</div>
				Remove
			</PrettyButton>
			<PrettyButton
				hoverBackground="var(--primary-a20)"
				hoverColor="var(--dark-a0)"
				onclick={redirectToFilteredContainers}
			>
				<div class="icon-holder">
					<Funnel width="1.2rem" height="1.2rem" />
				</div>
				View Containers
			</PrettyButton>
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
