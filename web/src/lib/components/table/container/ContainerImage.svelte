<script lang="ts">
	import TruncatedID from '$lib/components/table/TruncatedID.svelte';

	type Props = {
		image: string;
	};

	let { image }: Props = $props();

	let isOnlyID = $derived(image.startsWith('sha256:'));
	let isBothTagAndID = $derived(image.includes('@sha256'));

	let imageTag = $derived.by(() => {
		if (isOnlyID) {
			return '';
		}
		if (isBothTagAndID) {
			return image.split('@')[0];
		}
		return image;
	});

	let imageID = $derived.by(() => {
		if (isOnlyID) {
			return image;
		}
		if (isBothTagAndID) {
			return image.split('@sha256:')[1];
		}
		return '';
	});
</script>

{#if isOnlyID}
	<TruncatedID id={image} />
{:else if isBothTagAndID}
	{imageTag}
	<span class="separator">@</span>
	<TruncatedID id={imageID} />
{:else}
	{image}
{/if}
