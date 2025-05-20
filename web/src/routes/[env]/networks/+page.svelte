<script lang="ts">
	import { type Network, TNetworkResponse } from '$lib/types/docker/Network';
	import { URLPrefix } from '$lib';
	import EnvStore from '$lib/stores/EnvStore.svelte';
	import { DATA_FETCHING_INTERVAL } from '$lib/consts';
	import { sortDataByKey } from '$lib/utils/displayUtils';
	import type { TableColumn } from '$lib/types/Table';
	import TableHeader from '$lib/components/table/TableHeader.svelte';
	import { load } from '../+layout';
	import TruncatedID from '$lib/components/table/TruncatedID.svelte';
	import Table from '$lib/components/table/Table.svelte';
	import Badge from '$lib/components/fragments/Badge.svelte';

	let networkData = $state<Network[]>([]);
	let loading = $state(true);

	async function fetchData() {
		const res = await fetch(`${URLPrefix}/api/env/${EnvStore.name}/networks`);
		if (res.ok) {
			const data = await res.json();
			const parsed = TNetworkResponse.parse(data);

			networkData = Object.entries(parsed.Networks).map(([, v]) => v);
			loading = false;
		} else {
			console.error('Failed to fetch network data:', res.statusText);
		}
	}

	$effect(() => {
		if (EnvStore.name === '') {
			return;
		}

		fetchData();
		const interval = setInterval(() => {
			fetchData();
		}, DATA_FETCHING_INTERVAL);

		return () => {
			clearInterval(interval);
		};
	});

	let filter = $state('');
	let sortedBy = $state('Name');
	let sortedDirection = $state<'asc' | 'desc'>('asc');
	let filteredData = $derived.by(() => {
		const query = filter.toLowerCase().trim();
		if (query === '') {
			return networkData;
		}
		if (query === 'unused') {
			return networkData.filter((network) => !network.Containers.length);
		}
		if (query === 'used') {
			return networkData.filter((network) => network.Containers.length);
		}
		return networkData.filter((network) => {
			return (
				network.Name.toLowerCase().includes(query) ||
				network.ID.toLowerCase().startsWith(query) ||
				network.Driver.toLowerCase().startsWith(query)
			);
		});
	});

	let sortedData = $derived(sortDataByKey(filteredData, sortedBy, sortedDirection));

	const tableColumns: TableColumn[] = [
		{ label: 'ID', sortable: true },
		{ label: 'Name', sortable: true },
		{ label: 'Driver', sortable: true },
		{ label: 'Scope', sortable: true },
		{ label: '' }
	];
</script>

<svelte:head>
	<title>Networks - {EnvStore.name} - Shipyard</title>
</svelte:head>

<TableHeader title="Networks" bind:query={filter} />

<Table columns={tableColumns} data={sortedData} bind:sortedBy bind:sortedDirection {loading}>
	{#snippet Row(r: Network)}
		<td>
			<TruncatedID id={r.ID} />
			{#if !r.Containers.length}
				<Badge background="var(--yellow-a20)" color="var(--dark-a0)">Unused</Badge>
			{/if}
		</td>
		<td>{r.Name}</td>
		<td>{r.Driver}</td>
		<td>{r.Scope}</td>
	{/snippet}
</Table>
