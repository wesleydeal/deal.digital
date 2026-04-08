package content

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type rawDocument struct {
	FrontMatter FrontMatter
	Body        string
	Relative    string
	SourcePath  string
}

func LoadSite(cfg Config, opts Options) (*Site, error) {
	pages := []*Page{}
	sections := []*Section{}
	pageByRoute := map[string]*Page{}
	sectionByRoute := map[string]*Section{}

	if err := filepath.WalkDir(opts.ContentDir, func(fullPath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if filepath.Ext(entry.Name()) != ".md" {
			return nil
		}

		rel, err := filepath.Rel(opts.ContentDir, fullPath)
		if err != nil {
			return err
		}
		doc, err := parseMarkdownDocument(fullPath, rel)
		if err != nil {
			return fmt.Errorf("%s: %w", rel, err)
		}
		if doc.FrontMatter.Draft && !opts.IncludeDrafts {
			return nil
		}

		if path.Base(filepathToSlash(rel)) == "_index.md" {
			route := routeForSection(rel)
			title := fallbackTitle(doc.FrontMatter.Title, rel)
			if route == "/" && doc.FrontMatter.Title == "" {
				title = cfg.Title
			}
			section := &Section{
				Title:             title,
				Slug:              slugify(doc.FrontMatter.Title),
				Route:             route,
				Permalink:         permalink(cfg.BaseURL, route),
				SourcePath:        fullPath,
				RelativePath:      filepathToSlash(rel),
				ParentSection:     parentSectionFor(rel),
				Template:          defaultSectionTemplate(doc.FrontMatter.Template, routeForSection(rel)),
				PageTemplate:      normalizePageTemplate(doc.FrontMatter.PageTemplate),
				InsertAnchorLinks: doc.FrontMatter.InsertAnchorLinks,
				SortBy:            doc.FrontMatter.SortBy,
				Extra:             cloneMap(doc.FrontMatter.Extra),
				Body:              doc.Body,
				Checksum:          checksum(filepathToSlash(rel), doc.Body, fmt.Sprintf("%v", doc.FrontMatter.Extra)),
			}
			sections = append(sections, section)
			sectionByRoute[section.Route] = section
			return nil
		}

		route := routeForPage(rel)
		bundleDir := path.Dir(filepathToSlash(rel))
		if path.Base(filepathToSlash(rel)) != "index.md" {
			bundleDir = path.Join(path.Dir(filepathToSlash(rel)), strings.TrimSuffix(path.Base(filepathToSlash(rel)), ".md"))
		}
		page := &Page{
			Title:         fallbackTitle(doc.FrontMatter.Title, rel),
			Slug:          slugify(titleOrFile(doc.FrontMatter.Title, rel)),
			Route:         route,
			Permalink:     permalink(cfg.BaseURL, route),
			SourcePath:    fullPath,
			RelativePath:  filepathToSlash(rel),
			BundleDir:     bundleDir,
			ParentSection: parentSectionFor(rel),
			Template:      doc.FrontMatter.Template,
			PageTemplate:  normalizePageTemplate(doc.FrontMatter.PageTemplate),
			Author:        doc.FrontMatter.Author,
			Authors:       normalizeAuthors(doc.FrontMatter),
			Date:          doc.FrontMatter.Date,
			Updated:       doc.FrontMatter.Updated,
			Draft:         doc.FrontMatter.Draft,
			Raw:           doc.FrontMatter.Raw,
			TOCEnabled:    doc.FrontMatter.TOC,
			Taxonomies:    cloneTaxonomies(doc.FrontMatter.Taxonomies),
			Extra:         cloneMap(doc.FrontMatter.Extra),
			Body:          doc.Body,
			Checksum:      checksum(filepathToSlash(rel), doc.Body, fmt.Sprintf("%v", doc.FrontMatter.Extra)),
		}
		pages = append(pages, page)
		pageByRoute[route] = page
		return nil
	}); err != nil {
		return nil, err
	}

	if _, ok := sectionByRoute["/"]; !ok {
		root := &Section{
			Title:         cfg.Title,
			Route:         "/",
			Permalink:     permalink(cfg.BaseURL, "/"),
			Template:      "home",
			RelativePath:  "_index.md",
			SourcePath:    filepath.Join(opts.ContentDir, "_index.md"),
			ParentSection: "/",
			Extra:         map[string]any{},
		}
		sections = append(sections, root)
		sectionByRoute[root.Route] = root
	}

	for _, section := range sections {
		section.ParentSection = nearestSectionRoute(sectionByRoute, section.ParentSection, section.Route)
	}

	for _, page := range pages {
		page.ParentSection = nearestSectionRoute(sectionByRoute, page.ParentSection, page.Route)
		parent := sectionByRoute[page.ParentSection]
		if parent == nil {
			parent = sectionByRoute["/"]
		}
		parent.Pages = append(parent.Pages, page)
		if page.PageTemplate == "" {
			page.PageTemplate = normalizePageTemplate(parent.PageTemplate)
		}
		if page.Template == "" {
			page.Template = DefaultPageTemplate(parent.Route)
		}
	}

	for _, section := range sections {
		if section.Route == "/" {
			continue
		}
		parent := sectionByRoute[section.ParentSection]
		if parent == nil || parent == section {
			parent = sectionByRoute["/"]
		}
		parent.Subsections = append(parent.Subsections, section)
		if section.PageTemplate == "" && parent.PageTemplate != "" {
			section.PageTemplate = normalizePageTemplate(parent.PageTemplate)
		}
	}

	for _, section := range sections {
		sortPages(section.Pages, firstNonEmpty(section.SortBy, "title"))
		slices.SortFunc(section.Subsections, func(a, b *Section) int {
			return strings.Compare(a.Route, b.Route)
		})
	}
	sortPages(pages, "date")

	site := &Site{
		Config:          cfg,
		Pages:           pages,
		Sections:        sections,
		PageByRoute:     pageByRoute,
		SectionByRoute:  sectionByRoute,
		Taxonomies:      buildTaxonomies(cfg, pages),
		RecentPages:     slices.Clone(pages),
		BuildTime:       time.Now().UTC(),
		SearchIndexName: "search_index." + firstNonEmpty(cfg.DefaultLanguage, "en") + ".json",
	}
	return site, nil
}

