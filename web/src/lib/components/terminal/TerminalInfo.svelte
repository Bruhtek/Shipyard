<script lang="ts">
	import type { Terminal } from '$lib/terminal/TerminalStore.svelte';
	import Divider from '$lib/components/fragments/Divider.svelte';
	import CaretUp from '~icons/ph/caret-up';
	import Cube from '~icons/ph/cube';
	import CubeFocus from '~icons/ph/cube-focus';
	import TerminalStatusToIcon from '$lib/components/terminal/TerminalStatusToIcon.svelte';
	import type { ActionStatus } from '$lib/types/Action';

	type Props = {
		terminal: Terminal;
		show: boolean;
		cmdStatus: ActionStatus;
		shortFormWidth: number;
	}

	let {
		terminal,
	  	show,
	  	cmdStatus,
	  	shortFormWidth = $bindable()
	}: Props = $props();
</script>

<div class="info">
	<div class="left"
		 bind:clientWidth={shortFormWidth}
	>
		<span
			class="icon-container"
			class:rotated={show}
		>
			<CaretUp />
		</span>
		{terminal.Action}
		<TerminalStatusToIcon status={cmdStatus} />
	</div>
	<div class="right">
		<Cube height="100%" width="max-content" /> {terminal.Environment}
		<Divider />
		<CubeFocus height="120%" width="max-content" /> {terminal.Object}
	</div>
</div>

<style>
    .info {
        height: 2rem;
        width: 100%;
        display: flex;
        justify-content: space-between;
        padding: 0.3rem;
    }
    .left, .right {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }

    .icon-container {
        transition: transform 0.2s ease-in-out;
        display: block;
    }
    .rotated {
        transform: rotate(180deg);
    }
</style>