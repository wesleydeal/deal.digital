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
let searchShortcutRegistry = {};

function toggleSearch(event=null, force=false) {
	searchCtr = elid("search-container");
	if (!searchCtr) {
		document.body.insertAdjacentHTML('afterbegin', `
			<div id="search-container">
				<button id="search-close" aria-label="Close Navigator">🗙</button>
				<label for="search-box"><b>///// Navigator</b> <i>alpha one</i></label>
				<input id="search-box" type="text" placeholder="Type to search 🧭">
				<menu id="search-results">
				</menu>
				<div id="search-help">
					<p>Press <kbd>/</kbd> to open and <kbd>Esc</kbd> to close.
					<h2>!keywords</h2>
					<ul id="search-keyword-list">
						<li>Brave Search <samp>b</samp>
						<li>Google Search <samp>g</samp>
						<li>ChatGPT Search <samp>gpt</samp>
						<li>YouTube <samp>yt</samp>
						<li>eBay <samp>eb</samp>
						<li>Amazon <samp>am</samp>
						<li>Mozilla Dev <samp>mdn</samp>
						<li>Zola Docs <samp>zola</samp>
						<li>Tera Docs <samp>tera</samp>
						<li>Anna's Archive <samp>an</samp>
						<li>Mapple TV <samp>tv</samp>
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
		elid("search-close").addEventListener("click", () => elid("search-container").remove());
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
				return s.color !== '' || query.search('rand') > -1;
			},
			action: (event) => {
				let color = event.target.title;
				if (color.search('rand') > -1) {
					color = "#" + Array.from(crypto.getRandomValues(new Uint8Array(3))).map(b => b.toString(16).padStart(2, '0')).join('');
					console.log("Randomly selected", color);
				}
				document.documentElement.style.setProperty('--color-primary', color);
			}
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
			getURL: (query) => 'https://youtube.com/results?search_query=' + encodeURIComponent(query).replaceAll("%20", "+"),
			validate: (query) => true,
			action: null,
		},
		ChatGPT: {
			keywords: ['gpt', 'chatgpt'],
			desc: 'ChatGPT Search',
			icon: '/icons/search/chatgpt.png',
			suggest: true,
			getURL: (query) => 'https://chatgpt.com/?q=' + encodeURIComponent(query).replaceAll("%20", "+"),
			validate: (query) => true,
			action: null,
		},
		Amazon: {
			desc: 'Amazon',
			icon: '/icons/search/amazon.png',
			suggest: true,
			getURL: (query) => 'https://www.amazon.com/s?k=' + encodeURIComponent(query).replaceAll("%20", "+"),
			validate: (query) => true,
			action: null,
		},
		MDN: {
			desc: 'Mozilla Dev',
			icon: '/icons/search/mdn.png',
			suggest: false,
			getURL: (query) => 'https://developer.mozilla.org/en-US/search?q='  + encodeURIComponent(query).replaceAll("%20", "+"),
			validate: (query) => true,
			action: null,
		},
		AnnasArchive: {
			desc: "Anna's Archive",
			icon: '/icons/search/annas-archive.png',
			suggest: false,
			getURL: (query) => 'https://annas-archive.org/search?q='  + encodeURIComponent(query).replaceAll("%20", "+"),
			validate: (query) => true,
			action: null,
		},
		Zola: {
			desc: 'Zola Documentation',
			icon: '/icons/search/zola.png',
			suggest: false,
			getURL: (query) => 'https://search.brave.com/search?q=site%3Agetzola.org+' + encodeURIComponent(query).replaceAll("%20", "+"),
			validate: (query) => true,
			action: null,
		},
		mappletv: {
			desc: 'Mapple.tv',
			icon: '/icons/search/mapple.png',
			suggest: false,
			getURL: (query) => 'https://mapple.tv/search?q=' + encodeURIComponent(query).replaceAll("%20", "+"),
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
		"color": "SetColor",
		"am": "Amazon",
		"mdn": "MDN",
		"an": "AnnasArchive",
		"anna": "AnnasArchive",
		"annasarchive": "AnnasArchive",
		"book": "AnnasArchive",
		"zo": "Zola",
		"zola": "Zola",
		"tv": "mappletv",
		"mapple": "mappletv",
		"mappletv": "mappletv",
	};
	let query = searchBox.value;
	let resultCount = 0;
	let foundKeyword = false;
	searchShortcutRegistry = {}

	while(query[0] === "/") {
		query = query.substring(1);
	}

	let entries = []; //todo replace
	let providerQueries = [];
	searchResults.innerHTML = '';

	if (query.replace(" ","") === "") {
		return;
	}

	// TODO: handle !bangs and searches beginning with provider queries
	for (word of query.toLowerCase().split(" ").reverse()) {
		if (word.search("!") > -1) {
			let possibleKeyword = word.substring(word.search("!") + 1);
			if (Array.from(Object.keys(keywordMap)).includes(possibleKeyword)) {
				providerQueries.push({ providerName: keywordMap[possibleKeyword], query: query.replace(/!.*?( |$)/g, '') });
				foundKeyword = true;
			}
		}
	}

	if (Array.from(Object.keys(keywordMap)).includes(query.toLowerCase().split(" ")[0])) {
		providerQueries.push({
			providerName: keywordMap[query.toLowerCase().split(" ")[0]],
			query: query.search(" ") > -1 ? query.substring(query.search(" ")) : '',
		});
		//foundKeyword = true;
	}

	// TODO: add options for all currently nonexistent providers
	for (providerName in providers) {
		if (providers[providerName].validate(query) && providers[providerName].suggest){
			providerQueries.push({ providerName: providerName, query: query });
		}
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
		if (p.action) {
			elQueryLink.addEventListener('click', p.action);
			elQueryLink.title = q;
			elQueryLink.href = "#";
		} else {
			elQueryLink.href = p.getURL(q);
		}
		elQueryLink.innerHTML = '<b>' + p.desc + '</b>: ' + q;
		elQueryLink.classList.add("search-link");
		if (foundKeyword) elQueryLink.classList.add("instant");
		elResult.insertAdjacentElement("beforeend", elQueryLink);

		if (resultCount < 11) {
			let elShortcut = document.createElement("div");
			elShortcut.classList.add("search-shortcut");
			if (resultCount === 0) {
				elShortcut.innerHTML = '<kbd>ENTER</kbd>';
				searchShortcutRegistry['Enter'] = elQueryLink;
			} else {
				elShortcut.innerHTML = '<kbd>Alt</kbd> + <kbd>Shift</kbd> + <kbd>' + resultCount + '</kbd>';
				searchShortcutRegistry['Alt' + " !@#$%^&*()"[resultCount]] = elQueryLink;
			}
			elResult.insertAdjacentElement("beforeend", elShortcut);
		}

		searchResults.insertAdjacentElement("beforeend", elResult);
		resultCount++;
	}

	// TODO: implement search this site
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
			toggleSearch(null, true);
		} else if (e.key === 'Escape') {
			const searchBox = elid("search-box");
			if (searchBox.value === "") {
				elid("search-container").remove();
			} else {
				searchBox.value = "";
				updateSearch(null);
			}
		}

		if (document.activeElement === elid("search-box")) {
			let key = "Alt".repeat(e.altKey) + e.key;
			if (key in searchShortcutRegistry) {
				elQueryLink = searchShortcutRegistry[key];
				if (e.ctrlKey) {
					let original = elQueryLink.target;
					elQueryLink.target = "_blank";
					elQueryLink.click();
					elQueryLink.target = original;
				} else {
					elQueryLink.click();
				}
			}
			if (e.key === "ArrowUp") {

			}
		}
	});

	if (document.location.search.search("q=") > -1) {
		toggleSearch(null, true);
		searchString = document.location.search.split("q=")[1].split("#")[0].split("&")[0].replaceAll("+", "%20");
		elid("search-box").value = decodeURIComponent(searchString);
		updateSearch(null);
		document.querySelectorAll("a.search-link.instant")?.[0].click?.();
	}
}



load();
