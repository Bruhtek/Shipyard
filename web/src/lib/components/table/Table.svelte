<script lang="ts" generics="T extends { ID: string }">
	import type { TableColumn } from '$lib/types/Table';
	import type { Snippet } from 'svelte';
	import SortingButton from '$lib/components/table/SortingButton.svelte';
	import Shadow from '$lib/components/fragments/Shadow.svelte';

	type Props = {
		columns: TableColumn[];
		data: T[];
		sortedBy: string;
		sortedDirection: 'asc' | 'desc';
		loading?: boolean;
		Row: Snippet<[T]>;
	};

	let {
		columns,
		data,
		sortedBy = $bindable('ID'),
		sortedDirection = $bindable('asc'),
		Row,
		loading = false
	}: Props = $props();

	let selected = $state<string[]>([]);
</script>

<table class="table">
	<thead class="thead">
		<tr>
			<!-- This is for select -->
			<th></th>
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
			<tr class="t-row">
				<td>
					<input
						type="checkbox"
						checked={selected.includes(rowData.ID)}
						onchange={() => {
							if (selected.includes(rowData.ID)) {
								selected = selected.filter((id) => id !== rowData.ID);
							} else {
								selected = [...selected, rowData.ID];
							}
						}}
					/>
				</td>
				{@render Row(rowData)}
			</tr>
		{/each}
		{#if data.length === 0 && loading}
			{#each { length: 3 }}
				<tr class="t-row">
					<td>
						<input type="checkbox" checked={false} disabled />
					</td>
					{#each columns as col (col.label)}
						<td>
							<Shadow />
						</td>
					{/each}
				</tr>
			{/each}
		{/if}
	</tbody>
</table>

{#if data.length === 0 && !loading}
	<div class="no-data">
		<p>Empty</p>
	</div>
{/if}

<style>
	.no-data {
		text-align: center;
		padding: 1rem;
		border: 0.1rem solid var(--surface-tonal-a10);
		border-bottom-left-radius: var(--border-radius);
		border-bottom-right-radius: var(--border-radius);
	}
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
		background-color: var(--surface-tonal-a20);
	}
	.table-header {
		padding: 0.5rem;
	}
	.unsortable {
		opacity: 0.7;
	}

	.t-row :global(td) {
		border: 0.1rem solid var(--surface-tonal-a10);
		padding: 0.2rem 0.5rem;
	}
	.t-row:last-child :global(td:first-child) {
		border-bottom-left-radius: var(--border-radius);
	}
	.t-row:last-child :global(td:last-child) {
		border-bottom-right-radius: var(--border-radius);
	}

	.table :global(tr:hover) {
		background-color: var(--surface-tonal-a10);
	}
</style>