func parseMarkdownDocument(fullPath, relative string) (rawDocument, error) {
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return rawDocument{}, err
	}

	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	body := text
	front := FrontMatter{Taxonomies: map[string][]string{}, Extra: map[string]any{}}
	if strings.HasPrefix(text, "+++\n") {
		rest := text[len("+++\n"):]
		end := strings.Index(rest, "\n+++")
		if end < 0 {
			return rawDocument{}, fmt.Errorf("unterminated TOML frontmatter")
		}
		fmRaw := rest[:end]
		body = strings.TrimPrefix(rest[end+len("\n+++"):], "\n")
		fm, err := decodeFrontMatter(fmRaw)
		if err != nil {
			return rawDocument{}, err
		}
		front = fm
	}

	return rawDocument{
		FrontMatter: front,
		Body:        body,
		Relative:    filepathToSlash(relative),
		SourcePath:  fullPath,
	}, nil
}

func decodeFrontMatter(raw string) (FrontMatter, error) {
	var input map[string]any
	if _, err := toml.Decode(raw, &input); err != nil {
		return FrontMatter{}, err
	}

	fm := FrontMatter{
		Title:             lookupString(input, "title"),
		Date:              lookupTime(input, "date"),
		Updated:           lookupTime(input, "updated"),
		Draft:             lookupBool(input, "draft"),
		Template:          lookupString(input, "template"),
		PageTemplate:      lookupString(input, "page_template"),
		InsertAnchorLinks: lookupString(input, "insert_anchor_links"),
		SortBy:            lookupString(input, "sort_by"),
		Author:            lookupString(input, "author"),
		Authors:           lookupStringSlice(input, "authors"),
		Raw:               lookupBool(input, "raw"),
		TOC:               lookupBool(input, "toc"),
		Taxonomies:        map[string][]string{},
		Extra:             map[string]any{},
	}
	if fm.Author != "" && len(fm.Authors) == 0 {
		fm.Authors = []string{fm.Author}
	}
	if nested, ok := input["taxonomies"].(map[string]any); ok {
		for key, value := range nested {
			fm.Taxonomies[key] = anySliceToStrings(value)
		}
	}
	if nested, ok := input["extra"].(map[string]any); ok {
		for key, value := range nested {
			fm.Extra[key] = value
		}
	}

	known := map[string]struct{}{
		"title": {}, "date": {}, "updated": {}, "draft": {}, "template": {}, "page_template": {},
		"insert_anchor_links": {}, "sort_by": {}, "author": {}, "authors": {}, "taxonomies": {},
		"extra": {}, "raw": {}, "toc": {},
	}
	for key, value := range input {
		if _, ok := known[key]; ok {
			continue
		}
		fm.Extra[key] = value
	}
	return fm, nil
}

