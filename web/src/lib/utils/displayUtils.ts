export const sortDataByKey = <T extends object>(
	data: T[],
	sortedBy: string,
	direction: 'asc' | 'desc'
): T[] => {
	if (data.length === 0) return data;

	const key = sortedBy as keyof T;
	if (!(key in data[0])) {
		console.warn(`Key "${sortedBy}" does not exist in the data object.`);
		return data;
	}

	const sortDirection = direction === 'asc' ? 1 : -1;

	return data.toSorted((a, b) => {
		if (a[key] < b[key]) {
			return -1 * sortDirection;
		}
		if (a[key] > b[key]) {
			return 1 * sortDirection;
		}
		return 0;
	});
};
