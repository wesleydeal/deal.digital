// GLOBAL VARIABLES --------------------------------
const content = document.getElementsByClassName('content')[0];
const root = document.documentElement;

// UTILITY FUNCTIONS -------------------------------
const el = (selector) => document.querySelector(selector);
const elid = (id) => document.getElementById(id);

function load() {
	elid("btn_larger")?.addEventListener("click", () => {
		currentSize = getComputedStyle(content).fontSize;
		content.style.fontSize = 'calc(' + currentSize + ' * 1.0667)';
	});
	elid("btn_smaller")?.addEventListener("click", () => {
		currentSize = getComputedStyle(content).fontSize;
		content.style.fontSize = 'calc(' + currentSize + ' * 0.937)';
	});
	elid("btn_theme")?.addEventListener("click", () => {
		let currentDarkMode = getComputedStyle(root).getPropertyValue('--dark-mode') == 'true';
		console.log(currentDarkMode);
		root.style.setProperty('--dark-mode', !currentDarkMode);
	});
	elid("btn_toc")?.addEventListener("click", () => {
		tocdetails = document.querySelector('#toc details');
		tocdetails.open = (tocdetails.open) ? false : true;
	});
	elid("btn_top")?.addEventListener("click", () => {
		window.scrollTo({top: 0, behavior: 'smooth'});
	})
}

load();
