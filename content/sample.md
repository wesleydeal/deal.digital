+++
title = "Sample Page"
[extra]
toc = true
+++

This page is here to demonstrate all of the common elements used on this website. This is primarily intended for
my personal use during stylesheet development. Oh, and by the way, [here is a link](#) that goes nowhere.

## Colors
Set the **primary color** for this page using this form: <input id="primary-color-picker" type="color">

| Color Name | Color Value |
|:-|:-|
| \--color-primary | |
| \--color-primary-sat | |
| \--color-pop | |
| \--color-pop-bgcontrast | |
| \--color-h1 | |
| \--color-fg | |
| \--color-bg | |

<script>
document.getElementById('primary-color-picker').addEventListener('input', setColors);

function setColors() {
	let picker = document.getElementById('primary-color-picker')
	document.documentElement.style.setProperty('--color-primary',picker.value);
	let table = picker.parentElement.nextElementSibling;
	let demoRows = table.querySelectorAll('tbody tr');
	
	for (row of demoRows){
		let nameCell = row.querySelectorAll('td')[0]
		let valCell = row.querySelectorAll('td')[1]
		let varName = nameCell.innerHTML;
		let val = getComputedStyle(document.documentElement).getPropertyValue(varName);
		nameCell.style.setProperty('color', val);
		nameCell.style.setProperty('text-shadow', '1px 1px 0 var(--color-fg)');
		valCell.innerHTML = val;
		valCell.style.setProperty('background', val);
		valCell.style.setProperty('text-shadow', '1px 1px 0 var(--color-bg)');
	}
}

const rgbToHex = (r, g, b) => {
  return "#" + ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1);
}

(() => {
	let color = getComputedStyle(document.querySelector('header')).getPropertyValue('border-bottom-color');
	let [r,g,b] = color.substring(color.indexOf('(')+1, color.indexOf(')')).split(", ");
	let picker = document.getElementById('primary-color-picker')
	picker.value = rgbToHex(r*1,g*1,b*1);
	setColors();
})();
</script>

## Code

There's two primary ways to write code. You can use inline code by quoting it with ``` `single backticks` ```
or a multiline ``` <pre><code> ``` surrounded by triple backticks, optionally placing a language on the first
line after the backticks. For example, here's some HTML.

```html
<!doctype html>
<html lang="en">
<meta charset="utf-8">
<title>Hello, world!</title>
<main>
	<h1>Hello, world.</h1>
	<p>This is a complete sample document within my sample document. A meta sample, so to speak. And it has a rather wide line of text.
</main>
```

## Blockquotes

Sometimes, you want a blockquote that doesn't quote anyone.

> Who wrote this?

But more often, you'll want to attribute the quote to an author (or even *--Anonymous*). Here's what that would look like, with a longer quote.

# Here is a long level 1 heading, just as a demonstration of how the text wraps
Below this, I'm going to put in some lower level headings. I don't generally use levels below an &lt;h3&gt;
## Heading level 2
### Heading level 3
#### Heading level 4
##### Heading level 5