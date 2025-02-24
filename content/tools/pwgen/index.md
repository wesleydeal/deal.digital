+++
author = "Wesley Deal"
title = "Local Password Generator"
date = 2025-02-21
updated = 2025-02-24
raw = true
[taxonomies]
tags = ["software", "security"]
[extra]
shorttitle = "Password Generator"
+++

<style type="text/css">
#length_val{
	width: 3em;
	background: transparent;
	border: none;
	border-bottom: 1px solid var(--color-fg);
	border-radius: 0;
	font: inherit;
	padding: 0;
	text-align: center;
}
input[type=number]::-webkit-inner-spin-button {
    opacity: 1
}
#settings-wrapper {
	display: flex;
	gap: 10px;
	flex-direction: column;
}
#settings-wrapper div.length-container{
	display: flex;
	gap: 5px;
	flex-wrap: wrap;
}
#length{
	width: 256px;
}
#result{
	font-family: "Red Hat Mono",Consolas,monospace;
	font-size: 1.5em;
	border: 1px solid var(--color-fg);
	outline: 2px solid var(--color-accent);
	margin: 10px 0;
	overflow-x: auto;
	max-width: min(100%,1600px);
}
#result.smaller{
	font-size: 1em;
}
#result ul{
	display: flex;
	flex-wrap: wrap;
	gap: 24px 20px;
	justify-content: space-evenly;
	padding: 0 10px;
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
.content{
	display: flex;
	max-width: 100%;
	align-items: center;
	flex-direction: column;
	gap: 20px 0;
	min-height: 100vh;
}
.content > *:not(div#result) {
	width: calc(min(780px,100%));
	padding: 0 10px;
}
button#generate {
	font-size: 1.5em;
	color: var(--color-link);
	border: 2px solid var(--color-link);
	border-radius: 3px;
	background: var(--color-bg);
	padding: 5px 15px 10px;
	transition: border 250ms ease-in-out, color 250ms ease-in-out, background 250ms ease-in-out;
}
button#generate:hover {
	background: var(--color-hover);
	color: var(--color-bg);
	border: 2px solid var(--color-hover);
}
#controls{
	display: flex;
	flex-direction: row;
	align-items: end;
	justify-content: space-between;
	flex-wrap: wrap;
	gap: 10px;
}
select{
	border: 1px solid var(--color-fg);
	border-radius: 3px;
	font: inherit;
	background: none;
	padding: 2px;
}
</style>
<div id="controls">
	<div id="settings-wrapper">
		<div>
			<select id="mode" onchange="pwgen()">
				<option value="char">Random Characters</option>
				<option value="words">Random Words (NOT YET IMPLEMENTED)</option>
			</select>
		</div>
		<div class="length-container">
			<input type="range" id="length" min="8" max="128" value="16">
			<input type="number" id="length_val" value="16"></input>
   			<label for="length"> characters</label>
		</div>
		<div>
			<input type="checkbox" id="specials" value="specials" onchange="pwgen()">
			<label for="specials"> Special Characters</label><br>
			<input type="checkbox" id="ambiguous" value="ambiguous" onchange="pwgen()">
			<label for="ambiguous"> Ambiguous Characters</label><br>
			<input type="checkbox" id="spaces" value="spaces" onchange="pwgen()">
			<label for="spaces"> Spaces</label><br>
			<input type="checkbox" id="separators" value="separators" onchange="pwgen()" disabled>
			<label for="separators"> Separators (-_+=)  (NOT YET IMPLEMENTED)</label><br>
		</div>
	</div>
	<div id="genbutton-wrapper">
		<button type="button" id="generate">Generate <u>M</u>ore ↩</button>
	</div>
</div>

<div id="explanation">
	<p>50 passwords, freshly baked with <code>crypto.getRandomValues</code> on your local computer, are ready below. Click to instantly copy to your clipboard.
	<p>These password have <span id="entropy">0</span> bits of entropy and would take <span id="hashtime">0 sec</span> on average to guess with an <a href="https://gist.github.com/Chick3nman/09bac0775e6393468c2925c1e1363d5c">NVIDIA RTX 5090 if hashed with NTLM</a>.
</div>

<div id="result"></div>

<div id="license"><p>Permission to use, copy, modify, and/or distribute this software for
any purpose with or without fee is hereby granted.
<p>THE SOFTWARE IS PROVIDED “AS IS” AND THE AUTHOR DISCLAIMS ALL
WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES
OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE
FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY
DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN
AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT
OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.</div>

<script>
const maxUInt32 = (new Uint32Array([-1]))[0]
const ambiguousAlphaNumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
const unambiguousAlphaNumeric = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789";
const symbols = '`~!@#$%^&*()_+-=[]{}\\|;:\'",<.>/?';
const separators = "-_+=";
	
