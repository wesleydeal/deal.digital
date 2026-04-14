// STATE -------------------------------------------
const root = document.documentElement;
const imageExtensions = new Set(["svg", "png", "jpg", "jfif", "gif", "webp", "avif"]);

const state = {
	soundPlayers: {},
	searchUI: null,
	searchPromise: null,
	globalEventsBound: false,
	pendingMu: null,
};

// HELPERS -----------------------------------------
const byId = (id) => document.getElementById(id);
const currentContent = () => document.querySelector(".content");
const currentURL = () => new URL(window.location.href);
const parseHTML = (html) => new DOMParser().parseFromString(html, "text/html");

// NAVIGATION --------------------------------------
var mu = window.mu || new function() {
	this._lastUrl = null;
	this._prevUrl = null;
	this._abortCtrl = null;
	this._bar = null;
	this._prefetch = new Map();
	this._hoverTimer = null;
	this._prefetchTtl = 3000;
	this._confirmQuit = false;
	this._leaveText = "Are you sure you want to leave this page?";
	this._jsIncludes = {};
	this._morph = null;
	this._initialized = false;

	this._newCfg = function() {
		return ({
			history: true,
			mode: "replace",
			target: "body",
			source: "body",
			scroll: null,
		});
	};
	this.init = function() {
		if (!mu._morph && typeof window.Idiomorph !== "undefined" && typeof window.Idiomorph.morph === "function") {
			mu._morph = function(target, html, opts) {
				window.Idiomorph.morph(target, html, opts);
			};
		}
		var existingScripts = document.querySelectorAll("script[src]");
		for (var i = 0; i < existingScripts.length; i++)
			mu._jsIncludes[existingScripts[i].getAttribute("src")] = true;
		if (!mu._initialized) {
			document.addEventListener("click", mu._onClick);
			document.addEventListener("submit", mu._onSubmit);
			document.addEventListener("mouseover", mu._onMouseOver);
			document.addEventListener("mouseout", mu._onMouseOut);
			document.addEventListener("input", mu._onInput);
			window.addEventListener("popstate", mu._onPopState);
			window.addEventListener("beforeunload", mu._onBeforeUnload);
			mu._initialized = true;
		}
		window.history.replaceState({ mu: true, url: location.pathname + location.search }, "", location.pathname + location.search);
		mu._initTriggers(document.body);
		mu._emit("mu:init", { url: location.pathname + location.search });
	};

	this.load = function(url) {
		mu._loadExec(mu._resolveUrl(url) || url, mu._newCfg());
	};
	this._attr = function(el, name) {
		return (el.getAttribute("mu-" + name));
	};
	this._attrBool = function(el, name, fallback) {
		var v = mu._attr(el, name);
		if (v === null)
			return (fallback);
		if (v === "" || v === "true")
			return (true);
		if (v === "false")
			return (false);
		return (fallback);
	};

	this._resolveTarget = function(selector, sourceEl) {
		if (!selector || selector.indexOf("&") === -1 || !sourceEl)
			return (selector);
		if (!sourceEl.id)
			sourceEl.id = "mu-" + Math.random().toString(36).slice(2);
		return (selector.replace(/&/g, "#" + sourceEl.id));
	};

	this._resolveUrl = function(url) {
		if (!url || url.charAt(0) === "#")
			return (null);
		if (url.charAt(0) === "/" && url.charAt(1) !== "/")
			return (url);
		try {
			var parsed = new URL(url, document.baseURI);
			if (parsed.origin === window.location.origin)
				return (parsed.pathname + parsed.search + parsed.hash);
		} catch(e) {}
		return (null);
	};
	this._isHtmlUrl = function(url) {
		var path = url.split("#")[0].split("?")[0];
		var leaf = path.substring(path.lastIndexOf("/") + 1);
		var dot = leaf.lastIndexOf(".");
		return (dot === -1 || /html?$/i.test(leaf.substring(dot + 1)));
	};
	this._resolveMediaUrl = function(value, baseURL) {
		var t = value.trim();
		if (t && !/^[#/]/.test(t) && !/^[a-z][a-z\d+.-]*:/i.test(t)) {
			try {
				var resolved = new URL(t, baseURL);
				return (resolved.origin === window.location.origin ? resolved.pathname + resolved.search + resolved.hash : resolved.toString());
			} catch {}
		}
		return (value);
	};
	this._resolveMediaSrcset = function(value, baseURL) {
		return (value.split(",").map(function(candidate) {
			var trimmed = candidate.trim();
			if (!trimmed)
				return (trimmed);
			var parts = trimmed.split(/\s+/);
			var resolved = mu._resolveMediaUrl(parts[0], baseURL);
			return (parts.length > 1 ? resolved + " " + parts.slice(1).join(" ") : resolved);
		}).join(", "));
	};
	this._normalizeMedia = function(scope, baseURL) {
		if (!scope || !baseURL)
			return;
		var els = scope.querySelectorAll("img[src], img[srcset], source[src], source[srcset], video[poster]");
		for (var i = 0; i < els.length; i++) {
			if (els[i].hasAttribute("src"))
				els[i].setAttribute("src", mu._resolveMediaUrl(els[i].getAttribute("src"), baseURL));
			if (els[i].hasAttribute("srcset"))
				els[i].setAttribute("srcset", mu._resolveMediaSrcset(els[i].getAttribute("srcset"), baseURL));
			if (els[i].hasAttribute("poster"))
				els[i].setAttribute("poster", mu._resolveMediaUrl(els[i].getAttribute("poster"), baseURL));
		}
	};

	this._shouldProcess = function(el) {
		if (mu._attr(el, "disabled") === "true" || mu._attr(el, "disabled") === "")
			return (false);
		if (el.hasAttribute("target") || el.hasAttribute("download"))
			return (false);
		if ((el.tagName === "A" && el.hasAttribute("onclick")) || (el.tagName === "FORM" && el.hasAttribute("onsubmit")))
			return (false);
		var url = mu._resolveUrl(mu._attr(el, "url") || el.getAttribute("href") || el.getAttribute("action") || "");
		return (url !== null && (el.tagName !== "A" || mu._isHtmlUrl(url)));
	};

	this._elemCfg = function(el) {
		var cfg = mu._newCfg();
		var v;
		if ((v = mu._attr(el, "mode")) !== null)
			cfg.mode = v;
		if ((v = mu._attr(el, "target")) !== null)
			cfg.target = v;
		if ((v = mu._attr(el, "source")) !== null)
			cfg.source = v;
		if ((v = mu._attr(el, "url")) !== null)
			cfg._url = v;
		cfg.history = mu._attrBool(el, "history", cfg.history);
		cfg.scroll = mu._attrBool(el, "scroll", cfg.scroll);
		if ((v = mu._attr(el, "method")) !== null)
			cfg.method = v.toLowerCase();
		cfg.confirm = mu._attr(el, "confirm");
		cfg.patchHistory = mu._attrBool(el, "patch-history", false);
		if (cfg.mode !== "replace" && cfg.mode !== "update" && cfg.mode !== "patch") {
			if (mu._attr(el, "history") === null)
				cfg.history = false;
			if (mu._attr(el, "scroll") === null && cfg.scroll === null)
				cfg.scroll = false;
		}
		return (cfg);
	};

	this._copyScript = function(node) {
		var s = document.createElement("script");
		for (var j = 0; j < node.attributes.length; j++)
			s.setAttribute(node.attributes[j].name, node.attributes[j].value);
		s.textContent = node.textContent;
		return (s);
	};

	this._onClick = function(e) {
		var el = e.target.closest("[mu-url], a");
		if (!el || !mu._shouldProcess(el))
			return;
		if (mu._getTrigger(el) !== "click")
			return;
		if (e.ctrlKey || e.metaKey || e.shiftKey || e.altKey)
			return;
		e.preventDefault();
		mu._triggerAction(el, true);
	};
	this._onSubmit = function(e) {
		var form = e.target.closest("form");
		if (!form || !mu._shouldProcess(form))
			return;
		if (mu._getTrigger(form) !== "submit")
			return;
		if (!form.reportValidity())
			return;
		var validator = mu._attr(form, "validate");
		if (validator && typeof window[validator] === "function" && !window[validator](form))
			return;
		var cfg = mu._elemCfg(form);
		var method = cfg.method || (form.getAttribute("method") || "get").toLowerCase();
		cfg.method = method;
		var url = mu._resolveUrl(cfg._url || form.getAttribute("action"));
		if (!url)
			return;
		cfg._el = form;
		e.preventDefault();
		var formData = new FormData(form);
		var submitter = e.submitter;
		if (submitter && submitter.form === form && submitter.name)
			formData.append(submitter.name, submitter.value || "");
		if (method === "get") {
			var qs = new URLSearchParams(formData).toString();
			url = url + "?" + qs;
		} else {
			cfg.postData = form.enctype === "multipart/form-data" ? formData : new URLSearchParams(formData);
			if (mu._attr(form, "history") === null)
				cfg.history = false;
			if (cfg.scroll === null && cfg.mode !== "patch")
				cfg.scroll = true;
		}
		mu._confirmQuit = false;
		mu._loadExec(url, cfg);
	};
	this._onInput = function(e) {
		if (e.target.closest("form[mu-confirm-quit]"))
			mu._confirmQuit = true;
	};
	this._onMouseOver = function(e) {
		var el = e.target.closest("[mu-url], a");
		if (!el || !mu._shouldProcess(el))
			return;
		if (mu._getTrigger(el) !== "click")
			return;
		var method = mu._attr(el, "method");
		if (method && method.toLowerCase() !== "get")
			return;
		if (mu._attrBool(el, "prefetch", true) === false)
			return;
		var url = mu._resolveUrl(mu._attr(el, "url") || el.getAttribute("href"));
		if (!url)
			return;
		var existing = mu._prefetch.get(url);
		if (existing && (Date.now() - existing.ts) < mu._prefetchTtl)
			return;
		if (url === location.pathname + location.search)
			return;
		clearTimeout(mu._hoverTimer);
		mu._hoverTimer = setTimeout(function() {
			mu._hoverTimer = null;
			var cached = mu._prefetch.get(url);
			if (cached && (Date.now() - cached.ts) < mu._prefetchTtl)
				return;
			var promise = fetch(url, {
				headers: { "X-Requested-With": "XMLHttpRequest", "X-Mu-Prefetch": "1" }
			})
			.then(function(r) { return (r.ok ? r.text() : null); })
			.catch(function() { return (null); });
			mu._prefetch.set(url, { promise: promise, ts: Date.now() });
		}, 50);
	};
	this._onMouseOut = function() {
		if (mu._hoverTimer) {
			clearTimeout(mu._hoverTimer);
			mu._hoverTimer = null;
		}
	};
	this._onPopState = function(e) {
		var state = e.state;
		if (!state || !state.mu)
			return;
		var cfg = mu._newCfg();
		cfg.history = false;
		cfg.scroll = false;
		cfg._popstate = true;
		cfg._scrollPos = state.scrollX !== undefined ? { x: state.scrollX, y: state.scrollY } : null;
		mu._loadExec(state.url, cfg);
	};
	this._onBeforeUnload = function(e) {
		if (mu._confirmQuit) {
			e.preventDefault();
			e.returnValue = mu._leaveText;
		}
	};

	this._getTrigger = function(el) {
		var t = mu._attr(el, "trigger");
		if (t)
			return (t);
		t = el.tagName;
		if (t === "FORM")
			return ("submit");
		if (t === "INPUT" || t === "TEXTAREA" || t === "SELECT")
			return ("change");
		return ("click");
	};
	this._debounce = function(fn, delay) {
		var timer = null;
		return function() {
			clearTimeout(timer);
			timer = setTimeout(fn, delay);
		};
	};
	this._triggerAction = function(el, isClick) {
		if (!mu._shouldProcess(el))
			return;
		var cfg = mu._elemCfg(el);
		var url = mu._resolveUrl(cfg._url || el.getAttribute("href") || el.getAttribute("action"));
		if (!url)
			return;
		cfg._el = el;
		if (!cfg.method)
			cfg.method = "get";
		if (isClick) {
			if (cfg.confirm && !window.confirm(cfg.confirm))
				return;
			if (mu._confirmQuit) {
				if (!window.confirm(mu._leaveText))
					return;
				mu._confirmQuit = false;
			}
			if (cfg.method !== "get" && mu._attr(el, "history") === null)
				cfg.history = false;
		} else {
			cfg._trigger = true;
			if (mu._attr(el, "history") === null)
				cfg.history = false;
			if (mu._attr(el, "scroll") === null && cfg.scroll === null)
				cfg.scroll = false;
			var t = el.tagName;
			if (t === "INPUT" || t === "TEXTAREA" || t === "SELECT") {
				var form = el.closest("form");
				if (form) {
					if (cfg.method === "get") {
						var formData = new FormData(form);
						var qs = new URLSearchParams(formData).toString();
						url = url.split("?")[0] + "?" + qs;
					} else {
						cfg.postData = form.enctype === "multipart/form-data" ? new FormData(form) : new URLSearchParams(new FormData(form));
					}
				} else if (el.name && cfg.method === "get") {
					url = url.split("?")[0] + "?" + encodeURIComponent(el.name) + "=" + encodeURIComponent(el.value);
				}
			}
		}
		if (cfg.method === "sse") {
			mu._openSSE(url, el, cfg);
			return;
		}
		mu._loadExec(url, cfg);
	};
	this._initTriggers = function(container) {
		var els = container.querySelectorAll("[mu-url], [mu-trigger]");
		if (container.matches && container.matches("[mu-url], [mu-trigger]")) {
			var tmp = [container];
			for (var j = 0; j < els.length; j++)
				tmp.push(els[j]);
			els = tmp;
		}
		for (var i = 0; i < els.length; i++) {
			var el = els[i];
			if (el._mu_bound)
				continue;
			var trigger = mu._getTrigger(el);
			if (trigger === "click" || trigger === "submit")
				continue;
			var url = mu._attr(el, "url") || el.getAttribute("href") || el.getAttribute("action");
			if (!url)
				continue;
			el._mu_bound = true;
			var debounceMs = parseInt(mu._attr(el, "debounce"), 10) || 0;
			var repeatMs = parseInt(mu._attr(el, "repeat"), 10) || 0;
			var handler = (function(targetEl) {
				return function() { mu._triggerAction(targetEl); };
			})(el);
			if (debounceMs > 0)
				handler = mu._debounce(handler, debounceMs);
			if (repeatMs > 0) {
				handler = (function(targetEl, fn, ms) {
					var started = false;
					return function() {
						fn();
						if (!started) {
							started = true;
							targetEl._mu_interval = setInterval(fn, ms);
						}
					};
				})(el, handler, repeatMs);
			}
			if (trigger === "change") {
				el.addEventListener("input", handler);
			} else if (trigger === "blur") {
				var dedupHandler = (function(fn) {
					var lastFired = 0;
					return function() {
						var now = Date.now();
						if (now - lastFired < 50)
							return;
						lastFired = now;
						fn();
					};
				})(handler);
				el.addEventListener("change", dedupHandler);
				el.addEventListener("blur", dedupHandler);
			} else if (trigger === "focus") {
				el.addEventListener("focus", handler);
			} else if (trigger === "load") {
				handler();
			}
		}
	};
	this._cleanupTriggers = function(container) {
		var els = container.querySelectorAll ? container.querySelectorAll("*") : [];
		for (var i = 0; i < els.length; i++) {
			if (els[i]._mu_interval) {
				clearInterval(els[i]._mu_interval);
				els[i]._mu_interval = null;
			}
			if (els[i]._mu_sse) {
				els[i]._mu_sse.close();
				els[i]._mu_sse = null;
			}
			els[i]._mu_bound = false;
		}
		if (container._mu_interval) {
			clearInterval(container._mu_interval);
			container._mu_interval = null;
		}
		if (container._mu_sse) {
			container._mu_sse.close();
			container._mu_sse = null;
		}
		container._mu_bound = false;
	};
	this._openSSE = function(url, el, cfg) {
		if (el._mu_sse)
			el._mu_sse.close();
		var source = new EventSource(url);
		el._mu_sse = source;
		source.onmessage = function(e) {
			var detail = { url: url, html: e.data, config: cfg };
			if (!mu._emit("mu:before-render", detail))
				return;
			if (cfg.mode === "patch") {
				mu._renderPatch(detail.html, cfg);
			} else {
				mu._renderPage(detail.html, cfg);
			}
			mu._emit("mu:after-render", { url: url, finalUrl: url, mode: cfg.mode });
		};
		source.onerror = function() {
			mu._emit("mu:fetch-error", { url: url, fetchUrl: url, error: new Error("SSE connection error") });
		};
	};

	this._loadExec = async function(url, cfg) {
		if (!cfg._trigger && !cfg._popstate)
			mu._saveScroll();
		if (!mu._emit("mu:before-fetch", { url: url, fetchUrl: url, config: cfg, sourceElement: cfg._el || null }))
			return;
		var abortCtrl;
		if (cfg._trigger) {
			abortCtrl = new AbortController();
		} else {
			if (mu._abortCtrl)
				mu._abortCtrl.abort();
			abortCtrl = mu._abortCtrl = new AbortController();
		}
		if (!cfg._trigger)
			mu._showProgress();
		try {
			var html = null;
			var resp = null;
			var finalUrl = url;
			var method = cfg.method || "get";
			var cached = mu._prefetch.get(url);
			if (method === "get" && cached && cached.promise && (Date.now() - cached.ts) < mu._prefetchTtl) {
				html = await cached.promise;
				if (abortCtrl.signal.aborted)
					return;
			}
			if (!html) {
				var fetchOpts = {
					signal: abortCtrl.signal,
					headers: {
						"X-Requested-With": "XMLHttpRequest",
						"X-Mu-Mode": cfg.mode,
					},
				};
				if (method !== "get") {
					fetchOpts.headers["X-Mu-Method"] = fetchOpts.method = method.toUpperCase();
					if (cfg.postData)
						fetchOpts.body = cfg.postData;
				}
				resp = await fetch(url, fetchOpts);
				if (!resp.ok) {
					mu._emit("mu:fetch-error", { url: url, fetchUrl: url, status: resp.status, response: resp });
					return;
				}
				if (resp.redirected) {
					var redirected = new URL(resp.url);
					finalUrl = redirected.pathname + redirected.search;
				}
				html = await resp.text();
			}
			cfg._docUrl = resp && resp.redirected ? resp.url : new URL(finalUrl, window.location.href).toString();
			if (!cfg._trigger && method === "get")
				mu._prefetch.set(url, { promise: Promise.resolve(html), ts: Date.now() });
			var detail = { url: url, finalUrl: finalUrl, html: html, config: cfg };
			if (!mu._emit("mu:before-render", detail))
				return;
			if (cfg.mode === "patch") {
				cfg._addHistory = cfg.patchHistory;
			} else {
				cfg._addHistory = cfg.history;
				if (resp && resp.redirected)
					cfg._addHistory = true;
			}
			var postRender = function() {
				mu._prevUrl = mu._lastUrl;
				mu._lastUrl = finalUrl;
				if (cfg._addHistory)
					window.history.pushState({ mu: true, url: finalUrl }, "", finalUrl);
				if (cfg.mode !== "patch") {
					if (cfg._scrollPos) {
						window.scrollTo(cfg._scrollPos.x, cfg._scrollPos.y);
					} else if (cfg.scroll !== false) {
						var hashIdx = url.indexOf("#");
						if (hashIdx !== -1) {
							var anchor = document.getElementById(url.substring(hashIdx + 1));
							if (anchor)
								anchor.scrollIntoView({ behavior: "smooth" });
						} else {
							window.scrollTo(0, 0);
						}
					}
				}
				mu._confirmQuit = false;
				mu._emit("mu:after-render", { url: url, finalUrl: finalUrl, mode: cfg.mode });
			};
			var applyDom = cfg.mode === "patch"
				? function() { mu._renderPatch(detail.html, cfg); }
				: function() { mu._renderPage(detail.html, cfg); };
			if (!cfg._trigger && document.startViewTransition) {
				document.startViewTransition(applyDom).updateCallbackDone.then(postRender);
			} else {
				applyDom();
				postRender();
			}
		} catch (err) {
			if (err.name !== "AbortError")
				mu._emit("mu:fetch-error", { url: url, fetchUrl: url, error: err });
		} finally {
			if (!cfg._trigger)
				mu._hideProgress();
		}
	};

	this._renderPage = function(html, cfg) {
		var doc = parseHTML(html);
		mu._normalizeMedia(doc, cfg._docUrl);
		var sourceNode = cfg.source ? doc.querySelector(cfg.source) : null;
		if (!sourceNode)
			sourceNode = doc.body;
		var resolvedTarget = mu._resolveTarget(cfg.target, cfg._el);
		var targetNode = document.querySelector(resolvedTarget);
		if (!targetNode) {
			console.warn("[µJS] Target element '" + resolvedTarget + "' not found.");
			return;
		}
		mu._cleanupTriggers(targetNode);
		mu._applyMode(cfg.mode, targetNode, sourceNode);
		if (cfg._addHistory)
			mu._updateTitle(doc);
		mu._mergeHead(doc);
		var container = document.querySelector(resolvedTarget) || document.body;
		mu._runScripts(container);
		mu._initTriggers(container);
	};

	this._renderPatch = function(html, cfg) {
		var doc = parseHTML(html);
		mu._normalizeMedia(doc, cfg._docUrl);
		var fragments = doc.querySelectorAll("[mu-patch-target]");
		var patchedSelectors = [];
		for (var i = 0; i < fragments.length; i++) {
			var frag = fragments[i];
			var targetSel = mu._resolveTarget(frag.getAttribute("mu-patch-target"), cfg._el);
			var mode = frag.getAttribute("mu-patch-mode") || "replace";
			var targetNode = document.querySelector(targetSel);
			if (!targetNode) {
				console.warn("[µJS] Patch target '" + targetSel + "' not found.");
				continue;
			}
			mu._cleanupTriggers(targetNode);
			mu._applyMode(mode, targetNode, frag);
			if (mode !== "remove") {
				mu._runScripts(frag);
				patchedSelectors.push(targetSel);
			}
		}
		for (var j = 0; j < patchedSelectors.length; j++) {
			var patched = document.querySelector(patchedSelectors[j]);
			if (patched)
				mu._initTriggers(patched);
		}
	};

	this._applyMode = function(mode, targetNode, sourceNode) {
		var useMorph = mu._morph;
		switch (mode) {
			case "update":
				if (useMorph) {
					mu._morph(targetNode, sourceNode.innerHTML, { morphStyle: "innerHTML" });
				} else {
					targetNode.innerHTML = sourceNode.innerHTML;
				}
				break;
			case "prepend":
				targetNode.prepend(sourceNode);
				break;
			case "append":
				targetNode.append(sourceNode);
				break;
			case "before":
				targetNode.before(sourceNode);
				break;
			case "after":
				targetNode.after(sourceNode);
				break;
			case "remove":
				targetNode.remove();
				break;
			case "none":
				break;
			case "replace":
			default:
				if (targetNode.tagName === "BODY" && sourceNode.tagName === "BODY") {
					if (useMorph) {
						mu._morph(targetNode, sourceNode.innerHTML, { morphStyle: "innerHTML" });
					} else {
						targetNode.innerHTML = sourceNode.innerHTML;
					}
				} else if (useMorph) {
					mu._morph(targetNode, sourceNode.outerHTML, { morphStyle: "outerHTML" });
				} else {
					targetNode.replaceWith(sourceNode);
				}
				break;
		}
	};

	this._updateTitle = function(doc) {
		var el = doc.querySelector("title");
		if (el)
			document.title = el.textContent;
	};
	this._mergeHead = function(doc) {
		var selector = "link[rel='stylesheet'], style, script";
		var oldEls = document.head.querySelectorAll(selector);
		var newEls = doc.head.querySelectorAll(selector);
		var oldKeys = new Set();
		for (var i = 0; i < oldEls.length; i++)
			oldKeys.add(mu._elKey(oldEls[i]));
		for (var j = 0; j < newEls.length; j++) {
			if (oldKeys.has(mu._elKey(newEls[j])))
				continue;
			if (newEls[j].tagName.toUpperCase() === "SCRIPT") {
				var s = mu._copyScript(newEls[j]);
				if (s.hasAttribute("src"))
					mu._jsIncludes[s.getAttribute("src")] = true;
				document.head.appendChild(s);
			} else {
				document.head.appendChild(newEls[j].cloneNode(true));
			}
		}
	};
	this._elKey = function(el) {
		var tag = el.tagName.toUpperCase();
		if (tag === "LINK")
			return ("link:" + el.getAttribute("href"));
		if (tag === "STYLE")
			return ("style:" + el.textContent.substring(0, 100));
		if (tag === "SCRIPT")
			return ("script:" + (el.getAttribute("src") || el.textContent.substring(0, 100)));
		return (el.outerHTML);
	};

	this._runScripts = function(container) {
		var scripts = container.querySelectorAll("script");
		for (var i = 0; i < scripts.length; i++) {
			var old = scripts[i];
			if (mu._attr(old, "disabled") === "true" || mu._attr(old, "disabled") === "")
				continue;
			if (old.hasAttribute("src")) {
				var src = old.getAttribute("src");
				if (mu._jsIncludes[src])
					continue;
				mu._jsIncludes[src] = true;
			}
			old.parentNode.replaceChild(mu._copyScript(old), old);
		}
	};

	this._showProgress = function() {
		if (!mu._bar) {
			mu._bar = document.createElement("div");
			mu._bar.id = "mu-progress";
			mu._bar.style.cssText = "position:fixed;top:0;left:0;height:3px;background:#29d;z-index:99999;transition:width .3s ease;width:0";
		}
		document.body.appendChild(mu._bar);
		mu._bar.offsetWidth;
		mu._bar.style.width = "70%";
	};
	this._hideProgress = function() {
		if (!mu._bar)
			return;
		mu._bar.style.width = "100%";
		setTimeout(function() {
			mu._bar.style.transition = "none";
			mu._bar.style.width = "0";
			mu._bar.offsetWidth;
			mu._bar.style.transition = "width .3s ease";
			mu._bar.remove();
		}, 200);
	};

	this._saveScroll = function() {
		var state = window.history.state;
		if (state && state.mu) {
			state.scrollX = window.scrollX;
			state.scrollY = window.scrollY;
			window.history.replaceState(state, "", location.pathname + location.search + location.hash);
		}
	};

	this._emit = function(name, detail) {
		detail = detail || {};
		detail.lastUrl = mu._lastUrl;
		detail.previousUrl = mu._prevUrl;
		var ev = new CustomEvent(name, {
			bubbles: true,
			cancelable: true,
			detail: detail,
		});
		return (document.dispatchEvent(ev));
	};
};
window.mu = mu;

const isInternalURL = (url) => url.origin === window.location.origin;
const toSiteURL = (path) => new URL(path, window.location.href);
const isMuExcludedURL = (url) => isInternalURL(url) && (
	url.pathname === "/resume" ||
	url.pathname.startsWith("/resume/") ||
	url.pathname === "/archive" ||
	url.pathname.startsWith("/archive/")
);

function syncMuExcludedLinks(scope = document) {
	if (!scope?.querySelectorAll) {
		return;
	}

	for (const link of scope.querySelectorAll("a[href]")) {
		const href = link.getAttribute("href");
		if (!href) {
			continue;
		}

		let url;
		try {
			url = new URL(href, window.location.href);
		} catch {
			continue;
		}

		if (isMuExcludedURL(url)) {
			link.setAttribute("mu-disabled", "");
		}
	}
}

function syncDocumentChrome(doc) {
	const nextRoot = doc.documentElement;
	const nextBody = doc.body;
	const nextStyle = nextRoot.getAttribute("style");
	const nextPageColor = nextRoot.dataset.pageColor ?? "";

	if (nextStyle) {
		root.setAttribute("style", nextStyle);
	} else {
		root.removeAttribute("style");
	}

	if (nextPageColor) {
		root.dataset.pageColor = nextPageColor;
	} else {
		delete root.dataset.pageColor;
	}

	document.body.id = nextBody?.id || "";

	if (nextBody?.dataset?.route) {
		document.body.dataset.route = nextBody.dataset.route;
	} else {
		delete document.body.dataset.route;
	}

	if (nextBody?.dataset?.pageColor) {
		document.body.dataset.pageColor = nextBody.dataset.pageColor;
	} else {
		delete document.body.dataset.pageColor;
	}
}

function playSound(url, volume = 1) {
	const start = Date.now();
	if (!state.soundPlayers[url]) {
		state.soundPlayers[url] = new Audio(url);
	}

	const player = state.soundPlayers[url];
	player.volume = volume;

	const playIt = () => {
		if ((Date.now() - start) > 300) {
			return;
		}
		player.currentTime = 0;
		player.play();
	};

	if (player.readyState >= HTMLMediaElement.HAVE_CURRENT_DATA) {
		playIt();
		return;
	}

	player.addEventListener("canplaythrough", playIt, { once: true });
}

function navigateTo(path) {
	const url = toSiteURL(path);
	if (!isInternalURL(url) || isMuExcludedURL(url)) {
		window.location.assign(url.toString());
		return;
	}

	mu.load(url.pathname + url.search + url.hash);
}

async function loadSearchUI() {
	if (state.searchUI) {
		return state.searchUI;
	}

	if (!state.searchPromise) {
		state.searchPromise = import("/search.js")
			.then(({ createSearch }) => createSearch({ navigateTo, isMuExcludedURL }))
			.then((searchUI) => {
				state.searchUI = searchUI;
				return searchUI;
			});
	}

	return state.searchPromise;
}

function prefetchSearchUI() {
	void loadSearchUI();
}

function openSearch(query = "") {
	void loadSearchUI().then((searchUI) => searchUI.open(query));
}

function toggleSearch(query = "") {
	void loadSearchUI().then((searchUI) => searchUI.toggle(query));
}

function syncSearchFromLocation(options) {
	if (!currentURL().searchParams.get("q")) {
		return;
	}
	void loadSearchUI().then((searchUI) => searchUI.syncFromLocation(options));
}

// PAGE ENHANCEMENTS -------------------------------
function updateTOCState() {
	const links = document.querySelectorAll('#toc a[href^="#"]');
	const tocDiv = document.querySelector("#toc div:first-child");
	let current = null;

	for (const link of links) {
		const target = document.querySelector(link.getAttribute("href"));
		if (target?.getBoundingClientRect().bottom <= window.innerHeight / 4) {
			current = link;
		}
	}

	for (const link of links) {
		link.classList.toggle("scroll-current", link === current);
	}

	if (tocDiv?.scrollWidth > tocDiv?.clientWidth && current) {
		tocDiv.scrollTo({ top: 0, left: current.offsetLeft - 16, behavior: "smooth" });
	}
}

function initializePageContent() {
	for (const link of document.querySelectorAll(".content a:has(img)")) {
		if (link.dataset.zoomBound === "true") {
			continue;
		}
		if (!imageExtensions.has(link.href.split(".").pop())) {
			continue;
		}

		link.dataset.zoomBound = "true";
		link.addEventListener("click", (event) => event.preventDefault());
	}

	for (const image of document.querySelectorAll(".content :not(a):not(.no-zoom) img:not(.no-zoom)")) {
		if (image.dataset.zoomBound === "true") {
			continue;
		}

		image.dataset.zoomBound = "true";
		image.addEventListener("click", (event) => {
			const target = event.currentTarget;
			const parentHref = target?.parentElement?.href;

			if (parentHref) {
				const extension = parentHref.split(".").pop();
				if (!imageExtensions.has(extension)) {
					return;
				}

				if (target.src !== parentHref) {
					const imageStyle = getComputedStyle(target);
					target.style.width = imageStyle.width;
					target.style.height = imageStyle.height;
					target.addEventListener("load", () => target.click(), { once: true });
					target.src = parentHref;
					return;
				}
			}

			playSound("/sounds/stone2.ogg");
			if (target.classList.contains("zoomed")) {
				target.classList.remove("zoomed");
				target.style.removeProperty("transform");
				return;
			}

			target.classList.add("zoomed");
			if (!target.classList.contains("pixelated") && target.naturalWidth < (root.clientWidth / 2) && target.naturalHeight < (root.clientHeight / 2)) {
				target.classList.add("pixelated");
			}

			const bounds = target.getBoundingClientRect();
			const translateX = (root.clientWidth / 2) - (bounds.x + (bounds.width / 2));
			const translateY = (root.clientHeight / 2) - (bounds.y + (bounds.height / 2));
			const scale = Math.min(root.clientWidth / target.clientWidth, root.clientHeight / target.clientHeight, 20);
			target.style.transform = `translate(${translateX}px, ${translateY}px) scale(${scale})`;

			const unzoom = () => {
				target.classList.remove("zoomed");
				target.style.removeProperty("transform");
				document.removeEventListener("scroll", unzoomAfterScroll);
				document.removeEventListener("click", unzoom);
				playSound("/sounds/stone2.ogg");
			};

			const unzoomAfterScroll = () => {
				const latestBounds = target.getBoundingClientRect();
				if (latestBounds.top <= 0 || latestBounds.bottom >= window.innerHeight) {
					unzoom();
				}
			};

			document.addEventListener("scroll", unzoomAfterScroll);
			window.setTimeout(() => document.addEventListener("click", unzoom), 100);
		});
	}
}

function syncSearchFromLocation({ allowInstant = false } = {}) {
	const query = currentURL().searchParams.get("q");
	if (!query) {
		return;
	}

	openSearch(query);
	if (allowInstant) {
		document.querySelector("a.search-link.instant")?.click?.();
	}
}

function toggleTheme() {
	const currentDarkMode = getComputedStyle(document.body).getPropertyValue("color") === "rgb(255, 255, 255)";
	root.classList.add(currentDarkMode ? "light" : "dark");
	root.classList.remove(currentDarkMode ? "dark" : "light");
	playSound("/sounds/KDE_Click_2.ogg", 1);
}

// EVENTS ------------------------------------------
function handleDocumentClick(event) {
	const actionTarget = event.target.closest("#btn_larger, #btn_smaller, #btn_theme, #btn_toc, #btn_top, #btn_search, #search-link");
	if (!actionTarget) {
		return;
	}

	if (actionTarget.id === "search-link") {
		event.preventDefault();
	}

	switch (actionTarget.id) {
		case "btn_larger": {
			const content = currentContent();
			if (!content) {
				return;
			}
			const currentSize = getComputedStyle(content).fontSize;
			content.style.fontSize = `calc(${currentSize} * 1.0667)`;
			return;
		}
		case "btn_smaller": {
			const content = currentContent();
			if (!content) {
				return;
			}
			const currentSize = getComputedStyle(content).fontSize;
			content.style.fontSize = `calc(${currentSize} * 0.937)`;
			return;
		}
		case "btn_theme":
			toggleTheme();
			return;
		case "btn_toc": {
			const tocDetails = document.querySelector("#toc details");
			if (tocDetails) {
				tocDetails.open = !tocDetails.open;
			}
			return;
		}
		case "btn_top":
			window.scrollTo({ top: 0, behavior: "smooth" });
			return;
		case "btn_search":
		case "search-link":
			void toggleSearch();
			return;
	}
}

function handleSearchPrefetch(event) {
	if (event.target?.closest?.("#btn_search, #search-link")) {
		prefetchSearchUI();
	}
}

function handleKeyUp(event) {
	state.searchUI?.handleKeyUp(event);
}

function handleKeyDown(event) {
	if (event.key === "/" && document.activeElement?.tagName !== "INPUT") {
		event.preventDefault();
		void openSearch(window.getSelection().toString().replaceAll("\n", ""));
		return;
	}
	state.searchUI?.handleKeyDown(event);
}

function handleMuBeforeRender(event) {
	if (!event.detail?.html || event.detail.mode === "patch") {
		state.pendingMu = null;
		return;
	}

	state.pendingMu = parseHTML(event.detail.html);
}

function handleMuAfterRender() {
	if (state.pendingMu) {
		syncDocumentChrome(state.pendingMu);
		state.pendingMu = null;
	}

	syncMuExcludedLinks();
	initializePageContent();
	updateTOCState();
	syncSearchFromLocation();
}

function bindGlobalEvents() {
	if (state.globalEventsBound) {
		return;
	}

	state.globalEventsBound = true;
	document.addEventListener("scrollend", updateTOCState);
	document.addEventListener("click", handleDocumentClick);
	document.addEventListener("pointerover", handleSearchPrefetch);
	document.addEventListener("focusin", handleSearchPrefetch);
	document.addEventListener("keyup", handleKeyUp);
	document.addEventListener("keydown", handleKeyDown);
	document.addEventListener("mu:before-render", handleMuBeforeRender);
	document.addEventListener("mu:after-render", handleMuAfterRender);
}

// INIT --------------------------------------------
function load() {
	bindGlobalEvents();
	syncMuExcludedLinks();
	initializePageContent();
	mu.init();
	updateTOCState();

	const navigationType = window.performance.getEntriesByType("navigation")[0]?.type;
	syncSearchFromLocation({ allowInstant: navigationType === "navigate" });
}

load();
