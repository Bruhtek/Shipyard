<script lang="ts">
	import { URLPrefix } from '$lib';
	import EnvStore from '$lib/stores/EnvStore.svelte';
	import Table from '$lib/components/table/Table.svelte';
	import type { TableColumn } from '$lib/types/Table';
	import TruncatedID from '$lib/components/table/TruncatedID.svelte';
	import { type Image, TImageResponse } from '$lib/types/docker/Image';
	import HumanReadableSize from '$lib/components/table/HumanReadableSize.svelte';
	import ImageActionButtons from '$lib/components/table/image/ImageActionButtons.svelte';
	import Badge from '$lib/components/fragments/Badge.svelte';
	import TableHeader from '$lib/components/table/TableHeader.svelte';
	import { DATA_FETCHING_INTERVAL } from '$lib/consts';
	import { sortDataByKey } from '$lib/utils/displayUtils';

	let imageData = $state<Image[]>([]);
	let loading = $state(true);

	async function fetchData() {
		const res = await fetch(`${URLPrefix}/api/env/${EnvStore.name}/images`);
		if (res.ok) {
			const data = await res.json();
			const parsed = TImageResponse.parse(data);

			imageData = Object.entries(parsed.Images).map(([, v]) => v);
			loading = false;
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
		}, DATA_FETCHING_INTERVAL);

		return () => {
			clearInterval(interval);
		};
	});

	let filter = $state('');
	let sortedBy = $state('Repository');
	let sortedDirection = $state<'asc' | 'desc'>('asc');
	let filteredData = $derived.by(() => {
		const query = filter.toLowerCase().trim();
		if (query === '') {
			return imageData;
		}
		if (query === 'unused') {
			return imageData.filter((image) => !image.Used);
		}
		if (query === 'used') {
			return imageData.filter((image) => image.Used);
		}
		return imageData.filter((image) => {
			return (
				image.ID.toLowerCase().startsWith(query) ||
				image.Repository.toLowerCase().includes(query) ||
				image.Tag.toLowerCase().includes(query)
			);
		});
	});

	let sortedData = $derived(sortDataByKey(filteredData, sortedBy, sortedDirection));

	const tableColumns: TableColumn[] = [
		{ label: 'ID', sortable: true },
		{ label: 'Repository', sortable: true },
		{ label: 'Tag', sortable: true },
		{ label: 'Size', sortable: true },
		{ label: '' }
	];
</script>

<svelte:head>
	<title>Images - {EnvStore.name} - Shipyard</title>
</svelte:head>

<TableHeader title="Images" bind:query={filter} />

<Table columns={tableColumns} data={sortedData} bind:sortedBy bind:sortedDirection {loading}>
	{#snippet Row(r: Image)}
		<td>
			<TruncatedID id={r.ID} />
			{#if !r.Used}
				<Badge background="var(--yellow-a20)" color="var(--dark-a0)">Unused</Badge>
			{/if}
		</td>
		<td>
			{r.Repository}
		</td>
		<td>{r.Tag}</td>
		<td class="align-right">
			<HumanReadableSize size={r.Size} />
		</td>
		<td class="set-width">
			<ImageActionButtons id={r.ID} repo={r.Repository} tag={r.Tag} />
		</td>
	{/snippet}
</Table>

<style>
	.align-right {
		text-align: right;
	}
	td.set-width {
		width: 2rem;
		padding: 0;
	}
</style>
