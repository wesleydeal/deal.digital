// GLOBAL VARIABLES --------------------------------
const content = document.getElementsByClassName('content')[0];
const root = document.documentElement;

// UTILITY FUNCTIONS -------------------------------
const el = (selector) => document.querySelector(selector);
const elid = (id) => document.getElementById(id);
const inrange = (x, start, end) => (x >= start) && (x <= end);
function uniq(a) {
    var seen = {};
    return a.filter((item) => seen.hasOwnProperty(item) ? false : (seen[item] = true));
}

// SEARCH ------------------------------------------
function toggleSearch(event=null, force=false) {
	searchCtr = elid("search-container");
	if (!searchCtr) {
		document.body.insertAdjacentHTML('afterbegin', `
			<div id="search-container">
				<p>This feature is unfinished and broken.</p>
				<label for="search-box"><b>🧭 Navigator</b>: Press <kbd>/</kbd>, type your query, and press <kbd>ENTER</kbd> or select a provider.</label>
				<input id="search-box" type="text" placeholder="Type to search">
				<menu id="search-results">
				</menu>
			</div>
		`);
		elid("search-box").focus();
		elid("search-box").addEventListener('keyup', updateSearch);
	} else if (force) {
		elid("search-box").placeholder = elid("search-box").value;
		elid("search-box").value = "";
		elid("search-box").focus();
	} else {
		elid("search-container").remove();
	}
}

function updateSearch(event) {
	const searchBox = elid("search-box");
	const searchResults = elid("search-results");
	const providers = {
		Google: {
			keywords: ['g', 'google'],
			desc: 'Google Search',
			icon: '/icons/search/google.png',
			getURL: (query) => "https://google.com/search?q=" + encodeURIComponent(query).replaceAll("%20", "+"),
		},
		eBay: {
			keywords: ['e', 'eb', 'ebay'],
			desc: 'eBay',
			icon: '/icons/search/ebay.png',
			getURL: (query) => "https://www.ebay.com/sch/i.html?_nkw=" + encodeURIComponent(query).replaceAll("%20", "+"),
		},
		YouTube: {
			keywords: ['y', 'yt', 'youtube'],
			desc: 'YouTube',
			icon: '/icons/search/youtube.png',
			getURL: (query) => '',
		},
		ChatGPT: {
			keywords: ['gpt', 'chatgpt'],
			desc: 'ChatGPT Search',
			icon: '/icons/seach/chatgpt.png',
			getURL: (query) => '',
		}
	}
	const keywordMap = {
		"g": "Google",
		"google": "Google",
		"e": "eBay",
		"eb": "eBay",
		"ebay": "eBay",
		"y": "YouTube",
		"yt": "YouTube",
		"gpt": "ChatGPT",
		"chatgpt": "ChatGPT",
	};
	let query = searchBox.value;
	let htmlresults = "";
	let resultCount = 0;

	while(query[0] === "/") {
		query = query.substring(1);
	}

	console.log("Query:", query);

	let entries = []; //todo replace
	let selectedProviders = [];

	// TODO: handle !bangs and searches beginning with provider queries
	for (word of query.split(" ").reverse()) {
		if (word.search("!") > -1) {
			let possibleKeyword = word.substring(word.search("!") + 1);
			if (Array.from(Object.keys(keywordMap)).includes(possibleKeyword)) {
				selectedProviders.push(keywordMap[possibleKeyword]);
			}
		}
	}

	for (provider of selectedProviders) {
		const p = providers[provider];
		const q = query.replace(/!.*?( |$)/g, '');
		console.log(q);
		let htmlresult = '<div class="search-result" id="search-result-' + resultCount + '">' +
			'<img class="search-icon" src="' + p.icon + '">' +
			'<a class="search-query" href="' + p.getURL(q) + '">' +
			'<b>' + p.desc + '</b>: ' + q + '</a>' +
			(resultCount < 10 ?
				'<span class="search-shortcut"' + (resultCount > 0 ? '<kbd>Alt</kbd> + <kbd>' + resultCount +'</kbd>' : '<kbd>ENTER</kbd>') + '</span>'
				: '') +
			'</div>';
		htmlresults += htmlresult;
		resultCount++;
	}

	selectedProviders = [];

	if (Array.from(Object.keys(keywordMap)).includes(query.split(" ")[0])) {
		selectedProviders.push(keywordMap[query.split(" ")[0]]);
	}

	for (provider of selectedProviders) {
		const p = providers[provider];
		const q = query.substring(query.search(' '));
		let htmlresult = '<a href="' + p.getURL(q) + '">' +
			'<div class="search-result">';
		htmlresult += '<img class="search-icon" src="' + p.icon + '">';
		htmlresult += '<span class="search-query"><b>' + p.desc + '</b>: ' + q + '</span>';
		htmlresult += '</a></div>'
		htmlresults += htmlresult;
		resultCount++;
	}

	// TODO: implement search this site

	// TODO: add options for all currently nonexistent providers
	for (providerName in providers) {
		const p = providers[providerName];
		const q = query;
		let htmlresult = '<a href="' + p.getURL(q) + '">' +
			'<div class="search-result">';
		htmlresult += '<img class="search-icon" src="' + p.icon + '">';
		htmlresult += '<span class="search-query"><b>' + p.desc + '</b>: ' + q + '</span>';
		htmlresult += '</a></div>'
		htmlresults += htmlresult;
		resultCount++;
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
		toggleSearch(null, true);
	}
});

load();
