<script lang="ts">
	type Props = {
		id: string;
	};

	let { id }: Props = $props();

	const length = $derived(id.length);
	const truncated = $derived(id.slice(0, 12));
</script>

{#if length > 12}
	<span class="truncated-id">
		{truncated}
	</span>
{:else}
	<span class="full-id">{id}</span>
{/if}
