<script lang="ts">
	import { URLPrefix } from '$lib';
	import EnvStore from '$lib/stores/EnvStore.svelte';
	import Table from '$lib/components/table/Table.svelte';
	import type { TableColumn } from '$lib/types/Table';
	import TableHeader from '$lib/components/table/TableHeader.svelte';
	import { DATA_FETCHING_INTERVAL } from '$lib/consts';
	import { sortDataByKey } from '$lib/utils/displayUtils';
	import TerminalStore from '$lib/terminal/TerminalStore.svelte';
	import { type Stack, type StackWithID, TStackResponse } from '$lib/types/docker/Stack';

	let stackData = $state<StackWithID[]>([]);
	let loading = $state(true);
	let abortController: AbortController | null = null;

	async function fetchData() {
		if (abortController) {
			abortController.abort();
		}
		abortController = new AbortController();

		const res = await fetch(`${URLPrefix}/api/env/${EnvStore.name}/stacks`, {
			signal: abortController.signal
		});
		if (res.ok) {
			const data = await res.json();
			const parsed = TStackResponse.parse(data);

			stackData = Object.entries(parsed.Stacks).map(([, v]) => {
				v.ID = v.ID || v.ConfigFiles || v.Name; // Ensure ID is set for each stack
				return v as StackWithID;
			});
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

	let searchQuery = $state('');
	let sortedBy = $state('Name');
	let sortedDirection = $state<'asc' | 'desc'>('asc');
	let filter = $state({});

	let filteredData = $derived.by(() => {
		let query = searchQuery.trim().toLowerCase();
		if (query === '') {
			return stackData;
		}
		if (query.includes(':')) {
			const parts = query.split(' ');
			return stackData.filter((stack) => {
				return parts.every((part) => {
					const [key, value] = part.split(':');
					if (key === 'name') {
						return stack.Name.toLowerCase().includes(value.toLowerCase());
					} else if (key === 'state') {
						return stack.Status.toLowerCase().includes(value.toLowerCase());
					}
					return false;
				});
			});
		}

		return stackData.filter((stack) => {
			return (
				stack.Name.toLowerCase().includes(query) ||
				stack.ConfigFiles.toLowerCase().includes(query)
			);
		});
	});

	let sortedData = $derived(sortDataByKey(filteredData, sortedBy, sortedDirection));

	const tableColumns: TableColumn[] = [
		{ label: 'Name', sortable: true },
		{ label: 'Status', sortable: true },
		{ label: '' }
	];
</script>

<svelte:head>
	<title>Stacks - {EnvStore.name} - Shipyard</title>
</svelte:head>

<TableHeader title="Stacks" bind:query={searchQuery} />

<Table
	columns={tableColumns}
	data={sortedData}
	bind:filter
	bind:sortedBy
	bind:sortedDirection
	{loading}
>
	{#snippet Row(r: StackWithID)}
		<td>{r.Name}</td>
		<td>{r.Status}</td>
		<td class="set-width"> A </td>
	{/snippet}
</Table>

<style>
	td.set-width {
		width: 2rem;
		padding: 0;
	}
</style>
