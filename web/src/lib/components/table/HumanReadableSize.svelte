<script lang="ts">
	type Props = {
		size: number;
	};

	let { size }: Props = $props();

	const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB'];
	const unitIndex = $derived(Math.floor(Math.log(size) / Math.log(1024)));
	const unit = $derived.by(() => {
		if (unitIndex < 0 || unitIndex >= units.length) {
			return 'B';
		}
		return units[unitIndex];
	});
	const humanReadableSize = $derived.by(() => {
		if (size === 0) {
			return '0 B';
		}
		const sizeInUnit = size / Math.pow(1024, unitIndex);
		return `${sizeInUnit.toFixed(2)} ${unit}`;
	});
</script>

<span class="human-readable-size">
	{humanReadableSize}
</span>
