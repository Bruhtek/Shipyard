<script lang="ts">
	import House from '~icons/ph/house';
	import QuestionMark from '~icons/ph/question-mark';

	const envTypeToIcon = (envType: string) => {
		switch (envType) {
			case 'local':
				return House;
			default:
				return QuestionMark;
		}
	}

	type Props = {
		envType: string;
	}

	let { envType }: Props = $props();

	const Icon = $derived(envTypeToIcon(envType));
</script>

<Icon width="1.4rem" height="1.4rem" />