package site

import (
	"bytes"
	"embed"
	"encoding/json"
	"html"
	"html/template"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	renderhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

type Renderer struct {
	md        goldmark.Markdown
	templates *template.Template
}

type templateData struct {
	Site         *Site
	Page         *Page
	Section      *Section
	Term         *TaxonomyTerm
	Title        string
	CurrentRoute string
	BuildTime    time.Time
}

var (
	tagRE        = regexp.MustCompile(`<[^>]+>`)
	headingHTMLR = regexp.MustCompile(`(?s)<h([2-6]) id="([^"]+)">(.*?)</h[2-6]>`)
)

var (
	inlineMarkdownOnce sync.Once
	inlineMarkdownMD   goldmark.Markdown
)

func NewRenderer() (*Renderer, error) {
	funcs := template.FuncMap{
		"html":           func(v string) template.HTML { return template.HTML(v) },
		"inlineMarkdown": inlineMarkdown,
		"date": func(v *time.Time) string {
			if v == nil {
				return ""
			}
			return v.Format("2006-01-02")
		},
		"pageTitle": displayTitle,
		"extraString": func(m map[string]any, key string) string {
			if s, ok := stringValue(m[key]); ok {
				return s
			}
			return ""
		},
		"extraBool":       extraBool,
		"colorFor":        colorFor,
		"documentTitle":   documentTitle,
		"bodyID":          bodyID,
		"showBreadcrumbs": showBreadcrumbs,
		"breadcrumbsHTML": breadcrumbsHTML,
		"siteTreeHTML":    siteTreeHTML,
		"sourceURL":       sourceURL,
		"shouldRenderTaxonomy": func(site *Site, name string) bool {
			for _, tax := range site.Config.Taxonomies {
				if tax.Name == name {
					return tax.ShouldRender()
				}
			}
			return true
		},
	}

	tmpl, err := template.New("site").Funcs(funcs).ParseFS(templateFS, "templates/*.tmpl")
	if err != nil {
		return nil, err
	}

	md := goldmark.New(
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
	)

	return &Renderer{md: md, templates: tmpl}, nil
}

func (r *Renderer) RenderSite(site *Site) error {
	for _, section := range site.Sections {
		htmlBody, _, _, err := r.renderMarkdown(section.Body, 0)
		if err != nil {
			return err
		}
		section.HTML = template.HTML(htmlBody)
	}

	for _, page := range site.Pages {
		var htmlBody, plain string
		var headings []Heading
		var err error

		if strings.EqualFold(strings.TrimSuffix(page.Template, ".html"), "raw") {
			htmlBody = r.preprocessShortcodes(page.Body, 0)
			plain = collapseWhitespace(stripTags(htmlBody))
		} else if page.Raw {
			htmlBody = r.preprocessShortcodes(page.Body, 0)
			plain = collapseWhitespace(stripTags(htmlBody))
		} else {
			htmlBody, headings, plain, err = r.renderMarkdown(page.Body, 0)
			if err != nil {
				return err
			}
		}

		page.HTML = template.HTML(htmlBody)
		page.Plain = plain
		page.Headings = headings
	}
	return nil
}

func (r *Renderer) renderMarkdown(input string, depth int) (string, []Heading, string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", nil, "", nil
	}
	processed := r.preprocessShortcodes(input, depth)
	source := []byte(processed)
	reader := text.NewReader(source)
	doc := r.md.Parser().Parse(reader)

	var buf bytes.Buffer
	if err := r.md.Renderer().Render(&buf, source, doc); err != nil {
		return "", nil, "", err
	}

	html := buf.String()
	html = addHeadingAnchors(html)
	headings := extractHeadingsFromHTML(html)
	plain := collapseWhitespace(stripTags(html))
	return html, headings, plain, nil
}

func addHeadingAnchors(input string) string {
	return headingHTMLR.ReplaceAllString(input, `<h$1 id="$2"><a class="anchor" href="#$2" aria-label="Anchor link for: $2">#</a>$3</h$1>`)
}

func extractHeadingsFromHTML(input string) []Heading {
	headings := []Heading{}
	for _, match := range headingHTMLR.FindAllStringSubmatch(input, -1) {
		level, err := strconv.Atoi(match[1])
		if err != nil || level != 2 {
			continue
		}
		title := collapseWhitespace(stripTags(match[3]))
		if title == "" {
			continue
		}
		title = strings.TrimPrefix(title, "# ")
		headings = append(headings, Heading{
			Level: level,
			ID:    match[2],
			Title: title,
		})
	}
	return headings
}

func stripTags(v string) string {
	return html.UnescapeString(tagRE.ReplaceAllString(v, " "))
}

func inlineMarkdown(input string) template.HTML {
	if strings.TrimSpace(input) == "" {
		return ""
	}
	inlineMarkdownOnce.Do(func() {
		inlineMarkdownMD = goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				extension.Footnote,
				extension.Linkify,
				extension.Typographer,
			),
			goldmark.WithRendererOptions(
				renderhtml.WithUnsafe(),
			),
		)
	})

	var buf bytes.Buffer
	source := []byte(strings.TrimSpace(input))
	if err := inlineMarkdownMD.Convert(source, &buf); err != nil {
		return template.HTML(html.EscapeString(input))
	}

	out := strings.TrimSpace(buf.String())
	if strings.HasPrefix(out, "<p>") && strings.HasSuffix(out, "</p>") {
		out = strings.TrimPrefix(out, "<p>")
		out = strings.TrimSuffix(out, "</p>")
	}
	return template.HTML(out)
}

func (r *Renderer) RenderPageDocument(site *Site, page *Page) ([]byte, error) {
	if strings.EqualFold(strings.TrimSuffix(page.Template, ".html"), "raw") {
		return []byte(string(page.HTML)), nil
	}
	name := page.PageTemplate
	if name == "" {
		name = defaultPageTemplate(page.Route)
	}
	data := templateData{
		Site:         site,
		Page:         page,
		Title:        page.Title,
		CurrentRoute: page.Route,
		BuildTime:    site.BuildTime,
	}
	var buf bytes.Buffer
	if err := r.templates.ExecuteTemplate(&buf, name, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (r *Renderer) RenderSectionDocument(site *Site, section *Section) ([]byte, error) {
	name := section.Template
	if name == "" {
		name = "list"
	}
	data := templateData{
		Site:         site,
		Section:      section,
		Title:        section.Title,
		CurrentRoute: section.Route,
		BuildTime:    site.BuildTime,
	}
	var buf bytes.Buffer
	if err := r.templates.ExecuteTemplate(&buf, name, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (r *Renderer) RenderTaxonomyDocument(site *Site, name string, term *TaxonomyTerm) ([]byte, error) {
	data := templateData{
		Site:         site,
		Term:         term,
		Title:        term.Name,
		CurrentRoute: term.Route,
		BuildTime:    site.BuildTime,
	}
	var buf bytes.Buffer
	if err := r.templates.ExecuteTemplate(&buf, "taxonomy", data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func taxonomyTermsSorted(terms map[string]*TaxonomyTerm) []*TaxonomyTerm {
	out := make([]*TaxonomyTerm, 0, len(terms))
	for _, term := range terms {
		out = append(out, term)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
