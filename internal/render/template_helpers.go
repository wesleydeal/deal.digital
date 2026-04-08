package render

import (
	"fmt"
	"html"
	"html/template"
	"sort"
	"strings"

	"deal.digital/internal/content"
)

func colorFor(page *content.Page, section *content.Section) string {
	if page != nil {
		if s, ok := content.StringValue(page.Extra["color"]); ok && s != "" {
			return s
		}
	}
	if section != nil {
		if s, ok := content.StringValue(section.Extra["color"]); ok && s != "" {
			return s
		}
	}
	return ""
}

func documentTitle(data templateData) string {
	switch {
	case data.Page != nil:
		return content.DisplayTitle(data.Page) + " - Wesley Deal"
	case data.Section != nil && data.Section.Route != "/" && data.Section.Title != "":
		return data.Section.Title + " - Wesley Deal"
	case data.Term != nil && data.Term.Name != "":
		return data.Term.Name + " - Wesley Deal"
	default:
		return "Wesley Deal"
	}
}

func bodyID(route string) string {
	trimmed := strings.Trim(route, "/")
	if trimmed == "" {
		return "generic"
	}
	return strings.ReplaceAll(trimmed, "/", "-")
}

func extraBool(m map[string]any, key string) bool {
	if m == nil {
		return false
	}
	value, ok := m[key]
	if !ok {
		return false
	}
	b, _ := value.(bool)
	return b
}

func showBreadcrumbs(site *content.Site, page *content.Page) bool {
	if site == nil || page == nil {
		return false
	}
	if !extraBool(site.Config.Extra, "breadcrumbs") {
		return false
	}
	if extraBool(page.Extra, "nobreadcrumbs") {
		return false
	}
	return true
}

func breadcrumbsHTML(site *content.Site, page *content.Page) template.HTML {
	if !showBreadcrumbs(site, page) {
		return ""
	}

	var parts []string
	parts = append(parts, `<a href="/">~</a>`)

	routes := ancestorRoutes(site, page.ParentSection)
	for _, route := range routes {
		section := site.SectionByRoute[route]
		if section == nil || section.Route == "/" {
			continue
		}
		parts = append(parts, `<span>/</span>`)
		parts = append(parts, fmt.Sprintf(`<a href="%s">%s</a>`, html.EscapeString(section.Route), html.EscapeString(section.Title)))
	}
	parts = append(parts, `<span>/</span>`)
	parts = append(parts, fmt.Sprintf(`<a href="%s">●</a>`, html.EscapeString(page.Route)))
	return template.HTML(strings.Join(parts, "\n"))
}

func ancestorRoutes(site *content.Site, route string) []string {
	var routes []string
	for route != "" && route != "/" {
		section := site.SectionByRoute[route]
		if section == nil {
			break
		}
		routes = append(routes, section.Route)
		if section.ParentSection == section.Route {
			break
		}
		route = section.ParentSection
	}
	for i, j := 0, len(routes)-1; i < j; i, j = i+1, j-1 {
		routes[i], routes[j] = routes[j], routes[i]
	}
	return routes
}

func sourceURL(page *content.Page, section *content.Section) string {
	base := "https://github.com/wesleydeal/deal.digital"
	switch {
	case page != nil && page.RelativePath != "":
		return base + "/blob/master/content/" + page.RelativePath
	case section != nil && section.RelativePath != "":
		return base + "/blob/master/content/" + section.RelativePath
	default:
		return base
	}
}

func siteTreeHTML(site *content.Site) template.HTML {
	if site == nil {
		return ""
	}
	root := site.SectionByRoute["/"]
	if root == nil {
		return `<a href="/">deal.digital/</a>`
	}

	var lines []string
	lines = append(lines, `<a href="/">deal.digital/</a>`)
	children := treeChildren(root.Subsections, root.Pages)
	renderTreeChildren(&lines, "", children)

	return template.HTML(strings.Join(lines, "<br>"))
}

type treeNode struct {
	isSection bool
	section   *content.Section
	page      *content.Page
}

func treeChildren(sections []*content.Section, pages []*content.Page) []treeNode {
	out := make([]treeNode, 0, len(sections)+len(pages))
	sort.SliceStable(sections, func(i, j int) bool { return sections[i].Route < sections[j].Route })
	sort.SliceStable(pages, func(i, j int) bool { return pages[i].Route < pages[j].Route })
	for _, section := range sections {
		out = append(out, treeNode{isSection: true, section: section})
	}
	for _, page := range pages {
		out = append(out, treeNode{page: page})
	}
	return out
}

func renderTreeChildren(lines *[]string, prefix string, nodes []treeNode) {
	for i, node := range nodes {
		connector := "├─"
		nextPrefix := prefix + "│ "
		if i == len(nodes)-1 {
			connector = "└─"
			nextPrefix = prefix + "  "
		}
		if node.isSection {
			*lines = append(*lines, prefix+connector+treeSectionLink(node.section))
			children := treeChildren(node.section.Subsections, node.section.Pages)
			if len(children) > 0 {
				renderTreeChildren(lines, nextPrefix, children)
			}
			continue
		}
		*lines = append(*lines, prefix+connector+treePageLink(node.page))
	}
}

func treeSectionLink(section *content.Section) string {
	return fmt.Sprintf(`<a href="%s">%s/</a>`, html.EscapeString(section.Route), html.EscapeString(lastRouteSegment(section.Route)))
}

func treePageLink(page *content.Page) string {
	return fmt.Sprintf(`<a href="%s">%s</a>`, html.EscapeString(page.Route), html.EscapeString(lastRouteSegment(page.Route)))
}

func lastRouteSegment(route string) string {
	trimmed := strings.Trim(route, "/")
	if trimmed == "" {
		return "deal.digital"
	}
	parts := strings.Split(trimmed, "/")
	return parts[len(parts)-1]
}
