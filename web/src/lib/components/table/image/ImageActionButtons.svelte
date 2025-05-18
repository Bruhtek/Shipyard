<script lang="ts">
	import PopupStore from '$lib/stores/CurrentPopup.svelte';
	import PrettyButton from '$lib/components/fragments/PrettyButton.svelte';

	import ArrowLineDown from '~icons/ph/arrow-line-down';
	import Trash from '~icons/ph/trash';
	import EnvStore from '$lib/stores/EnvStore.svelte';
	import ImageAction from '$lib/websocket/actions/Image';
	import Popup from '$lib/components/fragments/Popup/Popup.svelte';
	import PopupShowButton from '$lib/components/fragments/Popup/PopupShowButton.svelte';

	type Props = {
		id: string;
		repo: string;
		tag: string;
	};

	let { id, repo, tag }: Props = $props();

	let popupShown = $derived(PopupStore.popup === id);

	function handleClick(action: string) {
		if (!EnvStore.name) {
			return;
		}

		let objectId = id;
		if (action == 'pull') {
			objectId = repo + ':' + tag;
		}

		ImageAction(EnvStore.name, action, objectId);
		PopupStore.clear();
	}

	let updatePermitted = $derived(
		repo !== '' && tag !== '' && repo !== '<none>' && tag !== '<none>'
	);
</script>

<div class="container" class:shown={popupShown}>
	<PopupShowButton {id} />
	{#if popupShown}
		<Popup {id}>
			{#if updatePermitted}
				<PrettyButton
					hoverBackground="var(--primary-a20)"
					hoverColor="var(--light-a0)"
					onclick={() => handleClick('pull')}
				>
					<div class="icon-holder">
						<ArrowLineDown width="1.2rem" height="1.2rem" />
					</div>
					Update
				</PrettyButton>
			{/if}
			<PrettyButton
				hoverBackground="var(--red-a20)"
				hoverColor="var(--dark-a0)"
				onclick={() => handleClick('rm')}
			>
				<div class="icon-holder">
					<Trash width="1.2rem" height="1.2rem" />
				</div>
				Remove
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
