+++
title="Designing for print & beyond 	with HTML/CSS: building my web-native resume"
date=2024-04-28
updated=2024-04-28
draft=true
+++

TL;DR: Building a responsive [resume](/resume) taught me that modern web browsers are tremendously flexible without
need for frontend or backend frameworks and can be adapted quite well to typesetting for print, with some caveats.

## The Problem

What is the best way to create a document that:
* can easily be viewed online or in print,
* can be read on a smartphone,
* looks consistently good across platforms,
* allows easy change tracking,
* and can be tweaked on the fly on almost any platform?

For a while, my answer to this question was Google Docs. It runs in a browser on my phone or laptop,
it can export to PDF and to print quite easily, and it can even responsively shrink your page for
reading and writing on mobile.

There are a couple of main problems with this approach for a resume. For one, Google Docs is not
built for responsive design. You can't set breakpoints that change the layout and style of parts
of the page on smaller devices, so you're stuck trying to come up with something that looks good
enough on a small screen without negatively impacting your full-size layout.

Second, Google Docs styles don't cascade. That means that changes to the the or paragraph spacing tend to involve
lots of clicks around the document to ensure that properties are all set the same way.

Third, Google Docs is limited in style and interaction capabilities. I found that ATS scanning my
resume would almost universally fail to handle ligatures like the connected "*ffi*" in *efficient*.
But there's no way to turn these off other than to pick a font that doesn't support them. Or let's
say I would like to add a quick button to change the font to a sans-serif. Though achievable with
custom scripting, it's clear this type of interactivity is not what Docs is built for.

Thankfully, there is a stack built for responsive design, cascading style sheets, and advanced
styling and scripting. And it's already installed on all your devices.

## A Blank Slate

We'll start with the minimum reasonable HTML/CSS file:

```html
<!doctype html>
<html lang="en">
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>My Resume</title>
	<style>
	</style>
	<body>
		<main id="resume">
		<!-- we'll be adding stuff in here next -->
		</main>
	</body>
</html>
```

We'll structure our resume with semantic elements. Our resume is an article with a header, sections, and a footer.
The sections consist of a top-level heading and then may contain either subheadings or content. Our job titles are
subheadings with unordered lists of accomplishments. We'll use `<summary>` and `<details>` to make them collapsible.
Our education & certifications section has subheadings (we
will want them a bit smaller, so we use a different tag) and optionally some associated details. Our skills section
is just a few paragraphs beginning with a bolded category and a comma-separated list of skills. In the heading,
we want our name and job title in a seperate area from our contact information, so we make them separate `<div>`s.

```html
<header>
	<div id="namebox">
		<p class="name">John Doe
		<p class="jobtitle">Widget Developer
	</div>
	<div id="contact">
		<p><a href="tel:555-555-5555">(555) 555-5555</a>
		<p><a href="mailto:john.doe@example.com">john.doe@example.com</a>
		<p>Anchorage, Alaska
	</div>
</header>
<section id="experience">
	<h1>Work Experience</h1>
	<details open>
		<summary>
			<h2><b>Widget Developer</b> at NuWidget</h2>
			<span class="dates">May 2017 &ndash; Present</span>
		</summary>
		<ul>
			<li>Spearheaded the development of a cutting-edge widget customization tool, resulting in a 30% increase in user engagement within the first quarter of release.</li>
			<li>Collaborated with cross-functional teams to streamline widget integration processes, reducing deployment time by 40%.</li>
			<li>Implemented innovative widget performance optimization techniques, enhancing load times by 50% and ensuring seamless user experiences across all platforms.</li>
			<li>Led a team of developers in creating a library of reusable widget components, improving development efficiency by 25% and fostering consistency in design and functionality.</li>
			<li>Conducted comprehensive research on emerging widget technologies and industry trends, driving continuous improvement and maintaining NuWidget's competitive edge in the market.</li>
		</ul>
	</details>
	<details open>
		<summary>
			<h2><b>Junior Widget Developer</b> at Acme Widgets Inc.</h2>
			<span class="dates">December 2015 &ndash; July 2017</span>
		</summary>
		<ul>
			<li>Orchestrated the migration of legacy widget systems to a modern architecture, resulting in a 60% reduction in maintenance costs and improved scalability.</li>
			<li>Championed accessibility initiatives by integrating ADA-compliant features into widget designs, expanding user base and ensuring inclusivity for all customers.</li>
		</ul>
	</details>
</section>
<section id="education-certs">
	<h1>Education & Certifications</h1>
	<h3><span><b>Bachelor of Science in Computer Science</b>, NuTech University</span> <span class="dates">2016</span></h3>
	<p>Relevant Coursework: Advanced Software Engineering, Web Development, Data Structures and Algorithms
</section>
<section id="skills">
	<h1>Skills</h1>
	<p><strong>General</strong>: Problem-solving, Communication
	<p><strong>Web Dev</strong>: JavaScript, HTML, CSS, Python, React, Vue.js, Git, Agile methodologies, Jira, Responsive design, Cross-browser compatibility, Problem-solving, Communication
</section>
<footer>
	<p>Written and typeset in HTML/CSS. Updated <time datetime="2024-04-28">2024-04-28</time>. References available upon request.
</footer>
```

Alright, so we have a structure for our resume. Now let's fill in our `<style>`sheet.

```css
body
```