func lookupString(input map[string]any, key string) string {
	if s, ok := stringValue(input[key]); ok {
		return s
	}
	return ""
}

func stringValue(v any) (string, bool) {
	switch value := v.(type) {
	case string:
		return value, true
	case fmt.Stringer:
		return value.String(), true
	default:
		return "", false
	}
}

func StringValue(v any) (string, bool) {
	return stringValue(v)
}

func lookupBool(input map[string]any, key string) bool {
	value, ok := input[key]
	if !ok {
		return false
	}
	b, _ := value.(bool)
	return b
}

func lookupStringSlice(input map[string]any, key string) []string {
	value, ok := input[key]
	if !ok {
		return nil
	}
	return anySliceToStrings(value)
}

func anySliceToStrings(value any) []string {
	switch items := value.(type) {
	case []string:
		return slices.Clone(items)
	case []any:
		out := make([]string, 0, len(items))
		for _, item := range items {
			if s, ok := stringValue(item); ok && s != "" {
				out = append(out, s)
			}
		}
		return out
	default:
		return nil
	}
}

func lookupTime(input map[string]any, key string) *time.Time {
	value, ok := input[key]
	if !ok {
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		t := v.UTC()
		return &t
	case string:
		for _, format := range []string{time.DateOnly, time.RFC3339, "2006-01-02 15:04:05", "2006-01-02T15:04:05"} {
			if parsed, err := time.Parse(format, v); err == nil {
				t := parsed.UTC()
				return &t
			}
		}
	default:
		text := strings.TrimSpace(fmt.Sprint(v))
		for _, format := range []string{time.DateOnly, time.RFC3339, "2006-01-02 15:04:05", "2006-01-02T15:04:05"} {
			if parsed, err := time.Parse(format, text); err == nil {
				t := parsed.UTC()
				return &t
			}
		}
	}
	return nil
}

func normalizeAuthors(fm FrontMatter) []string {
	if len(fm.Authors) > 0 {
		return slices.Clone(fm.Authors)
	}
	if fm.Author != "" {
		return []string{fm.Author}
	}
	return nil
}

func buildTaxonomies(cfg Config, pages []*Page) map[string]map[string]*TaxonomyTerm {
	terms := map[string]map[string]*TaxonomyTerm{}
	configByName := map[string]TaxonomyConfig{}
	for _, taxonomy := range cfg.Taxonomies {
		configByName[taxonomy.Name] = taxonomy
		terms[taxonomy.Name] = map[string]*TaxonomyTerm{}
	}
	for _, page := range pages {
		pageTerms := cloneTaxonomies(page.Taxonomies)
		if len(page.Authors) > 0 {
			pageTerms["author"] = append(pageTerms["author"], page.Authors...)
		}
		for name, values := range pageTerms {
			if _, ok := terms[name]; !ok {
				terms[name] = map[string]*TaxonomyTerm{}
			}
			for _, value := range values {
				if value == "" {
					continue
				}
				slug := slugify(value)
				if terms[name][slug] == nil {
					terms[name][slug] = &TaxonomyTerm{
						Name:   value,
						Slug:   slug,
						Route:  "/" + slugify(name) + "/" + slug + "/",
						Config: configByName[name],
					}
				}
				terms[name][slug].Pages = append(terms[name][slug].Pages, page)
			}
		}
	}
	for _, taxTerms := range terms {
		for _, term := range taxTerms {
			sortPages(term.Pages, "date")
		}
	}
	return terms
}

