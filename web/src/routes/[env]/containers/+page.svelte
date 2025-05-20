<script lang="ts">
	import { type Container, TContainerResponse } from '$lib/types/docker/Container';
	import { URLPrefix } from '$lib';
	import EnvStore from '$lib/stores/EnvStore.svelte';
	import Table from '$lib/components/table/Table.svelte';
	import type { TableColumn } from '$lib/types/Table';
	import TruncatedID from '$lib/components/table/TruncatedID.svelte';
	import ContainerState from '$lib/components/table/container/ContainerState.svelte';
	import ContainerActionButtons from '$lib/components/table/container/ContainerActionButtons.svelte';
	import ContainerImage from '$lib/components/table/container/ContainerImage.svelte';
	import TableHeader from '$lib/components/table/TableHeader.svelte';
	import { DATA_FETCHING_INTERVAL } from '$lib/consts';
	import { sortDataByKey } from '$lib/utils/displayUtils';
	import TerminalStore from '$lib/terminal/TerminalStore.svelte';

	let containerData = $state<Container[]>([]);
	let loading = $state(true);

	async function fetchData() {
		const res = await fetch(`${URLPrefix}/api/env/${EnvStore.name}/containers`);
		if (res.ok) {
			const data = await res.json();
			const parsed = TContainerResponse.parse(data);

			containerData = Object.entries(parsed.Containers).map(([, v]) => v);
			loading = false;
		} else {
			console.error('Failed to fetch container data:', res.statusText);
		}
	}

	$effect(() => {
		if (EnvStore.name === '') {
			return;
		}
		TerminalStore.subscribeActionFinished(fetchData);
		fetchData();
		const interval = setInterval(() => {
			fetchData();
		}, DATA_FETCHING_INTERVAL);

		return () => {
			TerminalStore.unsubscribeActionFinished(fetchData);
			clearInterval(interval);
		};
	});

	let filter = $state('');
	let sortedBy = $state('Name');
	let sortedDirection = $state<'asc' | 'desc'>('asc');
	let filteredData = $derived.by(() => {
		let query = filter.trim().toLowerCase();
		if (query === '') {
			return containerData;
		}
		return containerData.filter((container) => {
			return (
				container.Name.toLowerCase().includes(query) ||
				container.Image.toLowerCase().includes(query) ||
				container.ID.toLowerCase().startsWith(query)
			);
		});
	});

	let sortedData = $derived(sortDataByKey(filteredData, sortedBy, sortedDirection));

	const tableColumns: TableColumn[] = [
		{ label: 'ID', sortable: true },
		{ label: 'Name', sortable: true },
		{ label: 'Image', sortable: true },
		{ label: 'State', sortable: true },
		{ label: '' }
	];
</script>

<svelte:head>
	<title>Containers - {EnvStore.name} - Shipyard</title>
</svelte:head>

<TableHeader title="Containers" bind:query={filter} />

<Table columns={tableColumns} data={sortedData} bind:sortedBy bind:sortedDirection {loading}>
	{#snippet Row(r: Container)}
		<td>
			<TruncatedID id={r.ID} />
		</td>
		<td>{r.Name}</td>
		<td>
			<ContainerImage image={r.Image} />
		</td>
		<td>
			<ContainerState state={r.State} />
		</td>
		<td class="set-width">
			<ContainerActionButtons id={r.ID} name={r.Name} />
		</td>
	{/snippet}
</Table>

<style>
	td.set-width {
		width: 2rem;
		padding: 0;
	}
</style>
