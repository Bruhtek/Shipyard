<script lang="ts">
	import Cube from '~icons/ph/cube';
	import Stack from '~icons/ph/stack';
	import Disc from '~icons/ph/disc';
	import Database from '~icons/ph/database';
	import Network from '~icons/ph/network';
	import type { Component } from 'svelte';
	import EnvStore from '$lib/stores/EnvStore.svelte';
	import { page } from '$app/state';

	const links: [string, string, Component][] = [
		['Containers', '/containers', Cube],
		['Stacks', '/stacks', Stack],
		['Images', '/images', Disc],
		['Volumes', '/volumes', Database],
		['Networks', '/networks', Network]
	];

	let currentRoute = $derived(decodeURI(page.url.pathname));
</script>

<div class="objects">
	{#each links as l (l[1])}
		{@const Icon = l[2]}
		<a
			href="/{EnvStore.name}{l[1]}"
			class="link"
			class:active={currentRoute.startsWith(`/${EnvStore.name}${l[1]}`)}
		>
			<Icon />
			{l[0]}
		</a>
	{/each}
</div>

<style>
	.objects {
		padding: 0.25rem 0.5rem 1rem;
	}
	.link {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}
	.link.active {
		text-shadow: 0 0 0.1rem var(--primary-a50);
		color: var(--primary-a50);
	}
	.link.active :global(svg) {
		filter: drop-shadow(0 0 0.1rem var(--primary-a50));
	}
</style>
