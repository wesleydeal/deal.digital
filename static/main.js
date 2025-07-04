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
		SetColor: {
			desc: 'Set Primary Color',
			icon: null,
			suggest: true,
			getURL: (query) => '',
			validate: (query) => {
				let s = new Option().style;
				s.color = query;
				return s.color !== '';
			},
			action: () => document.documentElement.style.setProperty('--color-primary', document.getElementById("search-box").value),
		},
		Google: {
			keywords: ['g', 'google'],
			desc: 'Google Search',
			icon: '/icons/search/google.png',
			suggest: true,
			getURL: (query) => "https://google.com/search?q=" + encodeURIComponent(query).replaceAll("%20", "+"),
			validate: (query) => true,
			action: null,
		},
		eBay: {
			keywords: ['e', 'eb', 'ebay'],
			desc: 'eBay',
			icon: '/icons/search/ebay.png',
			suggest: true,
			getURL: (query) => "https://www.ebay.com/sch/i.html?_nkw=" + encodeURIComponent(query).replaceAll("%20", "+"),
			validate: (query) => true,
			action: null,
		},
		YouTube: {
			keywords: ['y', 'yt', 'youtube'],
			desc: 'YouTube',
			icon: '/icons/search/youtube.png',
			suggest: true,
			getURL: (query) => '',
			validate: (query) => true,
			action: null,
		},
		ChatGPT: {
			keywords: ['gpt', 'chatgpt'],
			desc: 'ChatGPT Search',
			icon: '/icons/seach/chatgpt.png',
			suggest: true,
			getURL: (query) => '',
			validate: (query) => true,
			action: null,
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
	let providerQueries = [];
	searchResults.innerHTML = '';

	// TODO: handle !bangs and searches beginning with provider queries
	for (word of query.split(" ").reverse()) {
		if (word.search("!") > -1) {
			let possibleKeyword = word.substring(word.search("!") + 1);
			if (Array.from(Object.keys(keywordMap)).includes(possibleKeyword)) {
				providerQueries.push({ providerName: keywordMap[possibleKeyword], query: query.replace(/!.*?( |$)/g, '') });
			}
		}
	}

	if (Array.from(Object.keys(keywordMap)).includes(query.split(" ")[0])) {
		providerQueries.push({ providerName: keywordMap[query.split(" ")[0]], query: query.substring(query.search(" ")) });
	}

	// TODO: add options for all currently nonexistent providers
	for (providerName in providers) {
		if (providers[providerName].validate(query)){
			providerQueries.push({ providerName: providerName, query: query });
		}
	}

	for (item of providerQueries) {
		const p = providers[item.providerName];
		const q = item.query;
		console.log(p, q);
		let elResult = document.createElement("div");
		let elIcon = document.createElement("img");
		let elQueryLink = document.createElement("a");
		let elShortcut = document.createElement("div");
		elResult.id = "search-result-" + resultCount;
		elResult.classList.add("search-result");
		elIcon.src = p.icon;
		if (p.action) {
			elQueryLink.addEventListener('click', p.action);
			elQueryLink.href = "#";
		} else {
			elQueryLink.href = p.getURL(q);
		}
		elQueryLink.innerHTML = '<b>' + p.desc + '</b>: ' + q;
		elShortcut.classList.add("search-shortcut");
		elShortcut.innerHTML = (resultCount < 10 ?
				(resultCount > 0 ? '<kbd>Alt</kbd> + <kbd>' + resultCount + '</kbd>' : '<kbd>ENTER</kbd>')
				: '');
		//elResult.insertAdjacentElement("beforeend", elIcon);
		elResult.insertAdjacentElement("beforeend", elQueryLink);
		elResult.insertAdjacentElement("beforeend", elShortcut);
		searchResults.insertAdjacentElement("beforeend", elResult);
		resultCount++;
	}

	// TODO: implement search this site

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
