+++
author = "Wesley Deal"
title = "Local Password Generator"
date = 2025-02-21
updated = 2025-07-18
raw = true
[taxonomies]
tags = ["software", "security"]
[extra]
shorttitle = "Password Generator"
color = "#00ba89"
+++

<style type="text/css">
#passwords-exposition p{
	margin-top: 0;
}
#settings-wrapper {
	display: flex;
	gap: 10px;
	flex-direction: column;
	flex: 1 1 300px;
}
#controls{
	display: flex;
	flex-direction: row;
	align-items: end;
	justify-content: space-between;
	flex-wrap: wrap;
	gap: 10px;
	max-width: 80ch;
}
select{
	border: 1px solid var(--color-primary);
	border-radius: 3px;
	font: inherit;
	color: inherit;
	background: none;
	padding: 2px;
}
input {
	accent-color: var(--color-pop);
	color: inherit;
}
#settings-wrapper div.length-container{
	display: flex;
	gap: 5px;
	flex-wrap: wrap;
}
#length{
	flex: 1 0 120px;
}
#length_val{
	width: 3em;
	background: transparent;
	border: none;
	border-bottom: 1px solid var(--color-fg);
	font: inherit;
	padding: 0;
	text-align: center;
}
input[type=number]::-webkit-inner-spin-button {
	opacity: 1;
}
.word-settings{
	display: none;
}
#genbutton-wrapper {
	flex: 1 0 200px;
}
@keyframes bgmove {
	from {background-position: 0 0}
	to {background-position: 0 0, 0 0, 800px 0, 0 0}
}
button#generate {
	font: 350 1.5em/1.5 var(--typeface-body), sans-serif;
	color: inherit;
	border: 2px solid var(--color-primary);
	border-radius: 3px;
	background: #000;
	background: var(--button-overlay, linear-gradient(#0000, #fff2)), var(--color-button-important, var(--color-primary));
	padding: 5px 15px;
	transition: border 150ms ease-out, color 150ms ease-out, background 150ms ease-out;
	width: 100%;
	user-select: none;
	box-shadow: 0 2px 5px color-mix(in srgb, var(--color-fg) 20%, transparent);
	box-shadow: var(--button-shadow);
	border: var(--button-border);
	color: #fff;
	text-shadow: 1px 1px 0 #000;
	transition: all 70ms ease-out;
}
button#generate:hover {
	//background: var(--color-pop);
	//color: var(--color-bg);
	//border: 2px solid var(--color-pop);
	background: var(--button-highlight-overlay, linear-gradient(#fff2, #fff4)), var(--color-button-important, var(--color-primary));
	box-shadow: var(--button-shadow-hover);
}
button#generate:active {
	background: var(--button-highlight-overlay, linear-gradient(#fff2, #fff4)), color-mix(in oklch, var(--color-pop, #999), #000 45%);
	animation: bgmove 15s infinite linear;
}
#copied-caption{
	opacity: 0;
	font-size: 1.5em;
	color: var(--color-pop);
	font-weight: bold;
	transition: opacity 250ms ease-in-out;
	border-radius: 5px;
	border: 2px solid;
	pointer-events: none;
	user-select: none;
	margin-top: 10px;
	text-align: center;
}
#copied-caption.active{
	opacity: 1;
}
#result{
	font-family: Consolas,monospace;
	font-family: var(--typeface-mono,Consolas,monospace);
	font-size: 1.5em;
	border: 1px solid var(--color-fg);
	background: color-mix(in srgb, var(--color-bg) 80%, transparent);
	overflow-x: auto;
}
#result.smaller{
	font-size: 1em;
}
#result ul{
	display: flex;
	flex-wrap: wrap;
	gap: 24px 20px;
	justify-content: center;
	padding: 0 20px;
}
#result li{
	display: flex;
	width: auto;
	padding: 0;
	margin: 0;
	line-height: 1;
	white-space: pre;
}
#result li.copied{
	border-bottom: 2px solid var(--color-accent);
	margin-bottom: -2px;
	font-style: italic;
}
</style>
<div id="passwords-exposition">
	<p>With these settings, we'll generate 50 passwords with <span id="entropy">0</span> bits of entropy which, if hashed with NTLM and brute forced with <a href="https://gist.github.com/Chick3nman/09bac0775e6393468c2925c1e1363d5c">RTX 5090</a>s, would take <span id="hashtime">0 sec</span> on average to guess.
