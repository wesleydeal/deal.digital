+++
title = "Sample Page"
[extra]
toc = true
+++

This page is here to demonstrate all of the common elements used on this website. This is primarily intended for
my personal use during stylesheet development. Oh, and by the way, [here is a link](#) that goes nowhere.

This paragraph contains ***bold, *bold italic (some of which is [linked](#)),*** plain italic,* and ~strikethrough~ text.

## Colors
Set the **primary color** for this page using this or the field in the table below: <input id="primary-color-picker" type="color">

| Color Name | ... Color Preview ... | Color String |
|:-|:-:|:-|
| \--color-primary | Sample | <input id="primary-color-picker-text">
| \--color-primary-sat | Sample |
| \--color-pop | Sample |
| \--color-pop-bgcontrast | Sample |
| \--color-h1 | Sample |
| \--color-fg | Sample |
| \--color-bg | Sample |

<script>
document.getElementById('primary-color-picker').addEventListener('input', clickColor);
document.getElementById('primary-color-picker-text').addEventListener('input', typeColor);

function updateColorTable() {
	let textInput = document.getElementById('primary-color-picker-text');
	let table = textInput.parentElement.parentElement.parentElement;
	let demoRows = table.querySelectorAll('tbody tr');
	
	for (row of demoRows){
		let nameCell = row.querySelectorAll('td')[0]
		let previewCell = row.querySelectorAll('td')[1]
		let valCell = row.querySelectorAll('td')[2]
		let varName = nameCell.innerHTML;
		let val = getComputedStyle(document.documentElement).getPropertyValue(varName);
		nameCell.style.setProperty('color', val);
		nameCell.style.setProperty('text-shadow', '1px 1px 0 var(--color-fg)');
		if(varName == "--color-primary") {
			if (document.activeElement != textInput) {
				textInput.value = val;
			}
		} else {
			valCell.innerHTML = val;
		}
		previewCell.style.setProperty('background', val);
	}
}

function clickColor() {
	let picker = document.getElementById('primary-color-picker');
	document.documentElement.style.setProperty('--color-primary',picker.value);
	updateColorTable();
}

function typeColor() {
	let val = document.getElementById('primary-color-picker-text').value
	if (!CSS.supports('color',val)) {
		return;
	}
	document.documentElement.style.setProperty('--color-primary', val);
	updateColorTable();
}

updateColorTable();
</script>
<style>
:has(#primary-color-picker) + table tbody tr td:nth-child(2) {
	
}
</style>

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
	<p>This is a complete sample HTML document within my sample document. A meta sample, so to speak. And it has a rather ~w~i~d~e~ line of text. A line of text that might compress the size of the left sidebar if you're on a wide enough screen. Actually, that may be a bit of a bug. The line of text in question, the one you're currently reading, is so wide that even on a 4k monitor set to 96 DPI, you should still have to scroll just a little bit in order to see the whole thing.
</main>
```

## Blockquotes

Sometimes, you want a blockquote that doesn't quote anyone.

> Who wrote this?

But more often, you'll want to attribute the quote to an author (or even *--Anonymous*). Here's what that would look like, with a longer quote.

{% blockquote(author="Warren Buffet",date="2025-06-23") %}
Freak the fuck out and panic sell everything right now.

It's fucking over.
{% end %}

# Domain *Eukaryota*
Organisms whose cells possess a membrane-bound nucleus.

## Kingdom *Animalia*
Multicellular eukaryotes which generally consume organic material, breathe oxygen, move, and reproduce sexually.

### Phylum *Chordata*
Animals possessing a notochord at some stage in development.

#### Class *Mammalia*
Warm-blooded vertebrates with hair/fur and mammary glands.

##### Order *Carnivora*
Meat-eating mammals with specialized teeth and claws.

###### Family *Felidae*
The biological family of cats, known for retractable claws and keen predation skills.

**Snow Leopard**, *Panthera uncia*  
The snow leopard — a solitary, mountain-dwelling big cat native to Central and South Asia.

