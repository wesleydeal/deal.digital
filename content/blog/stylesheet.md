+++
title = "I rewrote my stylesheet"
date = 2099-01-01
draft = true
[extra]
shorttitle = "New Stylesheet"
color = "teal"
+++
I'm no trained designer; I don't know many of the words, the principles, the tricks of the trade.
I have certainly made a lot of graphic design choices at work, not all of them well-reasoned.
But I find myself constantly fascinated by shape, spacing, contrast, and type, both in print and
on screen.

Neither do I consider myself much of a coder. I have opinions on some programming concepts, I read Hacker News,
and I'll exercise my fingers typing a sub-1000 line script.

But sometimes...

## Inspiration Strikes

I'm enthralled with [pre-2010s computer interfaces](http://toastytech.com/guis/index.html) and [themes](https://macthemes.garden/).

As graphical computing developed, there was such a diversity of design and interface paradigms. Low resolution
and color depth created tight constraints for designers, and with them beautiful artifacts such as dither patterns,
bitmap fonts, and almost-but-not-quite analogues of real world objects. Because the field was new, they were free to
experiment with different UX patterns.

Even into the early 2000s as UI/UX design fields began to mature, amateur and professional designs pushed the bounds
of creativity and [sanity](https://wmpskinsarchive.neocities.org/). Early Mac OS X versions were adorned with pinstripes
and [brushed metal](https://cdsassets.apple.com/live/6GJYWVAV/user/ma543_macosx10_3_welcome.pdf); web "2.0" buttons
got gradients, highlights, shadows, reflections, and stripes; music players [adorned with aesthetics and characters
from every form of media](https://skins.webamp.org/). Windows saw hard edges transformed into slick [watercolor](https://www.betaarchive.com/wiki/index.php?title=File:Windows_Whistler_2257_Professional_2257Pro1stboot1.png)
then into shiny, [plasticky](https://www.reddit.com/media?url=https%3A%2F%2Fpreview.redd.it%2Fw038pvfct6hd1.png%3Fwidth%3D1024%26format%3Dpng%26auto%3Dwebp%26s%3Db9374cc14fa18c9bc7f4f57bd19c54869c2b81a7)
theme with soft edges that was again revised into [gloss](https://microsoft.fandom.com/wiki/Windows_XP_themes?file=Windows_XP_Royale.png), then glass.

I recall installing on my first computer a theme that sort of blended elements from this whole era into one, [Watercolor Emico
by Jamush](https://www.deviantart.com/jamush/art/Watercolor-Emico-Blue-87650169). It was this theme that stuck out most in my head
as I began to develop this website's new style.

## Goals

I didn't set out any formal requirements, but a loose idea of what I wanted started to take shape early on.

* make unusual choices that stand out a bit from common patterns in modern computing
* this is **not** flat design: shadows, bevels, gradients, and shine all have a place when used tastefully
* adapt responsively to small and large screens
* remain lightweight
* use modern browser features
* be playful, but optimize for reading

## Color and oklch()

One thing that I didn't decide at an early stage was a color palatte. Typically, my posts are image-light,
so to have some form of visual interest I prefer to pick a pleasant color or two. The trouble is, I change
the palette once every couple of years on my site, and it never feels like I've made a final or correct
choice.

So what if I **didn't** choose a brand color?

In the last couple of years, a new color space [called OKLCH](https://evilmartians.com/chronicles/oklch-in-css-why-quit-rgb-hsl) was added to CSS.
This has several interesting qualities that I made use of:

1. Its components are expressed as lightness, chroma (~saturation), and hue
2. Different hues with the same lightness and chroma are visually 'as bright' and 'as saturated'.
3. Different lightnesses with the same chroma and hue are visually 'the same color' with more or less light.
4. Different chroma values with the same lightness and hue are visually 'as bright' and 'the same color', just more or less gray.
5. You can perform expressions on these variables directly in CSS.

So I elected to create a set of rules that takes any input color and uses them to determine an appropriate
color palette for the page (and added a per-page primary color variable to my site template).

## What we left behind