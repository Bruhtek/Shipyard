import { dev } from '$app/environment';

const url = dev ? 'localhost:4000' : window.location.host;
export const URLPrefix = window.location.protocol + '//' + url;

export const toProperCase = (str: string) => {
	const words = str.split(/[\s_]+/);
	return words.map((word) => word.charAt(0).toUpperCase() + word.slice(1)).join(' ');
};
