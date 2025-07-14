+++
title = "Design Principles for my Website"
date = 2025-07-14
updated = 2025-07-14
[extra]
color = "#f2b762"
subtitle = "OR: codifying and post-hoc rationalizing the choices I've made thus far"
+++

I'm a tech guy, and as a general rule, tech guys don't make great UIs. Nonetheless,
I enjoy tinkering with websites and graphic design as a bit of a hobby. This site
has become something of a playground for me to exercise my CSS skills, and I've
decided it's time to put in writing some of the ideas that have been in my head
while I mess with it.

## A. Technical
1. Do much with little: a standard text-only blog page should require \<512KB of resources on first load.
2. Chained resource loading should be avoided if possible on initial page load. Allow required resources to be downloaded in parallel.
3. Avoid automatically loading resources which will not be used for this page if they are of substantial size.
4. The page should gracefully degrade for older/simpler browsers. Specifically, the inline stylesheet should allow for
a pleasant reading experience in even Firefox 3 and the structure should do the same for lynx.
5. Avoid loading images where possible.

## B. Aesthetic
1. Eschew corporate "flat design" and reflexive minimalism.
2. Borrow elements from different eras of computing history, including 1980s terminals and graphical desktop environments from 1993-2009.
   1. IBM AS/400
   2. *nixen
   3. MS-DOS and compatibles
   4. Windows 3, 95, 98, 2000, XP MCE, Longhorn, Vista
   5. OS/2 Warp
   6. Apple II series
   7. Mac OS System 7, 8, 9, Rhapsody, Jaguar, Panther, Tiger, Leopard
   8. Apple print advertisements
   9. KDE 3.5
   10. GNOME 2 from Ubuntu 8.04--10.10
   11. Winamp [Skins 🦙](https://skins.webamp.org/)
   12. Mac OS [Themes](https://macthemes.garden/)
   13. XP [Themes](https://www.deviantart.com/jamush/art/Watercolor-Emico-Black-89206761)
3. Gradients are underrated.
4. Not everything needs rounded corners.
5. Depth adds interest and helps the user feel grounded, if it is coherent.
6. Use relative colors, such that the page looks good with ANY selected primary color.
7. This isn't Neocities: keep styles coherent such that most of the site is obviously in the same design language.
8. This isn't Serious Business: a little bit of whimsy doesn't hurt

## C. Usability & Accessibility
1. Interactable elements should be obvious, and their function telegraphed by their appearance.
2. Headings and body text should be readable with a sane size and line width. The font size should follow the user's browser defaults.
3. Images need alt tags.
4. Semantic HTML elements are preferred to `<div>` and `<span>`
5. Consider screen sizes from 350x350 to 4K at standard DPI.
6. Contrast matters.
7. Neither light mode nor dark mode receives preferential treatment.