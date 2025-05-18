<script lang="ts">
	import PopupStore from '$lib/stores/CurrentPopup.svelte';
	import PrettyButton from '$lib/components/fragments/PrettyButton.svelte';

	import Play from '~icons/ph/play';
	import ArrowsClockwise from '~icons/ph/arrows-clockwise';
	import Stop from '~icons/ph/stop';
	import Trash from '~icons/ph/trash';
	import ContainerAction from '$lib/websocket/actions/Container';
	import EnvStore from '$lib/stores/EnvStore.svelte';
	import Popup from '$lib/components/fragments/Popup/Popup.svelte';
	import PopupShowButton from '$lib/components/fragments/Popup/PopupShowButton.svelte';

	type Props = {
		id: string;
		name: string;
	};

	let { id, name }: Props = $props();

	function handleClick(action: string) {
		if (!EnvStore.name) {
			return;
		}
		ContainerAction(EnvStore.name, action, name);
		PopupStore.clear();
	}

	let popupShown = $derived(PopupStore.popup === id);
</script>

<div class="container" class:shown={popupShown}>
	<PopupShowButton {id} />
	{#if popupShown}
		<Popup {id}>
			<PrettyButton
				hoverBackground="var(--green-a20)"
				hoverColor="var(--dark-a0)"
				onclick={() => handleClick('start')}
			>
				<div class="icon-holder">
					<Play width="1.2rem" height="1.2rem" />
				</div>
				Start
			</PrettyButton>
			<PrettyButton
				hoverBackground="var(--yellow-a20)"
				hoverColor="var(--dark-a0)"
				onclick={() => handleClick('restart')}
			>
				<div class="icon-holder">
					<ArrowsClockwise width="1.2rem" height="1.2rem" />
				</div>
				Restart
			</PrettyButton>
			<PrettyButton
				hoverBackground="var(--red-a20)"
				hoverColor="var(--dark-a0)"
				onclick={() => handleClick('stop')}
			>
				<div class="icon-holder">
					<Stop width="1.2rem" height="1.2rem" />
				</div>
				Stop
			</PrettyButton>
			<PrettyButton
				hoverBackground="var(--red-a20)"
				hoverColor="var(--dark-a0)"
				onclick={() => handleClick('remove')}
			>
				<div class="icon-holder">
					<Trash width="1.2rem" height="1.2rem" />
				</div>
				Delete
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
