// GLOBAL VARIABLES --------------------------------
const content = document.getElementsByClassName('content')[0];
const root = document.documentElement;
let ddIndex = null;
let Fuse = null;
let fuse = null;

// UTILITY FUNCTIONS -------------------------------
const elid = (id) => document.getElementById(id);
const inrange = (x, start, end) => (x >= start) && (x <= end);
const tru = () => true;
function uniq(a) {
    var seen = {};
    return a.filter((item) => seen.hasOwnProperty(item) ? false : (seen[item] = true));
}

// SEARCH ------------------------------------------
let searchShortcutRegistry = {};
const providers = {
	/*
	each top-level key represents the internal ID of a provider, with values as follows:
	- keywords (opt) - array of strings, a list of words which prioritize this query if used as the first word of a query or anywhere prefixed with a !bang
	- desc (req) - string, the displayed name of the search provider
	- icon (opt, unused) - string, a path to a square icon representing the search provider
	- suggestIf (opt) - function(string) returning boolean, determines whether to show as a suggested provider
	- getURL (req if no action) - function(string) returning string, convert the search terms into a search link from this provider
	- action (req if no getURL) - function(event), JS function to run on page if the link is clicked (for example, set the page color or dark/light mode)
	TODO: prioritizeIf (opt) - function(string) returning boolean, determines whether this is a particularly good suggestion for this query and should be placed near the top
	TODO: color (opt) - string, valid CSS color associated with this provider
	*/
	Brave: {
		keywords: ["b", "br", "brave"],
		desc: 'Brave Search',
		icon: '/icons/search/brave.png',
		getURL: (query) => 'https://search.brave.com/search?q=' + encodeURIComponent(query).replaceAll("%20", "+"),
		suggestIf: tru,
		color: '#f50',
	},
	ChatGPT: {
		keywords: ['gpt', 'chatgpt'],
		desc: 'ChatGPT Search',
		icon: '/icons/search/chatgpt.png',
		getURL: (query) => 'https://chatgpt.com/?q=' + encodeURIComponent(query).replaceAll("%20", "+"),
		suggestIf: tru,
		color: '#74AA9C',
	},
	Google: {
		keywords: ['g', 'google'],
		desc: 'Google Search',
		icon: '/icons/search/google.png',
		getURL: (query) => "https://google.com/search?q=" + encodeURIComponent(query).replaceAll("%20", "+"),
		suggestIf: tru,
		color: '#1368F4',
	},
	eBay: {
		keywords: ['e', 'eb', 'ebay'],
		desc: 'eBay',
		icon: '/icons/search/ebay.png',
		getURL: (query) => "https://www.ebay.com/sch/i.html?_nkw=" + encodeURIComponent(query).replaceAll("%20", "+"),
		suggestIf: tru,
		color: '#3665f3',
	},
	YouTube: {
		keywords: ['y', 'yt', 'youtube'],
		desc: 'YouTube',
		icon: '/icons/search/youtube.png',
		getURL: (query) => 'https://youtube.com/results?search_query=' + encodeURIComponent(query).replaceAll("%20", "+"),
		suggestIf: tru,
		color: '#f00',
	},
	Amazon: {
		keywords: ['am', 'amazon', 'amzn'],
		desc: 'Amazon',
		icon: '/icons/search/amazon.png',
		getURL: (query) => 'https://www.amazon.com/s?k=' + encodeURIComponent(query).replaceAll("%20", "+"),
		suggestIf: tru,
		color: '#f90',
	},
	MDN: {
		keywords: ['mdn'],
		desc: 'Mozilla Dev',
		icon: '/icons/search/mdn.png',
		getURL: (query) => 'https://developer.mozilla.org/en-US/search?q=' + encodeURIComponent(query).replaceAll("%20", "+"),
		color: '#8cb4ff',
	},
	AnnasArchive: {
		keywords: ['an','anna','annas','annasarchive','book'],
		desc: "Anna's Archive",
		icon: '/icons/search/annas-archive.png',
		getURL: (query) => 'https://annas-archive.org/search?q=' + encodeURIComponent(query).replaceAll("%20", "+"),
		color: '#0195ff',
	},
	Zola: {
		keywords: ['zola'],
		desc: 'Zola Docs',
		icon: '/icons/search/zola.png',
		getURL: (query) => 'https://search.brave.com/search?q=site%3Agetzola.org+' + encodeURIComponent(query).replaceAll("%20", "+"),
		color: '#191919',
	},
	Tera: {
		keywords: ['tera'],
		desc: 'Tera Docs',
		getURL: (query) => 'https://search.brave.com/search?q=site%3Ahttps%3A%2F%2Fkeats.github.io%2Ftera%2Fdocs%2F+' + encodeURIComponent(query).replaceAll("%20", "+"),
		color: '#de6262',
	},
	mappletv: {
		keywords: ['tv','mapple','mapple.tv'],
		desc: 'Mapple.TV',
		icon: '/icons/search/mapple.png',
		getURL: (query) => 'https://mapple.tv/search?q=' + encodeURIComponent(query).replaceAll("%20", "+"),
		color: '#fff',
	},
	SetColor: {
		keywords: ['color'],
		desc: 'Set Site Color',
		suggestIf: (query) => {
			let s = new Option().style;
			s.color = query;
			return s.color !== '' || query.search('rand') > -1;
		},
		action: (event) => {
			let color = event.target.title;
			if (color.search('rand') > -1) {
				color = "#" + Array.from(crypto.getRandomValues(new Uint8Array(3))).map(b => b.toString(16).padStart(2, '0')).join('');
				console.log("Randomly selected", color);
			}
			document.documentElement.style.setProperty('--color-primary', color);
			closeSearch();
		}
	},
	toggleDark: {
		keywords: ['dark', 'light', 'toggle'],
		desc: 'Set Dark Mode',
		suggestIf: (q) => (['dark', 'light', 'toggle'].includes(q.toLowerCase().replaceAll(' ',''))),
		action: (event) => {
			let q = elid("search-box").value;
			if (q.includes('toggle')) {
				let currentDarkMode = getComputedStyle(root).getPropertyValue('--dark-mode') == 'true';
				root.style.setProperty('--dark-mode', !currentDarkMode);
			} else if (q.includes('unset') || q.includes('none')) {
				root.style.setProperty('--dark-mode', null);
			} else if (q.includes('light')) {
				root.style.setProperty('--dark-mode', false);
			} else {
				root.style.setProperty('--dark-mode', true);
			}
			closeSearch();
		},
	}
};
let keywordMap = {};
for (providerName in providers) {
	for (kw of providers[providerName]?.keywords) {
		keywordMap[kw] = providerName;
	}
}

