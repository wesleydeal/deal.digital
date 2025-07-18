+++
title = "Sample Page"
[extra]
toc = true
+++

This page is here to demonstrate all of the common elements used on this website. This is primarily intended for
my personal use during stylesheet development. Oh, and by the way, [here is a link](#) that goes nowhere.

This paragraph contains **bold**, ***bold italic (some of which is [linked](#))***, *plain italic,* and ~strikethrough~ text.

This website is best viewed with any browser.

<img src="/badges/anydamn.gif" class="badge"><img src="/badges/right2repair.png" class="badge">

## Colors
Set the **primary color** for this page using this or the field in the table below: <input id="primary-color-picker" type="color">

| Color/Gradient Name | ... Color Preview ... | Color String |
|:-|:-:|:-|
| \--color-primary | Sample | <input id="primary-color-picker-text">
| \--color-primary-sat | Sample | ` `
| \--color-pop | Sample | ` `
| \--color-pop-bgcontrast | Sample | ` `
| \--color-link | Sample | ` `
| \--color-link-underline | Sample | ` `
| \--color-h1 | Sample | ` `
| \--color-h2 | Sample | ` `
| \--color-fg | Sample | ` `
| \--color-fg-hc | Sample | ` `
| \--color-bg | Sample | ` `
| \--color-bg-hc | Sample | ` `
| \--color-tinted-shadow | Sample | ` `
| \--bg | Sample | ` `
| \--header-bg | Sample | ` `

<script>
document.getElementById('primary-color-picker').addEventListener('input', clickColor);
document.getElementById('primary-color-picker-text').addEventListener('input', typeColor);

function updateColorTable() {
	let textInput = document.getElementById('primary-color-picker-text');
	let table = textInput.parentElement.parentElement.parentElement;
	let demoRows = table.querySelectorAll('tbody tr');

	for (row of demoRows){
		let nameCell = row.querySelectorAll('td')[0];
		let previewCell = row.querySelectorAll('td')[1];
		let valCell = row.querySelectorAll('td')[2];
		valCell = valCell.querySelector('code') || valCell;
		let varName = nameCell.innerHTML;
		let val = getComputedStyle(document.body).getPropertyValue(varName);
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

window.setTimeout(updateColorTable, 500);
</script>
<style>
:has(#primary-color-picker) + table tbody tr td:nth-child(2) {

}
</style>

## Unordered Lists
* typically comprised of bullet points
* can have multiple levels
  - each of which can have more levels
    * but more than 3 is inadvisable
  - which can be ordered
    1. though this can be a bit confusing
	2. typically begins with numbers rather than letters

## Ordered Lists
1. Typically numbered
2. Can have multiple levels
   1. which often alternate between letters and capital and lowercase numbers
      - in markdown, we still express these with numbers
3. Are used for situations such as:
   1. Requirements, for
      1. software projects
	  2. budget planning
   2. Lists of instructions that need to be followed from beginning to end, and which the reader may need to refer back to several times during the operation without losing their place.

## Images
This is an svg:
![svg diagram](/tools/tpmiddle/diagram.svg)

This is a normal image which links to a page:

[![sample image of librechat](/blog/librechat/librechat-example.png)](/blog/librechat/)

This is a reduced resolutiom image which links to a higher resolution version:

{{ fitimg(path='/blog/librechat/credentials.png', width=100) }}

## Code and \<pre\>

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

And here is some preformatted text, specifically a table with ascii box drawing characters, taken from [Suraj Kurapati](https://www.mail-archive.com/markdown-discuss@six.pairlist.net/msg01650.html)

<pre>
                     A more complex table-within-a table.
               An inner table showing a variety of headings and data items.
            ┌────────────────────────────────────────────────────────────────┐
            │                          Inner Table                           │
            ├───────────────────┬────────────────────────────────────────────┤
            │                   │                   Head1                    │
            │      CORNER       ├────────────────────────────┬───────────────┤
            │                   │                            │     Head3     │
            ├───────────┬───────┤           Head2            ├──────┬────────┤
            │   Head4   │ Head5 │                            │Head6 │        │
   Outer    ├───────────┼───────┼────────────────────────────┼──────┴────────┤
   Table    │A          │       │  • Lists can be table data │   Two Wide    │
            │           │       │  • Images can be table data│               │
            ├───────────┤       ├────────────────────────────┼──────┬────────┤
            │           │Two    │A Form in a table: Your age:│      │        │
            │           │Tall   │            [  ]            │  No  │Multiple│
            │HTML       │       │ What is your favorite ice  │border│line    │
            │Station    │       │           cream?           │Little│item    │
            │           │       │    [Chocolate         ]    │Table │        │
            │           │       │       [OK] [Cancel]        │      │        │
            └───────────┴───────┴────────────────────────────┴──────┴────────┘
</pre>

Sometimes a program has a sample output, which we represent with \<samp\> usually enclosed in \<pre\>. Here's the first few lines of what happens of you press <kbd>Ctrl</kbd> + <kbd>Shift</kbd> + <kbd>Alt</kbd> + <kbd>?</kbd>:

<pre><samp>Basic editing functions
=======================
Enter            AcceptLine          Accept the input or move to the next line if input is missing a closing token.
Shift+Enter      AddLine             Move the cursor to the next line without attempting to execute the input
Backspace        BackwardDeleteChar  Delete the character before the cursor
Ctrl+h           BackwardDeleteChar  Delete the character before the cursor
[...output truncated...]
</samp></pre>

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
