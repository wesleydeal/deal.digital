package render

import (
	"bytes"
	"encoding/json"
	"html"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"deal.digital/internal/content"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	renderhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

type Renderer struct {
	md          goldmark.Markdown
	templateDir string
	templates   *template.Template
}

type templateData struct {
	Site         *content.Site
	Page         *content.Page
	Section      *content.Section
	Term         *content.TaxonomyTerm
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

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"html":           func(v string) template.HTML { return template.HTML(v) },
		"inlineMarkdown": inlineMarkdown,
		"date": func(v *time.Time) string {
			if v == nil {
				return ""
			}
			return v.Format("2006-01-02")
		},
		"pageTitle": content.DisplayTitle,
		"extraString": func(m map[string]any, key string) string {
			if s, ok := content.StringValue(m[key]); ok {
				return s
			}
			return ""
		},
		"extraBool":       extraBool,
		"colorFor":        colorFor,
		"documentTitle":   documentTitle,
		"bodyID":          bodyID,
		"showBreadcrumbs": showBreadcrumbs,
		"primaryAuthor":   primaryAuthor,
		"breadcrumbsHTML": breadcrumbsHTML,
		"siteTreeHTML":    siteTreeHTML,
		"sourceURL":       sourceURL,
	}
}

func NewRenderer(templateDir string) (*Renderer, error) {
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

	r := &Renderer{
		md:          md,
		templateDir: templateDir,
	}
	if err := r.loadTemplates(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Renderer) RenderSite(site *content.Site) error {
	if err := r.loadTemplates(); err != nil {
		return err
	}
	for _, section := range site.Sections {
		htmlBody, _, _, err := r.renderMarkdown(section.Body, 0)
		if err != nil {
			return err
		}
		section.HTML = template.HTML(htmlBody)
	}

	for _, page := range site.Pages {
		var htmlBody, plain string
		var headings []content.Heading
		var err error

		if strings.EqualFold(strings.TrimSuffix(page.Template, ".html"), "raw") || page.Raw {
			htmlBody = r.preprocessShortcodes(page.Body, 0)
			plain = content.CollapseWhitespace(stripTags(htmlBody))
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

func (r *Renderer) loadTemplates() error {
	tmpl := template.New("site").Funcs(templateFuncs())
	parsed, err := tmpl.ParseGlob(filepath.Join(r.templateDir, "*.tmpl"))
	if err != nil {
		return err
	}
	r.templates = parsed
	return nil
}
func (r *Renderer) renderMarkdown(input string, depth int) (string, []content.Heading, string, error) {
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

	out := buf.String()
	out = addHeadingAnchors(out)
	headings := extractHeadingsFromHTML(out)
	plain := content.CollapseWhitespace(stripTags(out))
	return out, headings, plain, nil
}

func addHeadingAnchors(input string) string {
	return headingHTMLR.ReplaceAllString(input, `<h$1 id="$2"><a class="anchor" href="#$2" aria-label="Anchor link for: $2">#</a>$3</h$1>`)
}

func extractHeadingsFromHTML(input string) []content.Heading {
	headings := []content.Heading{}
	for _, match := range headingHTMLR.FindAllStringSubmatch(input, -1) {
		level, err := strconv.Atoi(match[1])
		if err != nil || level != 2 {
			continue
		}
		title := content.CollapseWhitespace(stripTags(match[3]))
		if title == "" {
			continue
		}
		title = strings.TrimPrefix(title, "# ")
		headings = append(headings, content.Heading{
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

func (r *Renderer) RenderPageDocument(site *content.Site, page *content.Page) ([]byte, error) {
	if strings.EqualFold(strings.TrimSuffix(page.Template, ".html"), "raw") {
		return trimDocument([]byte(page.HTML)), nil
	}
	name := page.PageTemplate
	if name == "" {
		name = content.DefaultPageTemplate(page.Route)
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
	return trimDocument(buf.Bytes()), nil
}

func (r *Renderer) RenderSectionDocument(site *content.Site, section *content.Section) ([]byte, error) {
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
	return trimDocument(buf.Bytes()), nil
}

func (r *Renderer) RenderTaxonomyDocument(site *content.Site, _ string, term *content.TaxonomyTerm) ([]byte, error) {
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
	return trimDocument(buf.Bytes()), nil
}

func (r *Renderer) RenderNotFoundDocument(site *content.Site) ([]byte, error) {
	data := templateData{
		Site:         site,
		Title:        "404 Not Found",
		CurrentRoute: "/404.html",
		BuildTime:    site.BuildTime,
	}
	var buf bytes.Buffer
	if err := r.templates.ExecuteTemplate(&buf, "404", data); err != nil {
		return nil, err
	}
	return trimDocument(buf.Bytes()), nil
}

func WriteJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func trimDocument(data []byte) []byte {
	return bytes.TrimSpace(data)
}

func taxonomyTermsSorted(terms map[string]*content.TaxonomyTerm) []*content.TaxonomyTerm {
	out := make([]*content.TaxonomyTerm, 0, len(terms))
	for _, term := range terms {
		out = append(out, term)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
