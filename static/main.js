// GLOBAL VARIABLES --------------------------------
const content = document.getElementsByClassName('content')[0];
const root = document.documentElement;

// UTILITY FUNCTIONS -------------------------------
const el = (selector) => document.querySelector(selector);
const elid = (id) => document.getElementById(id);
const inrange = (x, start, end) => (x >= start) && (x <= end);

// SEARCH ------------------------------------------
function toggleSearch(force=false) {
	searchCtr = elid("search-container");
	if (!searchCtr) {
		document.body.insertAdjacentHTML('afterbegin', `
			<div id="search-container">
				<input id="search-box" type="text" placeholder="Type to search">
				<menu id="search-results">
				</menu>
			</div>
		`);
		elid("search-box").focus();
		elid("search-box").addEventListener('keyup', updateSearch);
	} else if (force) {
		elid("search-box").focus();
	} else {
		elid("search-container").remove();
	}
}

function updateSearch(event) {
	const searchBox = elid("search-box");
	const searchResults = elid("search-results");
	let query = searchBox.value;

	while(query[0] === "/") {
		query = query.substring(1);
	}

	console.log("Query:", query);

	let entries = [];

	firstSpace = query.search(" ");

	if (inrange(firstSpace, 0, 3) || query[0] == "!") {
		switch(query.substring(0,firstSpace)) {
			case "g":
			case "!g":
				entries.push({
					type: "Google",
					desc: "Google Search",
					query: query.substring(firstSpace),
					url: "https://google.com/search?q=" + query.replace(" ", "+"),
				});
				break;
			case "yt":
				entries.push({
					type: "YouTube",
					desc: "YouTube Search",
					query: query.substring(firstSpace),
					url: "https://google.com/search?q=" + query.replace(" ", "+"),
				});
				break;
		}
	}

	let htmlresults = "";

	for (entry of entries) {
		let htmlresult = '<div class="search-result">';
		htmlresult += '<img class="search-icon" src="icons/search/' + entry.type.toLowerCase() + '.png">';
		htmlresult += '<span class="search-query">' + entry.desc + '</span>';
		htmlresult += '</div>'
		htmlresults += htmlresult;
	}

	searchResults.innerHTML = htmlresults;

	console.log(entries);
}



// ON LOAD -----------------------------------------
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
		window.scrollTo({ top: 0, behavior: 'smooth' });
	});
	elid("btn_search")?.addEventListener("click", toggleSearch);
}

document.addEventListener("keyup", (e) => {
	if (e.key === '/') {
		toggleSearch(1);
	}
});

load();