function toggleSearch(query="") {
	if (elid("search-container")) {
		closeSearch();
	} else {
		openSearch(query);
	}
}

function openSearch(query=null) {
	const url = new URL(window.location.href);
	console.log(query);
	if (query.type == 'click' || query === null) {
		query = "";
	}
	if (!elid("search-container")) {
		document.body.insertAdjacentHTML('afterbegin', `
			<div id="search-container">
				<button id="search-close" aria-label="Close Navigator">X</button>
				<label for="search-box"><b>///// Navigator</b> <i>alpha one</i></label>
				<input id="search-box" type="text" placeholder="Type to search 🧭">
				<menu id="search-results">
				</menu>
				<div id="search-help">
					<p>Press <kbd>/</kbd> to open and <kbd>Esc</kbd> clear or close.
					<h2>!keywords</h2>
					<ul id="search-keyword-list">
					</ul>
				  <h2>Examples</h2>
					<ul>
						<li><a href="#" onclick="document.documentElement.style.setProperty('--color-primary', 'aquamarine')">color aquamarine</a>
						<li><a href="https://www.ebay.com/sch/i.html?_nkw=+ibm+model+m+(bolt%2Cscrew)+(mod%2Cmodded)">eb ibm model m (bolt,screw) (mod,modded)</a>
						<li><a href="https://youtube.com/results?search_query=+moments+with+heavy+french+toast">yt moments with heavy french toast</a>
						<li><a href="https://chatgpt.com/?q=where+can+I+get+a+good+asada+burrito+nearby+">where can I get a good asada burrito nearby !gpt</a>
						<li><a href="https://annas-archive.org/search?q=mike+ma">mike ma !an</a>
					</ul>
				</div>
			</div>
		`);
		const keywordList = elid("search-keyword-list");
		for (providerName in providers) {
			if (providers[providerName].keywords) {
				const keywordEntry = document.createElement("li");
				keywordEntry.textContent = providers[providerName].desc;
				const keywordSamp = document.createElement("samp");
				keywordSamp.textContent = providers[providerName].keywords[0];
				keywordEntry.insertAdjacentElement("beforeend", keywordSamp);
				keywordEntry.addEventListener("click", () => openSearch(keywordSamp.textContent + " "));
				if (providers[providerName]?.color) {
					keywordEntry.style.setProperty('--c', providers[providerName]?.color);
				}
				keywordList.insertAdjacentElement("beforeend", keywordEntry);
			}
		}
		if (!url.searchParams.get("q")) {
			url.searchParams.set('q', query);
			window.history.replaceState(null, null, url);
		}
		elid("search-close").addEventListener("click", closeSearch);
		elid("search-box").focus();
		elid("search-box").addEventListener('keyup', updateSearch);
		if (query !== null) {
			elid("search-box").value = query;
			updateSearch();
		}
	} else {
		elid("search-box").placeholder = elid("search-box").value;
		elid("search-box").value = query;
	}
	if (query !== null) {
		elid("search-box").value = query;
		updateSearch();
	}
	elid("search-box").focus();
	url.searchParams.set('q', query);
	window.history.replaceState(null, null, url);
}

