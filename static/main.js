// GLOBAL VARIABLES --------------------------------
const content = document.getElementsByClassName('content')[0];
const root = document.documentElement;
const imageExtensions = ["svg", "png", "jpg", "jfif", "gif", "webp", "avif"];
let ddIndex = null;
let Fuse = null;
let fuse = Object;
fuse.search = async (...args) => {await initFuse(); return fuse.search(...args)}
let searchWindowPos;
let searchDrag = [];
let savedPage;
let soundPlayers = {};

// UTILITY FUNCTIONS -------------------------------
const elid = (id) => document.getElementById(id);
const inrange = (x, start, end) => (x >= start) && (x <= end);
const tru = () => true;
function uniq(a) {
    var seen = {};
    return a.filter((item) => seen.hasOwnProperty(item) ? false : (seen[item] = true));
}
function playSound(url, volume = 1) {
	const start = new Date();
	if (!soundPlayers?.[url]) soundPlayers[url] = new Audio(url);
	soundPlayers[url].volume = volume;
	const playIt = () => {
		if ((new Date() - start) > 300) return;
		soundPlayers[url].currentTime = 0;
		soundPlayers[url].play();
	};
	if (soundPlayers[url].readyState >= HTMLMediaElement.HAVE_CURRENT_DATA) {
		playIt();
	} else {
		soundPlayers[url].addEventListener("canplaythrough", playIt, { once: true });
	}
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
	+ ONE of the following
		- getURL - function(string) returning string, convert the search terms into a search link from this provider
		- action - function(event), JS function to run on page if the link is clicked (for example, set the page color or dark/light mode)
		- getURLs - function(string) returning array[array[string, string], ...] where each element is a [name, URL]
		- TODO NOT IMPLEMENTED - getActions - function(string) returning array[array[string, function(event)], ...] where each element is a [name, action]
	- TODO prioritizeIf (opt) - function(string) returning boolean, determines whether this is a particularly good suggestion for this query and should be placed near the top
	- color (opt) - string, valid CSS color associated with this provider
	- hide (opt) - boolean, if true don't list in !keywords
	*/
	Local: {
		keywords: ['d','dd','deal','deal.digital'],
		desc: 'deal.digital',
		getURLs: async (query) => {
			r = await fuse.search(query);
			if (r?.[0]?.score < 0.05 && query.length > 3 && r?.[0].item.url != window.location.href.split('?')[0].split('#')[0]) {
				await replacePage(r[0].item.url);
			}
			return r.map((result) => [result.item.title, result.item.url]);
		},
		//suggestIf: async (query) => (await fuse.search(query))?.[0]?.score < 0.1,
	},
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
	WolframAlpha: {
		keywords: ['wa','wolfram','wolframalpha'],
		desc: 'Wolphram|Alpha',
		getURL: (query) => 'http://www.wolframalpha.com/input/?i=' + encodeURIComponent(query),
		color: '#ee1f22',
		suggestIf: (query) => ["+", "-", "*", "/", "convert", "per"].reduce((a, s) => a + query.includes(s), 0),
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
			root.style.setProperty('--color-primary', color);
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
	},
	style: {
		keywords: ['style', 'stylesheet'],
		desc: 'Set Stylesheet',
		hide: true,
		action: (event) => {
			let ssLink = document.querySelector("link[rel='stylesheet'][as='style']");
			ssLink.href = "/" + event.target.title.replaceAll(" ","") + ".css";
			closeSearch();
		},
	},
	load: {
		keywords: ['load'],
		desc: 'Load Content From Internal Page',
		hide: true,
		action: replacePage,
	},
};
let keywordMap = {};
for (providerName in providers) {
	for (kw of providers[providerName]?.keywords) {
		keywordMap[kw] = providerName;
	}
}

