<script lang="ts">
	import InvisibleButton from '$lib/components/fragments/InvisibleButton.svelte';
	import DotsThree from '~icons/ph/dots-three';
	import PopupStore from '$lib/stores/CurrentPopup.svelte';
	import PrettyButton from '$lib/components/fragments/PrettyButton.svelte';

	import Play from '~icons/ph/play';
	import ArrowsClockwise from '~icons/ph/arrows-clockwise';
	import Stop from '~icons/ph/stop';
	import Trash from '~icons/ph/trash';
	import ContainerAction from '$lib/websocket/actions/Container';
	import EnvStore from '$lib/stores/EnvStore.svelte';

	type Props = {
		id: string;
		name: string;
	};

	let { id, name }: Props = $props();

	let popupShown = $derived(PopupStore.popup === id);

	function handleClick(action: string) {
		if (!EnvStore.name) {
			return;
		}
		ContainerAction(EnvStore.name, action, name);
		PopupStore.clear();
	}
</script>

<div class="container" class:shown={popupShown}>
	<InvisibleButton center={true} onclick={() => PopupStore.toggle(id)} class="show-button">
		<DotsThree width="1.5rem" height="1.5rem" />
	</InvisibleButton>
	<div class="popup">
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
	</div>
</div>

<style>
	.container {
		position: relative;
	}
	.popup {
		display: none;
		position: absolute;
		background-color: var(--surface-tonal-a20);
		z-index: 10;
		right: 0;
		top: 100%;
		border-radius: var(--border-radius);
		overflow: hidden;
	}
	.container.shown .popup {
		display: block;
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