</div>
<div id="controls">
	<div id="settings-wrapper">
		<div>
			<select id="mode" onchange="pwgen()">
				<option value="char">Random Characters</option>
				<option value="words" disabled>Random Words</option>
			</select>
		</div>
		<div class="length-container">
			<input type="range" id="length" min="8" max="128" value="16">
			<input type="number" id="length_val" value="16"></input>
   			<label for="length"> characters</label>
		</div>
		<div class="char-settings">
			<input type="checkbox" id="specials" value="specials" onchange="pwgen()">
			<label for="specials"> Special Characters</label><br>
			<input type="checkbox" id="ambiguous" value="ambiguous" onchange="pwgen()">
			<label for="ambiguous"> Ambiguous Characters</label><br>
			<input type="checkbox" id="spaces" value="spaces" onchange="pwgen()">
			<label for="spaces"> Spaces</label><br>
			<input type="checkbox" id="user-entropy" value="user-entropy" onchange="pwgen()">
			<label for="user-entropy"> Keyboard/Mouse Randomness</label><br>
		</div>
		<div class="word-settings">
			<input type="checkbox" id="separators" value="separators" onchange="pwgen()">
			<label for="separators"> Separators (-_+= )</label><br>
			<select id="capitalization">
				<option value="none">lowercase</option>
				<option value="all">ALL CAPS</option>
				<option value="first">First Letters</option>
				<option value="rand-first">random First letters</option>
				<option value="rand">raNDOm cAps</option>option>
			</select>
			<select id="numbers">
				<option value="none">No numbers</option>
				<option value="2">2 digits, start and/or end</option>
				<option value="4">4 digits, start and/or end</option>
				<option value="3">1-3 digits randomly placed</option>
			</select>
		</div>
	</div>
	<div id="genbutton-wrapper">
		<button type="button" id="generate">Generate <u>M</u>ore ↩</button>
	</div>
</div>

<div id="copied-caption">
	COPIED TO CLIPBOARD
</div>

<div id="result" class="fullbleed"></div>

<div id="alphabet-container">
	<p>from alphabet <code id="alphabet"></code></p>
</div>

<div id="information">
	<h2>About this tool</h2>
	<p>I use the JavaScript function <i><code>crypto.getRandomValues</code></i> to generate secure pseudorandom unsigned 32-bit integers in the range of [0, 4294967295].
		One Uint32 is used to pick each character of each password. On modern hardware, this is so blazingly fast that when you hold the button, the passwords update
		at the refresh rate of your display.
	<p>These passwords are generated entirely in your browser. You will note that this page contains no mechanism to send these passwords to any server.
		Generating 50 at once should act as a weak hedge against the possibility that the pseudorandom function on your system is compromised, because
		you can select both the time to generate passwords and the specific one you use at random. You can also opt to scramble the ordering of the alphabet used
		based on an array of numbers derived from the timestamps and values of your mouse and keyboard input.
		The main security advantage of this site for you is that if you view the page source, it should be short enough to parse and verify it does what it claims.
	<p>Password entropy is determined with the formula log₂(A<sup>L</sup>) where A is the alphabet and L is the length.
	<p>The time to crack is estimated by 2<sup>entropy</sup> ÷ hashRate ÷ 2. (We divide by 2 because, on average, the password will be found halfway through the search.)
	<p>For the hash rate, I've selected a ~$2000 GPU (RTX 5090) and a weak but common hash (NTLM) which is used by Windows to store local and network passwords.
		Ideally, your application should use a <em>much</em> better hash function which will take orders of magnitude longer to brute force.
</div>

