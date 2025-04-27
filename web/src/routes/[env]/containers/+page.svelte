<script lang="ts">
	import { type Container, TContainerResponse } from '$lib/types/docker/Container';
	import { URLPrefix } from '$lib';
	import EnvStore from '$lib/stores/EnvStore.svelte';
	import Table from '$lib/components/table/Table.svelte';
	import type { TableColumn } from '$lib/types/Table';
	import TruncatedID from '$lib/components/table/TruncatedID.svelte';
	import ContainerState from '$lib/components/table/ContainerState.svelte';

	let containerData = $state<Container[]>([]);

	async function fetchData() {
		const res = await fetch(`${URLPrefix}/api/env/${EnvStore.name}/containers`);
		if (res.ok) {
			const data = await res.json();
			const parsed = TContainerResponse.parse(data);

			containerData = Object.entries(parsed.Containers).map(([, v]) => v);
		} else {
			console.error('Failed to fetch container data:', res.statusText);
		}
	}

	$effect(() => {
		if (EnvStore.name === '') {
			return;
		}

		fetchData();
		const interval = setInterval(() => {
			fetchData();
		}, 3 * 1000);

		return () => {
			clearInterval(interval);
		};
	});

	let sortedBy = $state('ID');
	let sortedDirection = $state<'asc' | 'desc'>('asc');
	let sortedData = $derived.by(() => {
		const sortDirection = sortedDirection === 'asc' ? 1 : -1;
		if (containerData.length === 0) {
			return containerData;
		}
		const key = sortedBy as keyof Container;
		if (!(key in containerData[0])) {
			return containerData;
		}

		return containerData.toSorted((a, b) => {
			if (a[key] < b[key]) {
				return -1 * sortDirection;
			}
			if (a[key] > b[key]) {
				return 1 * sortDirection;
			}
			return 0;
		});
	});
	const tableColumns: TableColumn[] = [
		{ label: 'Name', sortable: true },
		{ label: 'Image', sortable: true },
		{ label: 'State', sortable: true },
		{ label: 'More' }
	];
</script>

Containers

<Table columns={tableColumns} data={sortedData} bind:sortedBy bind:sortedDirection>
	{#snippet Row(r: Container)}
		<tr class="t-row">
			<td>
				<TruncatedID id={r.ID} />
			</td>
			<td>{r.Name}</td>
			<td>{r.Image}</td>
			<td>
				<ContainerState state={r.State} />
			</td>
			<td> Button here </td>
		</tr>
	{/snippet}
</Table>

<style>
	.t-row td {
		border: 0.1rem solid var(--surface-tonal-a10);
		padding: 0.2rem 0.5rem;
	}
	.t-row:last-child td:first-child {
		border-bottom-left-radius: var(--border-radius);
	}
	.t-row:last-child td:last-child {
		border-bottom-right-radius: var(--border-radius);
	}
</style>
