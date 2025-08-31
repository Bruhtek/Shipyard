<script lang="ts">
	import type { ActionData } from '$lib/components/table/actionButtons/ActionData';
	import EnvStore from '$lib/stores/EnvStore.svelte';
	import PrettyButton from '$lib/components/fragments/PrettyButton.svelte';
	import X from '~icons/ph/x';
	import SelectionAll from '~icons/ph/selection-all';

	type Props = {
		actions: ActionData[];
		ids: string[];
		names: string[];
		clear: () => void;
		selectAll: () => void;
	};

	let { actions, ids, names, clear, selectAll }: Props = $props();

	function handleClick(
		callback: ActionData['onClick'],
		multiCallback: ActionData['onMultiClick']
	) {
		if (!EnvStore.name) {
			return;
		}
		if (multiCallback && ids.length > 1) {
			multiCallback(ids, names);
			clear();
			return;
		}
		for (let i = 0; i < ids.length; i++) {
			const id = ids[i];
			const name = names[i];
			callback(id, name);
		}
		clear();
	}
</script>

<div class="multi-action-container">
	<PrettyButton hoverBackground="var(--surface-a0)" hoverColor="var(--light-a0)" onclick={clear}>
		<div class="icon-holder">
			<X width="1.2rem" height="1.2rem" />
		</div>
	</PrettyButton>
	<PrettyButton
		hoverBackground="var(--surface-a0)"
		hoverColor="var(--light-a0)"
		onclick={selectAll}
	>
		<div class="icon-holder">
			<SelectionAll width="1.2rem" height="1.2rem" />
		</div>
	</PrettyButton>
	{#each actions as action (action.text)}
		{#if !action.singleOnly}
			<PrettyButton
				hoverBackground={action.hoverBackground}
				hoverColor={action.hoverColor}
				onclick={() => handleClick(action.onClick, action.onMultiClick)}
			>
				<div class="icon-holder">
					<action.icon width="1.2rem" height="1.2rem" />
				</div>
				{action.text}
			</PrettyButton>
		{/if}
	{/each}
</div>

<style>
	.multi-action-container {
		display: flex;
		position: sticky;
		background-color: var(--surface-tonal-a20);
		z-index: 10;
		bottom: 1rem;
		width: max-content;
		border-radius: var(--border-radius);
		overflow: hidden;
		box-shadow: 2px 2px 15px var(--dark-a0);
	}
	.icon-holder {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 1.3rem;
		height: 1.3rem;
	}
</style>