<div id="license">
	<h2>License & Warranty Disclaimer</h2>
	<p>Permission to use, copy, modify, and/or distribute this software for
	any purpose with or without fee is hereby granted.
	<p>THE SOFTWARE IS PROVIDED “AS IS” AND THE AUTHOR DISCLAIMS ALL
	WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES
	OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE
	FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY
	DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN
	AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT
	OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
</div>

<script>
const maxUInt32 = (new Uint32Array([-1]))[0]
const ambiguousAlphaNumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
const unambiguousAlphaNumeric = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789";
const symbols = '`~!@#$%^&*()_+-=[]{}\\|;:\'",<.>/?';
const separators = "-_+=";

var copiedTimeout;

const genButton = document.getElementById("generate");
const startGenSound = () => {playSound('/sounds/KDE_Event_1.ogg', 0.5)}
var genButtonHeld = false;
genButton.addEventListener('pointerdown', (e) => { genButtonHeld = true; startGenSound(); });
genButton.addEventListener('pointerup', (e) => { genButtonHeld = false; });
genButton.addEventListener('pointerleave', (e) => { genButtonHeld = false; });
genButton.addEventListener('pointercancel', (e) => { genButtonHeld = false; });
document.body.addEventListener('keydown', (event) => { if(event.code == "Enter" || event.code == "Space" || event.code == "KeyM") genButtonHeld = true; startGenSound(); });
document.body.addEventListener('keyup', (event) => { if(event.code == "Enter" || event.code == "Space" || event.code == "KeyM") genButtonHeld = false; });


// Set userEntropy[last 3 milliseconds of timestamp] to pressed or released keycode or screen X or Y position
// this creates an array where fairly random elements are determined by user input.
// later we will loop through this to swap elements of the alphabet from userEntropyIndex % alphabetLength to userEntropyValue % alphabetLength
var userEntropy = new Uint32Array(1000);
function addUserEntropy(entropy){
	var dateEntropy = Number(Date.now().toString().slice(-3));
	userEntropy[dateEntropy] += entropy;
}
document.body.addEventListener('keydown', (event) => {addUserEntropy(event.key.charCodeAt());});
document.body.addEventListener('keyup', (event) => {addUserEntropy(event.key.charCodeAt() + 1);});
document.body.addEventListener('pointermove', (event) => {addUserEntropy(event.screenX + event.screenY);});

var results = document.getElementById("result");
var length = document.getElementById("length");
var length_val = document.getElementById("length_val");
function onChangeLength(event) {
	if (event.target.value > 16) { // set password output font size
		results.classList.add("smaller");
	} else {
		results.classList.remove("smaller");
	}
	if (event.target == document.getElementById('length')){ // update the other HTML control to match
		document.getElementById('length_val').value=event.target.value;
	} else {
		document.getElementById('length').value=event.target.value;
	}
	pwgen();
}
function scrollAdjustLength(event){
	length.value -= Math.sign(event.deltaY);
	onChangeLength(event);
	event.preventDefault();
	event.stopPropagation();
}
length_val.addEventListener('wheel', scrollAdjustLength);
length_val.addEventListener('keyup', onChangeLength);
length_val.addEventListener('mouseup', onChangeLength);
length_val.addEventListener('pointerup', () => length_val.select());
length.addEventListener('input', onChangeLength);
length.addEventListener('wheel', scrollAdjustLength);

function secureRand(min, max, count=1) {
	var randInts = crypto.getRandomValues(new Uint32Array(count));
	var scale = max-min;
	var result = min;
	return randInts.map(r => min + Math.round(scale*(r/maxUInt32)));
}

