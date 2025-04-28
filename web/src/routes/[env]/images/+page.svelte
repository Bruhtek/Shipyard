<script lang="ts">
	import { URLPrefix } from '$lib';
	import EnvStore from '$lib/stores/EnvStore.svelte';
	import Table from '$lib/components/table/Table.svelte';
	import type { TableColumn } from '$lib/types/Table';
	import TruncatedID from '$lib/components/table/TruncatedID.svelte';
	import ContainerState from '$lib/components/table/container/ContainerState.svelte';
	import ContainerActionButtons from '$lib/components/table/container/ContainerActionButtons.svelte';
	import { type Image, TImageResponse } from '$lib/types/docker/Image';
	import HumanReadableSize from '$lib/components/table/HumanReadableSize.svelte';
	import ImageActionButtons from '$lib/components/table/image/ImageActionButtons.svelte';
	import Badge from '$lib/components/fragments/Badge.svelte';

	let imageData = $state<Image[]>([]);

	async function fetchData() {
		const res = await fetch(`${URLPrefix}/api/env/${EnvStore.name}/images`);
		if (res.ok) {
			const data = await res.json();
			const parsed = TImageResponse.parse(data);

			imageData = Object.entries(parsed.Images).map(([, v]) => v);
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
		if (imageData.length === 0) {
			return imageData;
		}
		const key = sortedBy as keyof Image;
		if (!(key in imageData[0])) {
			return imageData;
		}

		return imageData.toSorted((a, b) => {
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
		{ label: 'ID', sortable: true },
		{ label: 'Repository', sortable: true },
		{ label: 'Tag', sortable: true },
		{ label: 'Size', sortable: true },
		{ label: '' }
	];
</script>

Containers

<Table columns={tableColumns} data={sortedData} bind:sortedBy bind:sortedDirection>
	{#snippet Row(r: Image)}
		<td>
			<TruncatedID id={r.ID} />
		</td>
		<td>
			{r.Repository}
			{#if !r.Used}
				<Badge background="var(--yellow-a20)" color="var(--dark-a0)">Unused</Badge>
			{/if}
		</td>
		<td>{r.Tag}</td>
		<td class="align-right">
			<HumanReadableSize size={r.Size} />
		</td>
		<td>
			<ImageActionButtons id={r.ID} repo={r.Repository} tag={r.Tag} />
		</td>
	{/snippet}
</Table>

<style>
	.align-right {
		text-align: right;
	}
</style>
