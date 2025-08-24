<script lang="ts">
	import { type Network, TNetworkResponse } from '$lib/types/docker/Network';
	import { URLPrefix } from '$lib';
	import EnvStore from '$lib/stores/EnvStore.svelte';
	import { DATA_FETCHING_INTERVAL } from '$lib/consts';
	import { sortDataByKey } from '$lib/utils/displayUtils';
	import type { TableColumn } from '$lib/types/Table';
	import TableHeader from '$lib/components/table/TableHeader.svelte';
	import TruncatedID from '$lib/components/table/TruncatedID.svelte';
	import Table from '$lib/components/table/Table.svelte';
	import Badge from '$lib/components/fragments/Badge.svelte';
	import TerminalStore from '$lib/terminal/TerminalStore.svelte';
	import Trash from '~icons/ph/trash';
	import NetworkAction from '$lib/websocket/actions/Network';
	import { SvelteURL } from 'svelte/reactivity';
	import { goto } from '$app/navigation';
	import Funnel from '~icons/ph/funnel';
	import type { ActionData } from '$lib/components/table/actionButtons/ActionData';

	let networkData = $state<Network[]>([]);
	let loading = $state(true);
	let abortController: AbortController | null = null;

	async function fetchData() {
		if (abortController) {
			abortController.abort();
		}
		abortController = new AbortController();

		const res = await fetch(`${URLPrefix}/api/env/${EnvStore.name}/networks`, {
			signal: abortController.signal
		});

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
		if (query.includes(':')) {
			const parts = query.split(' ');
			return networkData.filter((network) => {
				return parts.every((part) => {
					const [key, value] = part.split(':');
					if (key === 'id') {
						return network.ID.toLowerCase().startsWith(value.toLowerCase());
					} else if (key === 'name') {
						return network.Name.toLowerCase().includes(value.toLowerCase());
					} else if (key === 'container') {
						return network.Containers.some(
							(container) =>
								container.Name.toLowerCase().includes(value.toLowerCase()) ||
								container.ID.toLowerCase().startsWith(value.toLowerCase())
						);
					} else if (key === 'status') {
						if (value === 'unused') {
							return !network.Containers.length;
						} else if (value === 'used') {
							return network.Containers.length > 0;
						}
						return false;
					} else if (key === 'driver') {
						return network.Driver.toLowerCase().startsWith(value.toLowerCase());
					} else if (key === 'scope') {
						return network.Scope.toLowerCase().startsWith(value.toLowerCase());
					}
					return false;
				});
			});
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

	const popupActions: ActionData[] = [
		{
			text: 'Remove',
			icon: Trash,
			onClick: (id: string) => {
				NetworkAction(EnvStore.name, 'remove', id);
			},
			hoverBackground: 'var(--red-a20)',
			hoverColor: 'var(--dark-a0)'
		},
		{
			text: 'View Containers',
			icon: Funnel,
			onClick: (id: string, name: string) => {
				const url = new SvelteURL(window.location.href);
				url.pathname = `/${EnvStore.name}/containers`;
				url.searchParams.set('query', `network:${name}`);
				goto(url.toString());
			},
			hoverBackground: 'var(--primary-a20)',
			hoverColor: 'var(--dark-a0)',
			singleOnly: true
		}
	];
</script>

<svelte:head>
	<title>Networks - {EnvStore.name} - Shipyard</title>
</svelte:head>

<TableHeader title="Networks" bind:query={filter} />

<Table
	columns={tableColumns}
	data={sortedData}
	bind:sortedBy
	bind:sortedDirection
	{loading}
	objName={(r) => r.Name}
	{popupActions}
>
	{#snippet Row(r: Network)}
		<td>
			<TruncatedID id={r.ID} />
			{#if ['bridge', 'host', 'none'].includes(r.Name)}
				<Badge background="var(--green-a20)" color="var(--dark-a0)">Built-in</Badge>
			{/if}
			{#if !r.Containers.length}
				<Badge background="var(--yellow-a20)" color="var(--dark-a0)">Unused</Badge>
			{/if}
		</td>
		<td>{r.Name}</td>
		<td>{r.Driver}</td>
		<td>{r.Scope}</td>
	{/snippet}
</Table>
