export type UserPreferences = {
	autoExpandActions: boolean;
	autoCollapseCompletedActions: boolean;
	keepDismissedActions: boolean;

	terminalWidth: number;
	terminalHeight: number;
};

type MapUserPreferences<T> = {
	[K in keyof T]: {
		default: T[K];
		label: string;
		tip?: string;
		limits?: { min?: number; max?: number }; // Only for number types, obviously
	};
};

type UserPreferencesMetadata = MapUserPreferences<UserPreferences>;
export const userPreferencesMetadata: UserPreferencesMetadata = {
	autoExpandActions: {
		default: true,
		label: 'Automatically expand new actions'
	},
	autoCollapseCompletedActions: {
		default: true,
		label: 'Automatically minimize completed actions'
	},
	keepDismissedActions: {
		default: true,
		label: 'Keep dismissed actions in the action list',
		tip: 'This only applies to expanded actions. Collapsed actions are always discarded.'
	},
	terminalWidth: {
		default: 80,
		label: 'Width of the terminal (in characters)',
		tip: 'Not recommended to change this to less than 80 - may cause glitches in rendering',
		limits: { min: 20, max: 200 }
	},
	terminalHeight: {
		default: 12,
		label: 'Height of the terminal (in characters)',
		limits: { min: 1, max: 30 }
	}
};

const defaultPreferences: UserPreferences = Object.fromEntries(
	Object.entries(userPreferencesMetadata).map(([key, value]) => [key, value.default])
) as UserPreferences;

class UserPreferencesStoreClass {
	p: UserPreferences = $state(defaultPreferences);

	constructor() {
		this.loadPreferences();
	}

	setPreference<K extends keyof UserPreferences>(key: K, value: UserPreferences[K]) {
		this.p[key] = value;
		this.persistPreferences();
	}

	persistPreferences() {
		localStorage.setItem('userPreferences', JSON.stringify(this.p));
	}

	loadPreferences() {
		const stored = localStorage.getItem('userPreferences');
		if (stored) {
			try {
				const parsed = JSON.parse(stored);
				this.p = { ...defaultPreferences, ...parsed };
			} catch (e) {
				console.error('Failed to parse user preferences from localStorage', e);
				this.p = defaultPreferences;
			}
		} else {
			this.p = defaultPreferences;
		}
	}
}

const UserPreferencesStore = new UserPreferencesStoreClass();
export default UserPreferencesStore;
