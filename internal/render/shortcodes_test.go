package render

import (
	"strings"
	"testing"

	"deal.digital/internal/content"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	renderhtml "github.com/yuin/goldmark/renderer/html"
)

func newTestRenderer() *Renderer {
	return &Renderer{
		md: goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				extension.Footnote,
				extension.Linkify,
				extension.Typographer,
				highlighting.NewHighlighting(
					highlighting.WithFormatOptions(
						chromahtml.WithClasses(true),
					),
				),
			),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
			goldmark.WithRendererOptions(
				renderhtml.WithUnsafe(),
			),
		),
	}
}

func TestRenderMarkdownRawBlockPreservesHTML(t *testing.T) {
	r := newTestRenderer()

	input := strings.TrimSpace(`
Before.

{% raw %}
<div class="demo">
	<span>**not markdown**</span>
</div>
{% endraw %}

After.
`)

	got, _, _, err := r.renderMarkdown(input, 0)
	if err != nil {
		t.Fatalf("render markdown: %v", err)
	}

	if !strings.Contains(got, "<p>Before.</p>") {
		t.Fatalf("missing markdown paragraph before raw block: %s", got)
	}
	if !strings.Contains(got, `<div class="demo">`) {
		t.Fatalf("missing raw html block: %s", got)
	}
	if !strings.Contains(got, `<span>**not markdown**</span>`) {
		t.Fatalf("raw block body was processed as markdown: %s", got)
	}
	if strings.Contains(got, "<p><div") {
		t.Fatalf("raw block was incorrectly wrapped in a paragraph: %s", got)
	}
	if !strings.Contains(got, "<p>After.</p>") {
		t.Fatalf("missing markdown paragraph after raw block: %s", got)
	}
}

func TestRenderMarkdownRawBlockSupportsLegacySyntax(t *testing.T) {
	r := newTestRenderer()

	input := strings.TrimSpace(`
{% raw () %}
<section data-kind="legacy">hello</section>
{% end %}
`)

	got, _, _, err := r.renderMarkdown(input, 0)
	if err != nil {
		t.Fatalf("render markdown: %v", err)
	}

	if !strings.Contains(got, `<section data-kind="legacy">hello</section>`) {
		t.Fatalf("legacy raw block did not render: %s", got)
	}
	if strings.Contains(got, "&lt;section") {
		t.Fatalf("legacy raw block was escaped: %s", got)
	}
}

func TestRenderMarkdownVariants(t *testing.T) {
	r := newTestRenderer()
	variants := []content.PageVariant{
		{ID: "general", Label: "General IT", Route: "/resume/"},
		{ID: "neteng", Label: "Network Engineering", Route: "/resume/neteng/"},
	}
	input := strings.TrimSpace(`
Shared {% variant(include="neteng") %}network{% endvariant %} experience.

{% variant(include="neteng") %}
* Network-only list item
{% endvariant %}

{% variant(exclude="neteng") %}General-only text.{% endvariant %}

{{ variant_menu() }}
`)

	general, _, _, err := r.renderMarkdownWithVariant(input, 0, &variantContext{current: "general", variants: variants})
	if err != nil {
		t.Fatalf("render general variant: %v", err)
	}
	if !strings.Contains(general, "Shared  experience.") || !strings.Contains(general, "General-only text.") {
		t.Fatalf("general content missing: %s", general)
	}
	if strings.Contains(general, "Network-only") || strings.Contains(general, "network experience") {
		t.Fatalf("network content leaked into general variant: %s", general)
	}
	if !strings.Contains(general, `<menu class="page-variants">`) || !strings.Contains(general, `<span aria-current="page">General IT</span>`) || !strings.Contains(general, `<a href="/resume/neteng/">Network Engineering</a>`) {
		t.Fatalf("general variant menu is incorrect: %s", general)
	}

	neteng, _, _, err := r.renderMarkdownWithVariant(input, 0, &variantContext{current: "neteng", variants: variants})
	if err != nil {
		t.Fatalf("render network variant: %v", err)
	}
	if !strings.Contains(neteng, "network experience") || !strings.Contains(neteng, "Network-only list item") {
		t.Fatalf("network content missing: %s", neteng)
	}
	if strings.Contains(neteng, "General-only") {
		t.Fatalf("general-only content leaked into network variant: %s", neteng)
	}
	if !strings.Contains(neteng, `<a href="/resume/">General IT</a>`) || !strings.Contains(neteng, `<span aria-current="page">Network Engineering</span>`) {
		t.Fatalf("network variant menu is incorrect: %s", neteng)
	}
}

func TestRenderMarkdownVariantsRejectInvalidUsage(t *testing.T) {
	r := newTestRenderer()
	context := &variantContext{current: "general", variants: []content.PageVariant{{ID: "general", Label: "General", Route: "/"}}}
	for _, input := range []string{
		`{% variant(include="missing") %}Nope{% endvariant %}`,
		`{% variant(include="general") %}Unclosed`,
		`{{ variant_menu() }}`,
	} {
		variant := context
		if input == `{{ variant_menu() }}` {
			variant = nil
		}
		if _, _, _, err := r.renderMarkdownWithVariant(input, 0, variant); err == nil {
			t.Fatalf("invalid variant shortcode was accepted: %s", input)
		}
	}
}
