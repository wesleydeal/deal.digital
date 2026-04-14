const root = document.documentElement;
const byId = (id) => document.getElementById(id);
const currentURL = () => new URL(window.location.href);
const plusQuery = (query) => encodeURIComponent(query).replaceAll("%20", "+");

const SEARCH_RESULT_LIMIT = 11;
const SEARCH_SHELL_HTML = `
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
`;
const SEARCH_HELP_HTML = `
	<menu id="search-results"></menu>
	<div id="search-help">
		<p>Press <kbd>/</kbd> to open and <kbd>Esc</kbd> to clear or close.
		<h2>!keywords</h2>
		<ul id="search-keyword-list"></ul>
		<h2>Examples</h2>
		<ul>
			<li><a href="#" onclick="document.documentElement.style.setProperty('--color-primary', 'aquamarine')">color aquamarine</a>
			<li><a href="https://www.ebay.com/sch/i.html?_nkw=+ibm+model+m+(bolt%2Cscrew)+(mod%2Cmodded)">eb ibm model m (bolt,screw) (mod,modded)</a>
			<li><a href="https://youtube.com/results?search_query=+moments+with+heavy+french+toast">yt moments with heavy french toast</a>
			<li><a href="https://chatgpt.com/?q=where+can+I+get+a+good+asada+burrito+nearby+">where can I get a good asada burrito nearby !gpt</a>
			<li><a href="https://annas-archive.org/search?q=mike+ma">mike ma !an</a>
		</ul>
	</div>
`;

let api = null;