func cloneMap(in map[string]any) map[string]any {
	if len(in) == 0 {
		return map[string]any{}
	}
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func cloneTaxonomies(in map[string][]string) map[string][]string {
	if len(in) == 0 {
		return map[string][]string{}
	}
	out := make(map[string][]string, len(in))
	for k, v := range in {
		out[k] = slices.Clone(v)
	}
	return out
}

func fallbackTitle(title, rel string) string {
	if title != "" {
		return title
	}
	return strings.TrimSuffix(path.Base(filepathToSlash(rel)), path.Ext(rel))
}

func titleOrFile(title, rel string) string {
	if title != "" {
		return title
	}
	name := strings.TrimSuffix(path.Base(filepathToSlash(rel)), path.Ext(rel))
	if name == "index" {
		return path.Base(path.Dir(filepathToSlash(rel)))
	}
	return name
}

func defaultSectionTemplate(templateName, route string) string {
	if templateName != "" {
		switch templateName {
		case "index.html":
			return "home"
		case "blog.html":
			return "blog_list"
		case "list.html":
			return "list"
		default:
			return strings.TrimSuffix(templateName, ".html")
		}
	}
	if route == "/" {
		return "home"
	}
	return "list"
}

func DefaultPageTemplate(route string) string {
	if strings.HasPrefix(route, "/blog/") {
		return "blog_page"
	}
	return "page"
}

func nearestSectionRoute(sectionByRoute map[string]*Section, candidate string, selfRoute string) string {
	candidate = strings.TrimSpace(candidate)
	if candidate == "" {
		candidate = "/"
	}
	if candidate == selfRoute {
		candidate = parentRoute(candidate)
	}

	for candidate != "" {
		if candidate != selfRoute {
			if _, ok := sectionByRoute[candidate]; ok {
				return candidate
			}
		}
		if candidate == "/" {
			return "/"
		}
		trimmed := strings.Trim(candidate, "/")
		if trimmed == "" {
			return "/"
		}
		parts := strings.Split(trimmed, "/")
		if len(parts) <= 1 {
			candidate = "/"
			continue
		}
		candidate = "/" + strings.Join(parts[:len(parts)-1], "/") + "/"
	}

	return "/"
}

func parentRoute(route string) string {
	if route == "/" {
		return "/"
	}
	trimmed := strings.Trim(route, "/")
	if trimmed == "" {
		return "/"
	}
	parts := strings.Split(trimmed, "/")
	if len(parts) <= 1 {
		return "/"
	}
	return "/" + strings.Join(parts[:len(parts)-1], "/") + "/"
}

func normalizePageTemplate(name string) string {
	switch strings.TrimSpace(name) {
	case "", "page", "page.html":
		return "page"
	case "blog_page", "blog-page", "blog-page.html":
		return "blog_page"
	case "raw", "raw.html":
		return "raw"
	default:
		return strings.TrimSuffix(name, ".html")
	}
}

func firstNonEmpty(items ...string) string {
	for _, item := range items {
		if strings.TrimSpace(item) != "" {
			return item
		}
	}
	return ""
}

func CollapseWhitespace(v string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(v)), " ")
}

func prependIfMissing(prefix string, body string) string {
	if strings.HasPrefix(body, prefix) {
		return body
	}
	var buf bytes.Buffer
	buf.WriteString(prefix)
	buf.WriteString(body)
	return buf.String()
}