async function replacePage(path) {
	if (path instanceof Event) {
		path = elid("search-box").value.replaceAll(" !load", "").split(" ").at(-1);
	}
	const response = await fetch(path);
	if (!response.ok) return;

	const html = await response.text();
	const doc = (new DOMParser()).parseFromString(html, 'text/html');
	document.querySelector('section.content').replaceWith(doc.querySelector('section.content'));
	if (doc.documentElement.style) {
		root.setAttribute('style', doc.documentElement.getAttribute('style'));
	}

	document.title = doc.title;
	window.history.replaceState(null, null, response.url);

	if (window.location.hash) {
		window.scrollTo({top: document.querySelector(window.location.hash).getBoundingClientRect().top, behavior: 'smooth'});
	} else {
		window.scrollTo({top: 0, behavior: 'smooth'});
	}
}

function toggleSearch(query="") {
	if (elid("search-container")) {
		if (elid("search-container").classList.contains('min')) {
			elid("search-container").classList.remove('min');
			elid("search-box").focus();
			//playSound('/sounds/KDE_Window_Shade_Up.ogg');
		} else {
			closeSearch();
		}
	} else {
		openSearch(query);
	}
}

async function initFuse(){
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

function openSearch(query=null) {
	const url = new URL(window.location.href);
	if (query.type == 'click' || query === null) {
		query = "";
	}
	if (!elid("search-container")) {
		//playSound('/sounds/KDE_Window_UnHide.ogg');
		document.body.insertAdjacentHTML('afterbegin', `
			<div id="search-container">
				<div id="search-titlebar">
					<label for="search-box"><b>Navigator</b> <i>alpha one</i></label>
					<div class="window-buttons">
						<button id="search-min" aria-label="Minimize Navigator"></button>
						<button id="search-max" aria-label="Maximize Navigator"></button>
						<button id="search-close" aria-label="Close Navigator"></button>
					</div>
				</div>
				<div id="search-inner">
					<input id="search-box" type="text" placeholder="Type to search 🧭" autocomplete="off">
				</div>
			</div>
		`);
		elid("search-box").focus();
		elid('search-inner').insertAdjacentHTML('beforeend', `
			<menu id="search-results">
			</menu>
			<div id="search-help">
				<p>Press <kbd>/</kbd> to open and <kbd>Esc</kbd> to clear or close.
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
		`);
		const keywordList = elid("search-keyword-list");
		for (providerName in providers) {
			if (providers[providerName].keywords && !(providers[providerName]?.hide)) {
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
		elid("search-close").addEventListener("click", closeSearch);
		elid("search-min").addEventListener("click", minimizeSearch);
		elid("search-max").addEventListener("click", toggleMaximizeSearch);
		elid("search-titlebar").addEventListener("dblclick", toggleMaximizeSearch);
		elid("search-titlebar").addEventListener("pointerdown", dragSearchStart);
		elid("search-box").addEventListener('input', updateSearch);
	} else {
		elid("search-box").placeholder = elid("search-box").value;
	}

	elid("search-box").value = query;
	updateSearch();
	elid("search-box").focus();
	url.searchParams.set('q', query);
	window.history.replaceState(null, null, url);
}

function closeSearch() {
	const url = new URL(window.location.href);
	url.searchParams.delete('q');
	window.history.replaceState(null, null, url);
	elid("search-container").classList.add('hidden');
	setTimeout(() => {elid("search-container").remove()}, 300);
	//playSound('/sounds/KDE_Window_Close.ogg')
}

function minimizeSearch() {
	elid("search-container").classList.add('min');
	//playSound('/sounds/KDE_Window_Shade_Down.ogg');
}

function toggleMaximizeSearch() {
	const sC = elid("search-container");
	const style = getComputedStyle(sC);

	if (sC.classList.contains('max')) {
		sC.classList.remove('max');
		for (pName in searchWindowPos) {
			sC.style.setProperty(pName, searchWindowPos[pName]);
		}
	} else {
		searchWindowPos = {
			top: style.top,
			left: style.left,
			width: style.width,
			height: style.height,
		};
		sC.style.removeProperty('top');
		sC.style.removeProperty('left');
		sC.style.removeProperty('width');
		sC.style.removeProperty('height');
		sC.classList.add('max');
	}
}

function onDragRelease(event) {
	searchDrag = [];
	document.removeEventListener('pointermove', dragSearchUpdate, {passive: false});
	elid("search-container").classList.remove('drag');
	root.style.removeProperty('touch-action');
	document.removeEventListener('pointerup', onDragRelease, {passive: false});
}

function dragSearchStart(event) {
	const sC = elid("search-container");
	const style = getComputedStyle(sC);
	if (document.querySelector("#search-titlebar .window-buttons").contains(event.target)) return;
	if (event.buttons != 1) return;

	if (sC.classList.contains('max')) { // don't start dragging until we've moved at least 4 pixels, if maximized
		let originalX = event.clientX;
		let originalY = event.clientY;
		const moveListener = (event) => {
			if (event.buttons != 1) {
				sC.removeEventListener('pointermove', moveListener);
				return;
			}
			if (((originalX + event.clientX)**2 + (originalY + event.clientY)**2)**0.5 > 4) {
				sC.style.setProperty('width', searchWindowPos.width);
				sC.style.setProperty('height', 'auto');
				sC.style.setProperty('top', 0);
				sC.style.setProperty('left', (event.clientX / window.innerWidth) * (window.innerWidth - parseFloat(searchWindowPos.width)) + "px");
				sC.classList.remove('max');
				dragSearchStart(event);
				sC.removeEventListener('pointermove', moveListener);
			}
		}
		sC.addEventListener('pointermove', moveListener);
	} else if (!searchDrag || searchDrag.length < 1) {
		searchDrag = [event.clientX - sC.offsetLeft, event.clientY - sC.offsetTop];
		elid("search-container").classList.add('drag');
		document.addEventListener('pointermove', dragSearchUpdate, {passive: false});
		document.addEventListener('pointerup', onDragRelease, {passive: false});
	}
}

function dragSearchUpdate(event) {
	const sC = elid("search-container");
	sC.style.setProperty('right', 'unset');
	sC.style.setProperty('left', event.clientX - searchDrag[0] + 'px');
	sC.style.setProperty('top', event.clientY - searchDrag[1] + 'px');
}

async function updateSearch(event=null) {
	searchShortcutRegistry = {}
	let startTime = Date.now();
	const searchBox = elid("search-box");
	const searchResults = elid("search-results");
	const url = new URL(window.location.href);

	let query = searchBox.value;
	let resultCount = 0;
	let foundKeyword = false;

	while(query[0] === "/") {
		query = query.substring(1);
	}

	let providerQueries = [];
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
		if (((await providers[providerName].suggestIf?.(query)) ?? false) && !(ignoreProviders.includes(providerName))){
			providerQueries.push({ providerName: providerName, query: query });
		}
	}

	for (item of providerQueries) {
		const p = providers[item.providerName];
		const q = item.query;
		let urls = await p.getURLs?.(q) ?? [[q, p.getURL?.(q)]];

		for ([title, resultUrl] of urls) {
			let elResult = document.createElement("div");
			elResult.id = "search-result-" + resultCount;
			elResult.classList.add("search-result");

			let elQueryLink = document.createElement("a");
			if (p?.action) {
				elQueryLink.addEventListener('click', p.action);
				elQueryLink.title = q;
				elQueryLink.href = "javascript:;";
			} else {
				elQueryLink.href = resultUrl;
			}
			elQueryLink.innerHTML = '<b>' + p.desc + '</b>: ' + title;
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

			if (p?.color) elResult.style.setProperty('--c', p.color);

			searchResults.insertAdjacentElement("beforeend", elResult);
			resultCount++;
		}
	}

	elid("search-results").insertAdjacentHTML("beforeend", "<p>Retrieved in " + (Date.now() - startTime) + "ms");
}



// ON LOAD -----------------------------------------
function load() {
	//if (performance.navigation.type === PerformanceNavigation.TYPE_NAVIGATE) playSound('/sounds/Woosh2.opus', .4);
	if (elid("toc")) {
		document.addEventListener('scrollend', () => {
			const links = document.querySelectorAll('#toc a[href^="#"]');
			const tocDiv = document.querySelector('#toc div:first-child');
			let current;

			for (link of links) {
				const target = document.querySelector(link.getAttribute('href'));
				if (target?.getBoundingClientRect().bottom <= window.innerHeight / 4) {
					current = link;
				}
			}
			for (link of links) {
				if (link !== current && link.classList.contains('scroll-current')) {
					link.classList.remove('scroll-current');
				}
			}

			if (!(current?.classList.contains('scroll-current'))) {
				current?.classList.add('scroll-current');
			}

			if (tocDiv.scrollWidth > tocDiv.clientWidth && current) {
				tocDiv.scrollTo({ top: 0, left: current.offsetLeft - 16, behavior: 'smooth' });
			}
		});
	}

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
		playSound('/sounds/KDE_Click_2.ogg', 1);
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

	for (a of document.querySelectorAll("section.content a:has(img)")) {
		if (imageExtensions.includes(a.href.split(".").pop())) {
			a.addEventListener("click", (e) => e.preventDefault());
		}
	}
	for (img of document.querySelectorAll("section.content :not(a):not(.no-zoom) img:not(.no-zoom)")) {
		img.addEventListener("click", (e) => {
			if (e.target?.parentElement?.href) {
				if (imageExtensions.includes(e.target.parentElement.href.split(".").pop())) {
					if (e.target.src != e.target.parentElement.href) {
						const iStyle = getComputedStyle(e.target);
						e.target.style.width = iStyle.width;
						e.target.style.height = iStyle.height;
						e.target.addEventListener('load', () => e.target.click(), { once: true });
						e.target.src = e.target.parentElement.href;
						return;
					}
				} else {
					return;
				}
			}
			playSound('/sounds/stone2.ogg');
			if (e.target.classList.contains("zoomed")) {
				e.target.classList.remove("zoomed");
				e.target.style.removeProperty('transform');
			} else {
				e.target.classList.add("zoomed");
				if (!(e.target.classList.contains('pixelated')) && e.target.naturalWidth < (root.clientWidth / 2) && e.target.naturalHeight < (root.clientHeight / 2)) e.target.classList.add('pixelated');
				let b = e.target.getBoundingClientRect();
				let targetTranslateX = (root.clientWidth / 2) - (b.x + (b.width / 2));
				let targetTranslateY = (root.clientHeight / 2) - (b.y + (b.height / 2));
				let targetScale = Math.min(root.clientWidth / e.target.clientWidth, root.clientHeight / e.target.clientHeight, 20);
				e.target.style.transform = 'translate(' + targetTranslateX + 'px, ' + targetTranslateY +'px) scale(' + targetScale + ')';

				const unzoom = () => {
					e.target.classList.remove("zoomed");
					e.target.style.removeProperty('transform');
					document.removeEventListener('scroll', unzoomAfterScroll);
					document.removeEventListener('click', unzoom);
					playSound('/sounds/stone2.ogg');
				}

				const unzoomAfterScroll = () => {
					let b = e.target.getBoundingClientRect();
					if (b.top <= 0 || b.bottom >= window.innerHeight) {
						unzoom();
					}
				}

				document.addEventListener('scroll', unzoomAfterScroll);
				window.setTimeout(() => document.addEventListener('click', unzoom), 100);
			}
		});
	}

	document.addEventListener("keyup", (e) => {
		if (e.key === 'Escape') {
			if (elid("search-box") && elid("search-box")?.value.replaceAll(' ','') === "") {
				closeSearch();
			} else if (elid("search-box")) {
				openSearch("");
			}
		}
	});
	document.addEventListener("keydown", (e) => {
		if (e.key === '/' && document.activeElement.tagName != "INPUT") {
			e.preventDefault();
			openSearch(window.getSelection().toString().replaceAll('\n',''));
		}
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
		openSearch((new URL(window.location.href)).searchParams.get("q"));
		if (window.performance.getEntriesByType("navigation")[0].type === "navigate") { // only click instant links when not back/forward/reload
			document.querySelectorAll("a.search-link.instant")?.[0]?.click?.();
		}
	}
}

load();