function pwgen() {
	var result = "";
	const modeElement = document.getElementById("mode");
	const mode = modeElement.options[modeElement.selectedIndex].value;
	const pwlen = parseInt(document.getElementById("length_val").value);
	const specials = document.getElementById("specials").checked;
	const ambiguous = document.getElementById("ambiguous").checked;
	const spaces = document.getElementById("spaces").checked;
	const count = 50;

	// prepare alphabet
	var charSet = ambiguous ? ambiguousAlphaNumeric : unambiguousAlphaNumeric;
	if(specials) {charSet += symbols;}
	if(spaces) {charSet += " ";}
	charSet = Array.from(charSet);
	if (document.getElementById("user-entropy").checked) {
		addUserEntropy(charSet.length+spaces+ambiguous+specials+pwlen);
		for (i=0; i<userEntropy.length; i++) {
			if (userEntropy[i]){
				var a = i % charSet.length;
				var b = userEntropy[i] % charSet.length;
				[charSet[a], charSet[b]] = [charSet[b], charSet[a]];
			}
		}
	}

	var resultInnerHTML = "<ul>";
	var randChars = Array.from(secureRand(0, charSet.length-1, pwlen*count));
	for (i=0; i<count; i++) {
		var pw = randChars.slice(i*pwlen,i*pwlen+pwlen).map(c => charSet[c]).join("");
		resultInnerHTML += "<li>" + pw.replaceAll("&","&amp;").replaceAll(">","&gt;").replaceAll("<","&lt;").replaceAll("\\",'&bsol;') + "</li>";
	}
	resultInnerHTML += "</ul>"
	document.getElementById("result").innerHTML = resultInnerHTML;

	for (el of document.querySelectorAll("#result li")){ // auto copy password on click
		el.addEventListener('click', function(e) {
			window.getSelection().removeAllRanges();
			r = document.createRange();
			r.selectNodeContents(e.target.closest('li'));
			window.getSelection().addRange(r);
			document.execCommand('copy');
			e.target.closest('li').classList.add('copied');
			document.getElementById("copied-caption").classList.add("active");
			window.clearTimeout(copiedTimeout);
			copiedTimeout = window.setTimeout(() => {document.getElementById("copied-caption").classList.remove("active");}, 1000);
			playSound('/sounds/Shot_Reverb.ogg', 0.5);
		});
	}

	// determine and display time to crack
	const hashRate = 3.401e+11;
	var entropy = Math.log2(Math.pow(charSet.length, pwlen));
	var hashTime = Math.pow(2, entropy) / hashRate / 2 / 3600 / 24 / 365;
	if (hashTime < 1/365/24) {
		hashTimeText = (hashTime*365*24*3600).toPrecision(3) + " GPU-seconds";
	} else if (hashTime < 1/365) {
		hashTimeText = (hashTime*365*24).toPrecision(3) + " GPU-hours";
	} else if (hashTime < 1) {
		hashTimeText = (hashTime*365).toPrecision(3) + " GPU-days";
	} else if (hashTime < 1000000) {
		hashTimeText = hashTime.toLocaleString("en-US", {notation: "compact", maximumSignificantDigits: 3})  + " GPU-years";
	} else if (hashTime < 1000000000) {
		hashTimeText = (hashTime / 1000000).toPrecision(3) + " million GPU-years";
	} else if (hashTime < 1000000000000) {
		hashTimeText = (hashTime / 1000000000).toPrecision(3) + " billion GPU-years";
	} else if (hashTime < 1000000000000000) {
		hashTimeText = (hashTime / 1000000000000).toPrecision(3) + " trillion GPU-years";
	} else {
		hashTimeText = hashTime.toPrecision(3) + " GPU-years";
	}
	document.getElementById("entropy").innerHTML = entropy.toPrecision(3);
	document.getElementById("hashtime").innerHTML = hashTimeText;
	document.getElementById("alphabet").innerHTML = charSet.join("").replaceAll("&","&amp;").replaceAll(">","&gt;").replaceAll("<","&lt;").replaceAll("\\",'&bsol;');
}

function monitorGenButton() { // generate passwords each frame
	if (genButtonHeld) pwgen();
	window.requestAnimationFrame(monitorGenButton);
}

monitorGenButton();
addUserEntropy(window.screen.availHeight); // init user based entropy with screen height
addUserEntropy(window.screen.availWidth); // and width
addUserEntropy(navigator.deviceMemory); // and the user's device memory
Array.from(navigator.userAgent).map(c => addUserEntropy(c.charCodeAt())); // and the character codes of the user agent
// (all of which are predictable, but at worst this is better than no scrambling)
pwgen();
</script>