var genButton = document.getElementById("generate");
var genButtonHeld = false;
genButton.addEventListener('pointerdown', (e) => { genButtonHeld = true; })
genButton.addEventListener('pointerup', (e) => { genButtonHeld = false; })
genButton.addEventListener('pointerleave', (e) => { genButtonHeld = false; })
genButton.addEventListener('pointercancel', (e) => { genButtonHeld = false; })
document.body.addEventListener('keydown', (event) => { if(event.code == "Enter" || event.code == "KeyM") genButtonHeld = true; })
document.body.addEventListener('keyup', (event) => { if(event.code == "Enter" || event.code == "KeyM") genButtonHeld = false; })

var results = document.getElementById("result");
var length = document.getElementById("length");
var length_val = document.getElementById("length_val");
function onChangeLength(event) {
	if (event.target.value > 16) {
		results.classList.add("smaller");
	} else {
		results.classList.remove("smaller");
	}
	if (event.target == document.getElementById('length')){
		document.getElementById('length_val').value=event.target.value;
	} else {
		document.getElementById('length').value=event.target.value;
	}
	pwgen();
}
length_val.addEventListener('wheel', scrollLength);
length_val.addEventListener('keyup', onChangeLength);
length_val.addEventListener('mouseup', onChangeLength);
length_val.addEventListener('pointerup', () => length_val.select());
length.addEventListener('input', onChangeLength);
length.addEventListener('wheel', scrollLength);

function scrollLength(event){
	length.value -= Math.sign(event.deltaY);
	length_val.value = length.value;
	pwgen();
	event.preventDefault();
	event.stopPropagation();
}
function secureRand(min, max, count=1) {
	var randInts = crypto.getRandomValues(new Uint32Array(count));
	var scale = max-min;
	var result = min;
	return randInts.map(r => min + Math.round(scale*(r/maxUInt32)));
}
function roundTo(value, decimals) {
	const m = Math.pow(10, decimals);
	return Math.round(m * value) / m;
}
function pwgen() {
	var result = "";
	var modeElement = document.getElementById("mode");
	var mode = modeElement.options[modeElement.selectedIndex].value;
	var pwlen = parseInt(document.getElementById("length_val").value);
	var specials = document.getElementById("specials").checked;
	var ambiguous = document.getElementById("ambiguous").checked;
	var spaces = document.getElementById("spaces").checked;
	const count = 50;
	var charSet = ambiguous ? ambiguousAlphaNumeric : unambiguousAlphaNumeric;
	if(specials) {charSet += symbols;}
	if(spaces) {charSet += " ";}
	var resultInnerHTML = "<ul>";
	var randChars = Array.from(secureRand(0, charSet.length-1, pwlen*count));
	for (i=0; i<count; i++) {
		var pw = randChars.slice(i*pwlen,i*pwlen+pwlen).map(c => charSet[c]).join("");
		resultInnerHTML += "<li>" + pw.replaceAll("&","&amp;").replaceAll(">","&gt;").replaceAll("<","&lt;").replaceAll("\\",'&bsol;') + "</li>";
	}
	resultInnerHTML += "</ul>"
	document.getElementById("result").innerHTML = resultInnerHTML;
	for (el of document.querySelectorAll("#result li")){
		el.addEventListener('click', function(e) {
			window.getSelection().removeAllRanges();
			r = document.createRange();
			r.selectNodeContents(e.target.closest('li'));
			window.getSelection().addRange(r);
			document.execCommand('copy');
			e.target.closest('li').classList.add('copied');
		});
	}
	const hashRate = 3.401e+11;
	var entropy = Math.log2(Math.pow(charSet.length, pwlen));
	var hashTime = Math.pow(2, entropy) / hashRate / 2 / 3600 / 24 / 365;
	var hashTimeText = "";
	if (hashTime < 1/365/24) {hashTimeText = roundTo(hashTime*365*24*3600, 0) + " seconds"} else
	if (hashTime < 1/365) {hashTimeText = roundTo(hashTime*365*24, 1) + " hours"} else
	if (hashTime < 1) { hashTimeText = roundTo(hashTime*365, 1) + " days"} else
	hashTimeText = roundTo(hashTime, 2) + " years";
	document.getElementById("entropy").innerHTML = roundTo(entropy, 2);
	document.getElementById("hashtime").innerHTML = hashTimeText;
}
function monitorGenButton() {
	if (genButtonHeld) pwgen();
	window.requestAnimationFrame(monitorGenButton);
}
monitorGenButton();
pwgen();

</script>
