<script lang="ts" generics="T extends { ID: string }">
	import type { TableColumn } from '$lib/types/Table';
	import type { Snippet } from 'svelte';
	import SortingButton from '$lib/components/table/SortingButton.svelte';
	import Shadow from '$lib/components/fragments/Shadow.svelte';
	import FilterPopup from '$lib/components/table/FilterPopup.svelte';
	import PopupActionMultiButton from '$lib/components/table/actionButtons/PopupActionMultiButton.svelte';
	import type { ActionData } from '$lib/components/table/actionButtons/ActionData';
	import PopupActionButton from '$lib/components/table/actionButtons/PopupActionButton.svelte';

	type Props = {
		columns: TableColumn[];
		data: T[];
		sortedBy: string;
		sortedDirection: 'asc' | 'desc';
		loading?: boolean;
		filter?: object;
		Row: Snippet<[T]>;
		objName: (obj: T) => string;
		popupActions: ActionData[];
	};

	let {
		columns,
		data,
		sortedBy = $bindable('ID'),
		sortedDirection = $bindable('asc'),
		filter = $bindable({}),
		Row,
		loading = false,
		objName = (obj: T) => obj.ID, // by default, use ID
		popupActions
	}: Props = $props();

	let selected = $state<string[]>([]);
	let selectedNames = $state<string[]>([]);

	function clearSelection() {
		selected = [];
		selectedNames = [];
	}
	function selectAll() {
		selected = data.map((d) => d.ID);
		selectedNames = data.map((d) => objName(d));
	}
</script>

<table class="table">
	<thead class="thead">
		<tr>
			<!-- This is for select -->
			<th></th>
			{#each columns as col (col.label)}
				<th class="table-header">
					<div class="header-container">
						{#if col.sortable}
							<SortingButton
								bind:current={sortedBy}
								bind:currentDirection={sortedDirection}
								sortByKey={col.key || col.label}
							>
								{col.label}
							</SortingButton>
						{:else}
							<span class="unsortable">
								{col.label}
							</span>
						{/if}
						{#if col.filterOptions}
							<FilterPopup id={col.key || col.label} />
						{/if}
					</div>
				</th>
			{/each}
		</tr>
	</thead>
	<tbody class="tbody">
		{#each data as rowData (rowData.ID)}
			<tr class="t-row" class:selected={selected.includes(rowData.ID)}>
				<td class="select-box">
					<label class="row-checkbox">
						<input
							type="checkbox"
							checked={selected.includes(rowData.ID)}
							onchange={() => {
								if (selected.includes(rowData.ID)) {
									selected = selected.filter((id) => id !== rowData.ID);
									selectedNames = selectedNames.filter(
										(name) => name !== objName(rowData)
									);
								} else {
									selected = [...selected, rowData.ID];
									selectedNames = [...selectedNames, objName(rowData)];
								}
							}}
						/>
					</label>
				</td>
				{@render Row(rowData)}
				<td class="popup-actions">
					<PopupActionButton
						id={rowData.ID}
						name={objName(rowData)}
						actions={popupActions}
					/>
				</td>
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

{#if selected.length > 0}
	<PopupActionMultiButton
		actions={popupActions}
		ids={selected}
		names={selectedNames}
		clear={clearSelection}
		{selectAll}
	/>
{/if}

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
	.thead {
		position: sticky;
		top: 3.5rem;
		z-index: 10;
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
	.header-container {
		width: 100%;
		height: 100%;
		display: flex;
		align-items: center;
		gap: 0.2rem;
	}
	.unsortable {
		opacity: 0.7;
		text-align: left;
	}

	.select-box {
		width: 2.5rem;
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
	.t-row:last-child {
		border-bottom-left-radius: var(--border-radius);
		border-bottom-right-radius: var(--border-radius);
	}
	.t-row.selected :global(td) {
		background-color: var(--surface-tonal-a20);
	}

	.table :global(tr:hover) {
		background-color: var(--surface-tonal-a10);
	}

	.row-checkbox {
		display: flex;
		justify-content: center;
		align-items: center;
		margin: -10px;
		padding: 10px;
	}

	td.popup-actions {
		width: 2rem;
		padding: 0;
	}
</style>
