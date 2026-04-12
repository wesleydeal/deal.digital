package render

import (
	"strings"
	"testing"

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
