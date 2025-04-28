<script lang="ts">
	type Props = {
		id: string;
	};

	let { id }: Props = $props();

	const length = $derived(id.length);
	let actualId = $derived.by(() => {
		if (id.startsWith('sha256:')) {
			return id.slice(7);
		}
		return id;
	});
	const truncated = $derived(actualId.slice(0, 12));
</script>

{#if length > 12}
	<span class="truncated-id">
		{truncated}
	</span>
{:else}
	<span class="full-id">{id}</span>
{/if}
