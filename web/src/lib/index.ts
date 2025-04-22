import { dev } from '$app/environment';

const url = dev ? 'localhost:4000' : window.location.host;
export const URLPrefix = window.location.protocol + "//" + url