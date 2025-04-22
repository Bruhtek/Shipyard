<script lang="ts">

	import EnvIcon from '$lib/components/sidebar/Misc/EnvIcon.svelte';
		import EnvStore from '$lib/stores/EnvStore.svelte';
		import EnvObjects from '$lib/components/sidebar/EnvObjects.svelte';

		type Props = {
		envType: string;
		name: string;
	}

	let { envType, name }: Props = $props();

	let objectsHeight = $state(120);
</script>

<div
	class="env-item"
	style="--objectsHeight: {objectsHeight}px;"
	class:active={name === EnvStore.name}
>
	<a
		href="/{name}"
		class="main"
	>
		<EnvIcon envType={envType} />
		{name}
	</a>

	<div
		class="objects"
		class:active={name === EnvStore.name}
		bind:clientHeight={objectsHeight}
	>
		<EnvObjects />
	</div>
</div>

<style>
    .env-item {
        --height: 1.4rem;
		max-height: calc(var(--height) + 1rem);
		transition: all 0.2s ease-in-out;
		overflow: hidden;
        display: flex;
		flex-direction: column;
    }
	.main {
        line-height: var(--height);
        height: calc(var(--height) + 1rem);
        display: flex;
        gap: 0.2rem;
        padding: 0.5rem;
        font-size: 1.1rem;
        align-items: center;
    }
	.env-item.active .main {
		text-shadow: 0 0 0.2rem var(--primary-a50);
		color: var(--primary-a50);
    }
    .env-item.active .main :global(svg) {
        filter: drop-shadow(0 0 0.2rem var(--primary-a50));
	}
    .env-item.active {
		max-height: calc(var(--height) + 1rem + var(--objectsHeight));
    }
</style>