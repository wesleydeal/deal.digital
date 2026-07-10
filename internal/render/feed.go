package render

import (
	"encoding/xml"
	"strings"
	"time"

	"deal.digital/internal/content"
)

const atomTimeFormat = "2006-01-02T15:04:05+00:00"

type atomFeed struct {
	XMLName   xml.Name      `xml:"feed"`
	Namespace string        `xml:"xmlns,attr"`
	Language  string        `xml:"xml:lang,attr,omitempty"`
	Title     string        `xml:"title"`
	Subtitle  string        `xml:"subtitle,omitempty"`
	Links     []atomLink    `xml:"link"`
	Generator atomGenerator `xml:"generator"`
	Updated   string        `xml:"updated"`
	ID        string        `xml:"id"`
	Entries   []atomEntry   `xml:"entry"`
}

type atomGenerator struct {
	Name string `xml:",chardata"`
}

type atomLink struct {
	Rel  string `xml:"rel,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`
	Href string `xml:"href,attr"`
}

type atomEntry struct {
	Language  string       `xml:"xml:lang,attr,omitempty"`
	Title     string       `xml:"title"`
	Published string       `xml:"published"`
	Updated   string       `xml:"updated"`
	Authors   []atomAuthor `xml:"author,omitempty"`
	Link      atomLink     `xml:"link"`
	ID        string       `xml:"id"`
	Content   atomContent  `xml:"content"`
}

type atomAuthor struct {
	Name string `xml:"name"`
}

type atomContent struct {
	Type string `xml:"type,attr"`
	Base string `xml:"xml:base,attr"`
	HTML string `xml:",chardata"`
}

func RenderFeed(site *content.Site, title, route string, pages []*content.Page) ([]byte, error) {
	entries, updated := feedEntries(site, pages)
	feedURL := absoluteURL(site.Config.BaseURL, route+"atom.xml")
	feed := atomFeed{
		Namespace: "http://www.w3.org/2005/Atom",
		Language:  site.Config.DefaultLanguage,
		Title:     title,
		Subtitle:  site.Config.Description,
		Links: []atomLink{
			{Rel: "self", Type: "application/atom+xml", Href: feedURL},
			{Rel: "alternate", Type: "text/html", Href: absoluteURL(site.Config.BaseURL, route)},
		},
		Generator: atomGenerator{Name: "deald"},
		Updated:   updated.Format(atomTimeFormat),
		ID:        feedURL,
		Entries:   entries,
	}
	data, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), append(data, '\n')...), nil
}

func feedEntries(site *content.Site, pages []*content.Page) ([]atomEntry, time.Time) {
	entries := make([]atomEntry, 0, len(pages))
	var updated time.Time
	for _, page := range pages {
		if page.Date == nil {
			continue
		}
		entryUpdated := page.Date.UTC()
		if page.Updated != nil {
			entryUpdated = page.Updated.UTC()
		}
		if updated.IsZero() || entryUpdated.After(updated) {
			updated = entryUpdated
		}
		pageURL := absoluteURL(site.Config.BaseURL, page.Route)
		entries = append(entries, atomEntry{
			Language:  site.Config.DefaultLanguage,
			Title:     page.Title,
			Published: page.Date.UTC().Format(atomTimeFormat),
			Updated:   entryUpdated.Format(atomTimeFormat),
			Authors:   feedAuthors(page, site.Config.Author),
			Link:      atomLink{Rel: "alternate", Type: "text/html", Href: pageURL},
			ID:        pageURL,
			Content: atomContent{
				Type: "html",
				Base: pageURL,
				HTML: string(page.HTML),
			},
		})
	}
	if updated.IsZero() {
		updated = site.BuildTime.UTC()
	}
	return entries, updated
}

func feedAuthors(page *content.Page, fallback string) []atomAuthor {
	authors := page.Authors
	if len(authors) == 0 && fallback != "" {
		authors = []string{fallback}
	}
	out := make([]atomAuthor, 0, len(authors))
	for _, author := range authors {
		if author != "" {
			out = append(out, atomAuthor{Name: author})
		}
	}
	return out
}

func absoluteURL(baseURL, route string) string {
	baseURL = strings.TrimRight(baseURL, "/")
	if route == "/" {
		return baseURL
	}
	return baseURL + route
}
