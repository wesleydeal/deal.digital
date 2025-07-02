// GLOBAL VARIABLES --------------------------------
const content = document.getElementsByClassName('content')[0];
const root = document.documentElement;


// UTILITY FUNCTIONS -------------------------------

function el(id) {
	return document.getElementById(id);
}

function load() {
	el("btn_min_tools")?.addEventListener("click", () => {
		let classes = el("ui-tools").classList;
		if (classes.contains("min")) classes.remove("min");
		else classes.add("min");
	});
	el("btn_larger")?.addEventListener("click", () => {
		currentSize = getComputedStyle(content).fontSize;
		content.style.fontSize = 'calc(' + currentSize + ' * 1.0667)';
	});
	el("btn_smaller")?.addEventListener("click", () => {
		currentSize = getComputedStyle(content).fontSize;
		content.style.fontSize = 'calc(' + currentSize + ' * 0.937)';
	});
	el("btn_theme")?.addEventListener("click", () => {
		let currentDarkMode = getComputedStyle(root).getPropertyValue('--dark-mode') == 'true';
		console.log(currentDarkMode);
		root.style.setProperty('--dark-mode', !currentDarkMode);
	});
	el("btn_toc")?.addEventListener("click", () => {
		tocdetails = document.querySelector('#toc details');
		tocdetails.open = (tocdetails.open) ? false : true;
	});
	el("btn_top")?.addEventListener("click", () => {
		window.scrollTo({top: 0, behavior: 'smooth'});
	})
}

load();