export function createSearch({ navigateTo, isMuExcludedURL }) {
	if (api) {
		return api;
	}

	const state = {
		fuse: null,
		fusePromise: null,
		searchWindowPos: null,
		searchDrag: [],
		searchShortcutRegistry: Object.create(null),
	};

	async function initFuse() {
		if (state.fuse) {
			return state.fuse;
		}

		if (!state.fusePromise) {
			state.fusePromise = Promise.all([
				import("/search_index.en.json", { with: { type: "json" } }),
				import("/fuse.min.mjs"),
			]).then(([indexModule, fuseModule]) => new fuseModule.default(indexModule.default, {
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
				],
			}));
		}

		state.fuse = await state.fusePromise;
		return state.fuse;
	}

	async function searchLocal(query) {
		const fuse = await initFuse();
		return fuse.search(query);
	}

	function setURLQuery(key, value) {
		const url = currentURL();
		if (value === null) {
			url.searchParams.delete(key);
		} else {
			url.searchParams.set(key, value);
		}
		window.history.replaceState(window.history.state, "", url);
	}

	function closeSearch() {
		const searchContainer = byId("search-container");
		if (!searchContainer) {
			return;
		}

		setURLQuery("q", null);
		searchContainer.classList.add("hidden");
		window.setTimeout(() => byId("search-container")?.remove(), 300);
	}

	function navigateFromSearch(event) {
		event?.preventDefault?.();
		const searchBox = byId("search-box");
		if (!searchBox) {
			return;
		}

		const path = searchBox.value.replaceAll(" !load", "").trim().split(/\s+/).at(-1);
		if (path) {
			navigateTo(path);
		}
	}

	const providers = {
		Local: {
			keywords: ["d", "dd", "deal", "deal.digital"],
			desc: "deal.digital",
			getURLs: async (query) => {
				const results = await searchLocal(query);
				const bestResult = results[0];
				const currentURLWithoutQuery = window.location.href.split("?")[0].split("#")[0];

				if (bestResult?.score < 0.05 && query.length > 3 && bestResult.item.url !== currentURLWithoutQuery) {
					navigateTo(bestResult.item.url);
				}

				return results.map((result) => [result.item.title, result.item.url]);
			},
		},
		Brave: {
			keywords: ["b", "br", "brave"],
			desc: "Brave Search",
			icon: "/icons/search/brave.png",
			getURL: (query) => `https://search.brave.com/search?q=${plusQuery(query)}`,
			suggestIf: () => true,
			color: "#f50",
		},
		ChatGPT: {
			keywords: ["gpt", "chatgpt"],
			desc: "ChatGPT Search",
			icon: "/icons/search/chatgpt.png",
			getURL: (query) => `https://chatgpt.com/?q=${plusQuery(query)}`,
			suggestIf: () => true,
			color: "#74AA9C",
		},
		Google: {
			keywords: ["g", "google"],
			desc: "Google Search",
			icon: "/icons/search/google.png",
			getURL: (query) => `https://google.com/search?q=${plusQuery(query)}`,
			suggestIf: () => true,
			color: "#1368F4",
		},
		eBay: {
			keywords: ["e", "eb", "ebay"],
			desc: "eBay",
			icon: "/icons/search/ebay.png",
			getURL: (query) => `https://www.ebay.com/sch/i.html?_nkw=${plusQuery(query)}`,
			suggestIf: () => true,
			color: "#3665f3",
		},
		YouTube: {
			keywords: ["y", "yt", "youtube"],
			desc: "YouTube",
			icon: "/icons/search/youtube.png",
			getURL: (query) => `https://youtube.com/results?search_query=${plusQuery(query)}`,
			suggestIf: () => true,
			color: "#f00",
		},
		Amazon: {
			keywords: ["am", "amazon", "amzn"],
			desc: "Amazon",
			icon: "/icons/search/amazon.png",
			getURL: (query) => `https://www.amazon.com/s?k=${plusQuery(query)}`,
			suggestIf: () => true,
			color: "#f90",
		},
		MDN: {
			keywords: ["mdn"],
			desc: "Mozilla Dev",
			icon: "/icons/search/mdn.png",
			getURL: (query) => `https://developer.mozilla.org/en-US/search?q=${plusQuery(query)}`,
			color: "#8cb4ff",
		},
		AnnasArchive: {
			keywords: ["an", "anna", "annas", "annasarchive", "book"],
			desc: "Anna's Archive",
			icon: "/icons/search/annas-archive.png",
			getURL: (query) => `https://annas-archive.org/search?q=${plusQuery(query)}`,
			color: "#0195ff",
		},
		Zola: {
			keywords: ["zola"],
			desc: "Zola Docs",
			icon: "/icons/search/zola.png",
			getURL: (query) => `https://search.brave.com/search?q=site%3Agetzola.org+${plusQuery(query)}`,
			color: "#191919",
		},
		Tera: {
			keywords: ["tera"],
			desc: "Tera Docs",
			getURL: (query) => `https://search.brave.com/search?q=site%3Ahttps%3A%2F%2Fkeats.github.io%2Ftera%2Fdocs%2F+${plusQuery(query)}`,
			color: "#de6262",
		},
		mappletv: {
			keywords: ["tv", "mapple", "mapple.tv"],
			desc: "Mapple.TV",
			icon: "/icons/search/mapple.png",
			getURL: (query) => `https://mapple.tv/search?q=${plusQuery(query)}`,
			color: "#fff",
		},
		WolframAlpha: {
			keywords: ["wa", "wolfram", "wolframalpha"],
			desc: "Wolphram|Alpha",
			getURL: (query) => `http://www.wolframalpha.com/input/?i=${encodeURIComponent(query)}`,
			color: "#ee1f22",
			suggestIf: (query) => ["+", "-", "*", "/", "convert", "per"].some((marker) => query.includes(marker)),
		},
		SetColor: {
			keywords: ["color"],
			desc: "Set Site Color",
			suggestIf: (query) => {
				const style = new Option().style;
				style.color = query;
				return style.color !== "" || query.includes("rand");
			},
			action: (event) => {
				let color = event.currentTarget.title;
				if (color.includes("rand")) {
					color = `#${Array.from(crypto.getRandomValues(new Uint8Array(3))).map((byte) => byte.toString(16).padStart(2, "0")).join("")}`;
				}
				root.style.setProperty("--color-primary", color);
				closeSearch();
			},
		},
		toggleDark: {
			keywords: ["dark", "light", "toggle"],
			desc: "Set Dark Mode",
			suggestIf: (query) => ["dark", "light", "toggle"].includes(query.toLowerCase().replaceAll(" ", "")),
			action: () => {
				const query = byId("search-box")?.value || "";
				if (query.includes("toggle")) {
					const currentDarkMode = getComputedStyle(root).getPropertyValue("--dark-mode") === "true";
					root.style.setProperty("--dark-mode", !currentDarkMode);
				} else if (query.includes("unset") || query.includes("none")) {
					root.style.removeProperty("--dark-mode");
				} else if (query.includes("light")) {
					root.style.setProperty("--dark-mode", false);
				} else {
					root.style.setProperty("--dark-mode", true);
				}
				closeSearch();
			},
		},
		style: {
			keywords: ["style", "stylesheet"],
			desc: "Set Stylesheet",
			hide: true,
			action: (event) => {
				const stylesheet = document.querySelector("link[rel='stylesheet'][as='style']");
				if (!stylesheet) {
					return;
				}
				stylesheet.href = `/${event.currentTarget.title.replaceAll(" ", "")}.css`;
				closeSearch();
			},
		},
		load: {
			keywords: ["load"],
			desc: "Load Content From Internal Page",
			hide: true,
			action: navigateFromSearch,
		},
	};

	const keywordMap = new Map(
		Object.entries(providers).flatMap(([providerName, provider]) => (
			(provider.keywords || []).map((keyword) => [keyword, providerName])
		)),
	);

	function populateSearchKeywordList(keywordList) {
		for (const provider of Object.values(providers)) {
			if (!provider.keywords || provider.hide || !keywordList) {
				continue;
			}

			const keywordEntry = document.createElement("li");
			const keywordSample = document.createElement("samp");

			keywordEntry.textContent = provider.desc;
			keywordSample.textContent = provider.keywords[0];
			keywordEntry.append(keywordSample);
			keywordEntry.addEventListener("click", () => openSearch(`${keywordSample.textContent} `));

			if (provider.color) {
				keywordEntry.style.setProperty("--c", provider.color);
			}

			keywordList.append(keywordEntry);
		}
	}

	async function collectSearchProviderQueries(query) {
		const providerQueries = new Map();
		let foundKeyword = false;

		for (const word of query.toLowerCase().split(" ").reverse()) {
			const bangIndex = word.indexOf("!");
			if (bangIndex < 0) {
				continue;
			}

			const providerName = keywordMap.get(word.slice(bangIndex + 1));
			if (!providerName) {
				continue;
			}

			providerQueries.set(providerName, query.replace(/!.*?( |$)/g, ""));
			foundKeyword = true;
		}

		const firstKeyword = keywordMap.get(query.toLowerCase().split(" ")[0]);
		if (firstKeyword) {
			providerQueries.set(firstKeyword, query.includes(" ") ? query.slice(query.indexOf(" ")) : "");
		}

		const ignoredProviders = new Set(providerQueries.keys());
		for (const [providerName, provider] of Object.entries(providers)) {
			const suggested = await provider.suggestIf?.(query);
			if (suggested && !ignoredProviders.has(providerName)) {
				providerQueries.set(providerName, query);
			}
		}

		return { providerQueries, foundKeyword };
	}

	function buildSearchLink(provider, query, title, resultURL, foundKeyword) {
		const result = document.createElement("div");
		result.className = "search-result";

		const link = document.createElement("a");
		link.className = "search-link";
		if (foundKeyword) {
			link.classList.add("instant");
		}

		if (provider.action) {
			link.href = "#";
			link.title = query;
			link.addEventListener("click", (event) => {
				event.preventDefault();
				provider.action(event);
			});
		} else {
			link.href = resultURL;
			try {
				if (isMuExcludedURL(new URL(resultURL, window.location.href))) {
					link.setAttribute("mu-disabled", "");
				}
			} catch {}
		}

		const providerLabel = document.createElement("b");
		providerLabel.textContent = provider.desc;
		link.append(providerLabel, `: ${title}`);
		result.append(link);

		result.addEventListener("click", () => link.click());
		if (provider.color) {
			result.style.setProperty("--c", provider.color);
		}

		return { result, link };
	}

	function appendSearchShortcut(result, link, resultCount) {
		if (resultCount >= SEARCH_RESULT_LIMIT) {
			return;
		}

		const shortcut = document.createElement("div");
		shortcut.className = "search-shortcut";

		if (resultCount === 0) {
			shortcut.innerHTML = "<kbd>ENTER</kbd>";
			state.searchShortcutRegistry.Enter = link;
		} else {
			shortcut.innerHTML = `<kbd>Alt</kbd> <kbd>${resultCount}</kbd>`;
			state.searchShortcutRegistry[`Alt${resultCount}`] = link;
		}

		result.append(shortcut);
	}

	async function renderSearchResults(searchResults, providerQueries, foundKeyword, startTime) {
		let resultCount = 0;

		for (const [providerName, providerQuery] of providerQueries.entries()) {
			const provider = providers[providerName];
			const urls = await provider.getURLs?.(providerQuery) ?? [[providerQuery, provider.getURL?.(providerQuery)]];

			for (const [title, resultURL] of urls) {
				const { result, link } = buildSearchLink(provider, providerQuery, title, resultURL, foundKeyword);
				appendSearchShortcut(result, link, resultCount);
				searchResults.append(result);
				resultCount++;
			}
		}

		const timing = document.createElement("p");
		timing.textContent = `Retrieved in ${Date.now() - startTime}ms`;
		searchResults.append(timing);
	}

	async function updateSearch() {
		state.searchShortcutRegistry = Object.create(null);

		const startTime = Date.now();
		const searchBox = byId("search-box");
		const searchResults = byId("search-results");
		if (!searchBox || !searchResults) {
			return;
		}

		let query = searchBox.value;
		while (query.startsWith("/")) {
			query = query.slice(1);
		}

		searchResults.replaceChildren();
		setURLQuery("q", query);

		if (query.replaceAll(" ", "") === "") {
			return;
		}

		const { providerQueries, foundKeyword } = await collectSearchProviderQueries(query);
		await renderSearchResults(searchResults, providerQueries, foundKeyword, startTime);
	}

	function minimizeSearch() {
		byId("search-container")?.classList.add("min");
	}

	function toggleMaximizeSearch() {
		const searchContainer = byId("search-container");
		if (!searchContainer) {
			return;
		}

		const style = getComputedStyle(searchContainer);
		if (searchContainer.classList.contains("max")) {
			searchContainer.classList.remove("max");
			for (const [propName, value] of Object.entries(state.searchWindowPos || {})) {
				searchContainer.style.setProperty(propName, value);
			}
			return;
		}

		state.searchWindowPos = {
			top: style.top,
			left: style.left,
			width: style.width,
			height: style.height,
		};
		searchContainer.style.removeProperty("top");
		searchContainer.style.removeProperty("left");
		searchContainer.style.removeProperty("width");
		searchContainer.style.removeProperty("height");
		searchContainer.classList.add("max");
	}

	function onDragRelease() {
		state.searchDrag = [];
		document.removeEventListener("pointermove", dragSearchUpdate, { passive: false });
		byId("search-container")?.classList.remove("drag");
		root.style.removeProperty("touch-action");
		document.removeEventListener("pointerup", onDragRelease, { passive: false });
	}

	function dragSearchStart(event) {
		const searchContainer = byId("search-container");
		const titlebarButtons = document.querySelector("#search-titlebar .window-buttons");
		if (!searchContainer || !titlebarButtons || titlebarButtons.contains(event.target) || event.buttons !== 1) {
			return;
		}

		if (searchContainer.classList.contains("max")) {
			const originalX = event.clientX;
			const originalY = event.clientY;
			const moveListener = (moveEvent) => {
				if (moveEvent.buttons !== 1) {
					searchContainer.removeEventListener("pointermove", moveListener);
					return;
				}

				const distance = Math.hypot(originalX - moveEvent.clientX, originalY - moveEvent.clientY);
				if (distance <= 4 || !state.searchWindowPos?.width) {
					return;
				}

				searchContainer.style.setProperty("width", state.searchWindowPos.width);
				searchContainer.style.setProperty("height", "auto");
				searchContainer.style.setProperty("top", "0");
				searchContainer.style.setProperty("left", `${(moveEvent.clientX / window.innerWidth) * (window.innerWidth - parseFloat(state.searchWindowPos.width))}px`);
				searchContainer.classList.remove("max");
				searchContainer.removeEventListener("pointermove", moveListener);
				dragSearchStart(moveEvent);
			};

			searchContainer.addEventListener("pointermove", moveListener);
			return;
		}

		if (state.searchDrag.length > 0) {
			return;
		}

		state.searchDrag = [event.clientX - searchContainer.offsetLeft, event.clientY - searchContainer.offsetTop];
		searchContainer.classList.add("drag");
		root.style.setProperty("touch-action", "none");
		document.addEventListener("pointermove", dragSearchUpdate, { passive: false });
		document.addEventListener("pointerup", onDragRelease, { passive: false });
	}

	function dragSearchUpdate(event) {
		const searchContainer = byId("search-container");
		if (!searchContainer || state.searchDrag.length < 2) {
			return;
		}

		searchContainer.style.setProperty("right", "unset");
		searchContainer.style.setProperty("left", `${event.clientX - state.searchDrag[0]}px`);
		searchContainer.style.setProperty("top", `${event.clientY - state.searchDrag[1]}px`);
	}

	function openSearch(query = "") {
		if (query instanceof Event || query === null) {
			query = "";
		}

		let searchContainer = byId("search-container");
		if (!searchContainer) {
			document.body.insertAdjacentHTML("afterbegin", SEARCH_SHELL_HTML);
			searchContainer = byId("search-container");
			byId("search-inner")?.insertAdjacentHTML("beforeend", SEARCH_HELP_HTML);
			populateSearchKeywordList(byId("search-keyword-list"));
			byId("search-close")?.addEventListener("click", closeSearch);
			byId("search-min")?.addEventListener("click", minimizeSearch);
			byId("search-max")?.addEventListener("click", toggleMaximizeSearch);
			byId("search-titlebar")?.addEventListener("dblclick", toggleMaximizeSearch);
			byId("search-titlebar")?.addEventListener("pointerdown", dragSearchStart);
			byId("search-box")?.addEventListener("input", updateSearch);
		} else {
			const searchBox = byId("search-box");
			if (searchBox) {
				searchBox.placeholder = searchBox.value;
			}
		}

		const searchBox = byId("search-box");
		if (!searchBox) {
			return;
		}

		searchBox.value = String(query);
		void updateSearch();
		searchBox.focus();
		setURLQuery("q", searchBox.value);
	}

	function toggleSearch(query = "") {
		const searchContainer = byId("search-container");
		if (searchContainer) {
			if (searchContainer.classList.contains("min")) {
				searchContainer.classList.remove("min");
				byId("search-box")?.focus();
			} else {
				closeSearch();
			}
			return;
		}

		openSearch(query);
	}

	function handleKeyUp(event) {
		if (event.key !== "Escape") {
			return;
		}

		const searchBox = byId("search-box");
		if (!searchBox) {
			return;
		}

		if (searchBox.value.replaceAll(" ", "") === "") {
			closeSearch();
			return;
		}

		openSearch("");
	}

	function handleKeyDown(event) {
		const searchBox = byId("search-box");
		if (document.activeElement !== searchBox) {
			return;
		}

		const key = event.altKey ? `Alt${event.key}` : event.key;
		const link = state.searchShortcutRegistry[key];
		if (!link) {
			return;
		}

		event.preventDefault();
		if (event.ctrlKey) {
			const originalTarget = link.target;
			link.target = "_blank";
			link.click();
			link.target = originalTarget;
			return;
		}

		link.click();
	}

	function syncFromLocation({ allowInstant = false } = {}) {
		const query = currentURL().searchParams.get("q");
		if (!query) {
			return;
		}

		openSearch(query);
		if (allowInstant) {
			document.querySelector("a.search-link.instant")?.click?.();
		}
	}

	api = {
		handleKeyDown,
		handleKeyUp,
		open: openSearch,
		syncFromLocation,
		toggle: toggleSearch,
	};
	return api;
}
