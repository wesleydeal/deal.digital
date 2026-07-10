package content

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"path"
	"sort"
	"strings"
	"time"
)

type FrontMatter struct {
	Title             string              `toml:"title"`
	Date              *time.Time          `toml:"date"`
	Updated           *time.Time          `toml:"updated"`
	Draft             bool                `toml:"draft"`
	Template          string              `toml:"template"`
	PageTemplate      string              `toml:"page_template"`
	InsertAnchorLinks string              `toml:"insert_anchor_links"`
	SortBy            string              `toml:"sort_by"`
	Author            string              `toml:"author"`
	Authors           []string            `toml:"authors"`
	Taxonomies        map[string][]string `toml:"taxonomies"`
	Variants          *VariantConfig      `toml:"variants"`
	Extra             map[string]any      `toml:"extra"`
	Raw               bool                `toml:"raw"`
	TOC               bool                `toml:"toc"`
}

type VariantConfig struct {
	Default string              `toml:"default"`
	Items   []VariantDefinition `toml:"items"`
}

type VariantDefinition struct {
	ID    string `toml:"id"`
	Label string `toml:"label"`
}

type PageVariant struct {
	ID    string
	Label string
	Route string
}

type Page struct {
	Title          string
	Slug           string
	Route          string
	Permalink      string
	SourcePath     string
	RelativePath   string
	BundleDir      string
	ParentSection  string
	Template       string
	PageTemplate   string
	Author         string
	Authors        []string
	Date           *time.Time
	Updated        *time.Time
	Draft          bool
	Raw            bool
	TOCEnabled     bool
	Variants       []PageVariant
	DefaultVariant string
	Taxonomies     map[string][]string
	Extra          map[string]any
	Body           string
	HTML           template.HTML
	Plain          string
	Headings       []Heading
	Checksum       string
}

type Section struct {
	Title             string
	Slug              string
	Route             string
	Permalink         string
	SourcePath        string
	RelativePath      string
	ParentSection     string
	Template          string
	PageTemplate      string
	InsertAnchorLinks string
	SortBy            string
	Extra             map[string]any
	Body              string
	HTML              template.HTML
	Pages             []*Page
	Subsections       []*Section
	Checksum          string
}

type TaxonomyTerm struct {
	Name   string
	Slug   string
	Route  string
	Pages  []*Page
	Config TaxonomyConfig
}

type Site struct {
	Config          Config
	Pages           []*Page
	Sections        []*Section
	PageByRoute     map[string]*Page
	SectionByRoute  map[string]*Section
	Taxonomies      map[string]map[string]*TaxonomyTerm
	RecentPages     []*Page
	BuildTime       time.Time
	SearchIndexName string
}

type Heading struct {
	Level int
	ID    string
	Title string
}

type searchDocument struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Body        string `json:"body"`
	Description string `json:"description,omitempty"`
}

func permalink(baseURL, route string) string {
	return strings.TrimRight(baseURL, "/") + route
}

func routeForPage(rel string) string {
	dir, file := path.Split(filepathToSlash(rel))
	if file == "index.md" {
		return "/" + strings.Trim(strings.TrimSuffix(dir, "/"), "/") + "/"
	}
	name := strings.TrimSuffix(file, path.Ext(file))
	full := path.Join(strings.Trim(dir, "/"), name)
	if full == "" {
		return "/"
	}
	return "/" + strings.Trim(full, "/") + "/"
}

func routeForSection(rel string) string {
	rel = filepathToSlash(rel)
	if rel == "_index.md" {
		return "/"
	}
	dir := path.Dir(rel)
	if dir == "." {
		return "/"
	}
	return "/" + strings.Trim(dir, "/") + "/"
}

func parentSectionFor(rel string) string {
	rel = filepathToSlash(rel)
	dir := path.Dir(rel)
	if path.Base(rel) == "_index.md" {
		dir = path.Dir(dir)
	}
	if dir == "." || dir == "/" {
		return "/"
	}
	return "/" + strings.Trim(dir, "/") + "/"
}

func filepathToSlash(v string) string {
	return strings.ReplaceAll(v, "\\", "/")
}

func checksum(parts ...string) string {
	h := sha256.New()
	for _, part := range parts {
		_, _ = h.Write([]byte(part))
		_, _ = h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))
}

func slugify(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	var b strings.Builder
	lastDash := false
	for _, r := range v {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		default:
			if !lastDash {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "item"
	}
	return out
}

func DisplayTitle(p *Page) string {
	if short, ok := stringValue(p.Extra["shorttitle"]); ok && short != "" {
		return short
	}
	return p.Title
}

func sortPages(pages []*Page, sortBy string) {
	sort.SliceStable(pages, func(i, j int) bool {
		a, b := pages[i], pages[j]
		switch sortBy {
		case "title":
			return strings.ToLower(DisplayTitle(a)) < strings.ToLower(DisplayTitle(b))
		case "date":
			if a.Date == nil && b.Date == nil {
				return a.Route < b.Route
			}
			if a.Date == nil {
				return false
			}
			if b.Date == nil {
				return true
			}
			if a.Date.Equal(*b.Date) {
				return a.Route < b.Route
			}
			return a.Date.After(*b.Date)
		default:
			return a.Route < b.Route
		}
	})
}

func PageDateString(page *Page) string {
	switch {
	case page.Date != nil:
		return page.Date.Format("2006-01-02")
	case page.Updated != nil:
		return page.Updated.Format("2006-01-02")
	default:
		return ""
	}
}

func requireString(m map[string]any, key string, fallback string) string {
	if s, ok := stringValue(m[key]); ok && s != "" {
		return s
	}
	return fallback
}

func SearchDocuments(site *Site) []searchDocument {
	docs := make([]searchDocument, 0, len(site.Pages))
	for _, page := range site.Pages {
		docs = append(docs, searchDocument{
			Title:       DisplayTitle(page),
			URL:         page.Route,
			Body:        CollapseWhitespace(page.Plain),
			Description: requireString(page.Extra, "subtitle", ""),
		})
	}
	sort.SliceStable(docs, func(i, j int) bool { return docs[i].URL < docs[j].URL })
	return docs
}

func (s *Site) String() string {
	return fmt.Sprintf("pages=%d sections=%d taxonomies=%d", len(s.Pages), len(s.Sections), len(s.Taxonomies))
}
