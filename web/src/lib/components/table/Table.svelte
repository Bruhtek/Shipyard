<script lang="ts" generics="T extends { ID: string }">
	import type { TableColumn } from '$lib/types/Table';
	import type { Snippet } from 'svelte';
	import SortingButton from '$lib/components/table/SortingButton.svelte';

	type Props = {
		columns: TableColumn[];
		data: T[];
		sortedBy: string;
		sortedDirection: 'asc' | 'desc';
		Row: Snippet<[T]>;
	};

	let {
		columns,
		data,
		sortedBy = $bindable('ID'),
		sortedDirection = $bindable('asc'),
		Row
	}: Props = $props();
</script>

<table class="table">
	<thead class="thead">
		<tr>
			<th class="table-header">
				<SortingButton
					bind:current={sortedBy}
					bind:currentDirection={sortedDirection}
					sortByKey="ID"
				>
					ID
				</SortingButton>
			</th>
			{#each columns as col (col.label)}
				<th class="table-header">
					{#if col.sortable}
						<SortingButton
							bind:current={sortedBy}
							bind:currentDirection={sortedDirection}
							sortByKey={col.label}
						>
							{col.label}
						</SortingButton>
					{:else}
						<span class="unsortable">
							{col.label}
						</span>
					{/if}
				</th>
			{/each}
		</tr>
	</thead>
	<tbody class="tbody">
		{#each data as rowData (rowData.ID)}
			{@render Row(rowData)}
		{/each}
	</tbody>
</table>

<style>
	.table {
		width: 100%;
		border-collapse: collapse;
	}
	.thead th:first-child {
		border-top-left-radius: var(--border-radius);
	}
	.thead th:last-child {
		border-top-right-radius: var(--border-radius);
	}
	.thead th {
		background-color: var(--surface-tonal-a10);
	}
	.table-header {
		padding: 0.5rem;
	}
	.unsortable {
		opacity: 0.5;
	}
</style>
