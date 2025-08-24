import type { Component } from 'svelte';

export type ActionData = {
	hoverBackground: string;
	hoverColor: string;
	text: string;
	icon: Component;
	onClick: (id: string, name: string) => void;
	disabled?: boolean;
	singleOnly?: boolean;
};
