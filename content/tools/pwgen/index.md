+++
author = "Wesley Deal"
title = "Local Password Generator"
date = 2025-02-21
updated = 2025-02-21
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
	border: 1px dashed #52abab;
	margin: 10px 0;
	overflow-x: auto;
	max-width: 100%;
}
#result.smaller{
	font-size: 1em;
}
#result ul{
	display: flex;
	flex-wrap: wrap;
	gap: 10px 20px;
	justify-content: space-evenly;
	padding: 0;
}
#result li{
	display: flex;
	width: auto;
}
#result li.copied{
	color: #5cc;
	font-style: italic;
}
.content{
	display: flex;
	max-width: 100%;
	align-items: center;
	flex-direction: column;
	gap: 20px 0;
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
			<select id="mode">
				<option value="char">Random Characters</option>
				<option value="words">Random Words (NOT YET IMPLEMENTED)</option>
			</select>
		</div>
		<div class="length-container">
			<input type="range" id="length" min="8" max="128" value="16">
			<input type="number" id="length_val" value="16" onmousewheel="document.getElementById('length').value=this.value; pwgen()"></input>
   			<label for="length"> characters</label>
		</div>
		<div>
			<input type="checkbox" id="specials" value="specials" onchange="pwgen()">
			<label for="specials"> Special Characters</label><br>
			<input type="checkbox" id="ambiguous" value="ambiguous" onchange="pwgen()">
			<label for="ambiguous"> Ambiguous Characters</label><br>
			<input type="checkbox" id="separators" value="separators">
			<label for="separators"> Separators (-_+=)  (NOT YET IMPLEMENTED)</label><br>
			<input type="checkbox" id="spaces" value="spaces">
			<label for="spaces"> Spaces</label><br>
		</div>
	</div>
	<div id="genbutton-wrapper">
		<button type="button" id="generate">Generate <u>M</u>ore ↩</button>
	</div>
</div>

</form>

<div id="explanation">30 passwords, freshly baked with <code>crypto.getRandomValues</code> on your local computer, are ready below. Click to instantly copy to your clipboard.</div>

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

const maxUInt32 = (function(){
	var toUnderflow = new Uint32Array(1);
	toUnderflow[0] -= 1;
	return toUnderflow[0];
})();
const ambiguousAlphaNumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
const unambiguousAlphaNumeric = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789";
const symbols = '`~!@#$%^&*()_+-=[]{}\|;:\'",<.>/?';
const separators = "-_+=";
	
var genButton = document.getElementById("generate");
var genButtonHeld = false;
genButton.addEventListener('pointerdown', (event) => { genButtonHeld = true; })
genButton.addEventListener('pointerup', (event) => { genButtonHeld = false; })
genButton.addEventListener('pointerleave', (event) => { genButtonHeld = false; })
genButton.addEventListener('pointercancel', (event) => { genButtonHeld = false; })
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
length_val.addEventListener('change', onChangeLength);
length.addEventListener('input', onChangeLength);
length_val.addEventListener('pointerup', () => length_val.select());
function secureRand(min, max) {
	var [randInt] = crypto.getRandomValues(new Uint32Array(1));
	var scale = max-min;
	var result = min;
	return Math.round(min + (scale*(randInt/maxUInt32)));
}
function pwgen() {
	var result = "";
	var modeElement = document.getElementById("mode");
	var mode = modeElement.options[modeElement.selectedIndex].value;
	var pwlen = parseInt(document.getElementById("length_val").value);
	var specials = document.getElementById("specials").checked;
	var ambiguous = document.getElementById("ambiguous").checked;
	var spaces = document.getElementById("spaces").checked;

	var charSet = ambiguous ? ambiguousAlphaNumeric : unambiguousAlphaNumeric;
	if(specials) {charSet += symbols;}
	if(spaces) {charSet += " ";}
	
	var resultInnerHTML = "<ul>";
	for (i=0; i<30; i++) {
		resultInnerHTML += "<li>" + Array.from({length: pwlen}, () => charSet[secureRand(0,charSet.length - 1)]).join("").replaceAll(">","&gt;").replaceAll("<","&lt;").replaceAll("\n","\\n");
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
}
function monitorGenButton() {
	if (genButtonHeld) pwgen();
	window.requestAnimationFrame(monitorGenButton);
}
monitorGenButton();
pwgen();

</script>