function closeSearch() {
	elid("search-container").remove();
	const url = new URL(window.location.href);
	url.searchParams.delete('q');
	window.history.replaceState(null, null, url);
}

async function updateSearch(event=null) {
	let startTime = Date.now();
	const searchBox = elid("search-box");
	const searchResults = elid("search-results");
	const url = new URL(window.location.href);

	let query = searchBox.value;
	let resultCount = 0;
	let foundKeyword = false;
	searchShortcutRegistry = {}

	while(query[0] === "/") {
		query = query.substring(1);
	}

	let providerQueries = [];
	let localResults = [];
	searchResults.innerHTML = '';

	url.searchParams.set('q', query);
	window.history.replaceState(null, null, url);

	if (query.replaceAll(" ","") === "") {
		return;
	}

	// prioritize providers with !bangs
	for (word of query.toLowerCase().split(" ").reverse()) {
		if (word.includes("!")) {
			let possibleKeyword = word.substring(word.search("!") + 1);
			if (Object.keys(keywordMap).includes(possibleKeyword)) {
				providerQueries.push({ providerName: keywordMap[possibleKeyword], query: query.replace(/!.*?( |$)/g, '') });
				foundKeyword = true;
			}
		}
	}

	// prioritize providers with a bangless keyword at the beginning
	if (Object.keys(keywordMap).includes(query.toLowerCase().split(" ")[0])) {
		providerQueries.push({
			providerName: keywordMap[query.toLowerCase().split(" ")[0]],
			query: query.search(" ") > -1 ? query.substring(query.search(" ")) : '',
		});
	}

	// add options for all currently nonexistent providers
	let ignoreProviders = providerQueries.map((p) => p.providerName);
	for (providerName in providers) {
		if ((providers[providerName].suggestIf?.(query) ?? false) && !(ignoreProviders.includes(providerName))){
			providerQueries.push({ providerName: providerName, query: query });
		}
	}

	// handle local queries
	if (query[0] === "&") {
		if (!ddIndex || !Fuse || !fuse) {
			ddIndex = (await import("/search_index.en.json", { with: { type: "json" } })).default;
			Fuse = (await import("/fuse.min.mjs")).default;

			const fuseOptions = {
				isCaseSensitive: false,
				includeScore: true,
				ignoreDiacritics: true,
				shouldSort: true,
				includeMatches: false,
				findAllMatches: false,
				minMatchCharLength: 2,
				location: 0,
				threshold: 0.2,
				distance: 100,
				useExtendedSearch: true,
				ignoreLocation: true,
				ignoreFieldNorm: false,
				fieldNormWeight: 1,
				keys: [
					{ name: "title", weight: 1 },
					{ name: "url", weight: 1 },
					{ name: "body", weight: 1 },
					{ name: "description", weight: 1 },
				]
			};

			fuse = new Fuse(ddIndex, fuseOptions);
		}

		localResults = fuse.search(query.substring(1));
	}

	for (result of localResults) {
		let elResult = document.createElement("div");
		elResult.id = "search-result-" + resultCount;
		elResult.classList.add("search-result", "search-result-local");

		let elResultLink = document.createElement("a");
		elResultLink.classList.add("search-link");
		elResultLink.innerHTML = result.item.title;
		elResultLink.href = result.item.url;

		elResult.insertAdjacentElement("beforeend", elResultLink);

		if (resultCount < 11) {
			let elShortcut = document.createElement("div");
			elShortcut.classList.add("search-shortcut");
			if (resultCount === 0) {
				elShortcut.innerHTML = '<kbd>ENTER</kbd>';
				searchShortcutRegistry['Enter'] = elResultLink;
			} else {
				elShortcut.innerHTML = '<kbd>Alt</kbd> <kbd>' + resultCount + '</kbd>';
				searchShortcutRegistry['Alt' + resultCount] = elResultLink;
			}
			elResult.insertAdjacentElement("beforeend", elShortcut);
		}

		elResult.addEventListener("click", () => elResultLink.click());
		searchResults.insertAdjacentElement("beforeend", elResult);

		resultCount++;
	}

	for (item of providerQueries) {
		const p = providers[item.providerName];
		const q = item.query;

		let elResult = document.createElement("div");
		elResult.id = "search-result-" + resultCount;
		elResult.classList.add("search-result");

		//let elIcon = document.createElement("img");
		//elIcon.src = p.icon;
		//elResult.insertAdjacentElement("beforeend", elIcon);

		let elQueryLink = document.createElement("a");
		if (p?.action) {
			elQueryLink.addEventListener('click', p.action);
			elQueryLink.title = q;
			elQueryLink.href = "javascript:;";
		} else {
			elQueryLink.href = p?.getURL(q);
		}
		elQueryLink.innerHTML = '<b>' + p.desc + '</b>: ' + q;
		elQueryLink.classList.add("search-link");
		if (foundKeyword) elQueryLink.classList.add("instant");
		elResult.insertAdjacentElement("beforeend", elQueryLink);

		elResult.addEventListener("click", () => elQueryLink.click());

		if (resultCount < 11) {
			let elShortcut = document.createElement("div");
			elShortcut.classList.add("search-shortcut");
			if (resultCount === 0) {
				elShortcut.innerHTML = '<kbd>ENTER</kbd>';
				searchShortcutRegistry['Enter'] = elQueryLink;
			} else {
				elShortcut.innerHTML = '<kbd>Alt</kbd> <kbd>' + resultCount + '</kbd>';
				searchShortcutRegistry['Alt' + resultCount] = elQueryLink;
			}
			elResult.insertAdjacentElement("beforeend", elShortcut);
		}

		searchResults.insertAdjacentElement("beforeend", elResult);
		resultCount++;
	}

	url.searchParams.set('q', query);
	window.history.replaceState(null, null, url);

	elid("search-results").insertAdjacentHTML("beforeend", "<p>Retrieved in " + (Date.now() - startTime) + "ms");
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
	elid("search-link")?.addEventListener("click", toggleSearch);

	document.addEventListener("keyup", (e) => {
		if (e.key === '/' && document.activeElement.tagName != "INPUT") {
			openSearch(window.getSelection().toString().replaceAll('\n',''));
		} else if (e.key === 'Escape') {
			if (elid("search-box") && elid("search-box")?.value.replaceAll(' ','') === "") {
				closeSearch();
			} else {
				openSearch("");
			}
		}
	});
	document.addEventListener("keydown", (e) => {
		if (document.activeElement === elid("search-box")) {
			let key = "Alt".repeat(e.altKey) + e.key;
			if (key in searchShortcutRegistry) {
				e.preventDefault();
				const link = searchShortcutRegistry[key];
				if (e.ctrlKey) { // open in new tab on CTRL
					const original = link.target;
					link.target = "_blank";
					link.click();
					link.target = original;
				} else {
					link.click();
				}
			}
			if (e.key === "ArrowUp") {

			}
		}
	});

	if (document.location.search.search("q=") > -1) {
		const url = new URL(window.location.href);
		searchString = url.searchParams.get("q");
		openSearch(searchString);
		updateSearch();
		if (window.performance.getEntriesByType("navigation")[0].type === "navigate") { // only click instant links when not back/forward/reload
			document.querySelectorAll("a.search-link.instant")?.[0]?.click?.();
		}
	}
}



load();
