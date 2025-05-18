<script lang="ts">
	import InvisibleButton from '$lib/components/fragments/InvisibleButton.svelte';
	import DotsThree from '~icons/ph/dots-three';
	import PopupStore from '$lib/stores/CurrentPopup.svelte';

	type Props = {
		id: string;
	};

	let { id }: Props = $props();
</script>

<InvisibleButton center={true} onclick={() => PopupStore.toggle(id)} class="show-button">
	<DotsThree width="1.5rem" height="1.5rem" />
</InvisibleButton>
