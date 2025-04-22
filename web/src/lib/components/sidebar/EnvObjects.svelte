<script lang="ts">
	import Cube from '~icons/ph/cube';
	import Stack from '~icons/ph/stack';
	import Disc from '~icons/ph/disc';
	import Database from '~icons/ph/database';
	import Network from '~icons/ph/network';

	import type { Component } from 'svelte';
	import EnvStore from '$lib/stores/EnvStore.svelte';

	const links: [string, string, Component][] = [
		["Containers", "/containers", Cube],
	  	["Stacks", '/stacks', Stack],
	  	["Images", "/images", Disc],
	  	["Volumes", "/volumes", Database],
	  	["Networks", "/networks", Network],
	]
</script>

<div class="objects">
	{#each links as l (l[1])}
		{@const Icon = l[2]}
		<a
			href="/{EnvStore.name}{l[1]}"
		   	class="link"
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
</style>