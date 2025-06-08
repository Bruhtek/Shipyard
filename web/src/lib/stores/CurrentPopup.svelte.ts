class PopupStoreClass {
	private currentPopup: string = $state('');

	constructor() {}

	get popup(): string {
		return this.currentPopup;
	}
	set popup(popup: string) {
		this.currentPopup = popup;
	}

	toggle(popup: string): void {
		if (this.currentPopup === popup) {
			this.currentPopup = '';
			console.debug('Popup closed:', popup);
		} else {
			this.currentPopup = popup;
			console.debug('Popup opened:', popup);
		}
	}
	clear(): void {
		this.currentPopup = '';
		console.debug('Popup cleared');
	}
}

const PopupStore = new PopupStoreClass();

export default PopupStore;
