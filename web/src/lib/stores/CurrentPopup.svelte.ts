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
		} else {
			this.currentPopup = popup;
		}
	}
	clear(): void {
		this.currentPopup = '';
	}
}

const PopupStore = new PopupStoreClass();

export default PopupStore;
