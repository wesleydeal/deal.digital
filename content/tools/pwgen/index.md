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
	width: 50px;
}
form {
	display: flex;
	gap: 5px;
	flex-direction: column;
}
form div{
	padding: 5px 10px;
}
form div.lengthfields{
	display: flex;
	gap: 5px;
}
#length{
	width: 256px;
}
#result{
	font-family: "Red Hat Mono",Consolas,monospace;
	font-size: 1.5em;
	border: 1px dashed #52abab;
	margin: 10px 0;
}
#result.smaller{
	font-size: 1em;
}
#result ul{
	display: flex;
	flex-wrap: wrap;
	gap: 10px 20px;
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
	max-width: initial;
	align-items: center;
	flex-direction: column;
	gap: 20px 0;
}
.content > *:not(div#result) {
	width: calc(min(780px,100%));
}
</style>
<form>
	<div>
		<select id="mode">
			<option value="char">Random Characters</option>
			<option value="words">Random Words (NOT YET IMPLEMENTED)</option>
		</select>
	</div>
	<div>
		<label for="length">Length: </label>
		<input type="range" id="length" min="1" max="128" onchange="document.getElementById('length_val').value=this.value; pwgen()" value="16">
		<input type="number" id="length_val" value="16" onchange="document.getElementById('length').value=this.value; pwgen()" onmousewheel="document.getElementById('length').value=this.value; pwgen()"></input>
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
	<div>
		<button type="button" id="generate">More</button>
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
genButton.addEventListener('touchstart', (event) => { genButtonHeld = true; })
genButton.addEventListener('touchend', (event) => { genButtonHeld = false; })
genButton.addEventListener('mousedown', (event) => { genButtonHeld = true; })
genButton.addEventListener('mouseup', (event) => { genButtonHeld = false; })

var results = document.getElementById("result");
var length = document.getElementById("length");
length.addEventListener('change', (event) => {
	if (length.value > 16) {
		length.classList.add("smaller");
	} else {
		length.classList.remove("smaller");
	}
});

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
	var pwlen = parseInt(document.getElementById("length").value);
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
