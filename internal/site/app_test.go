package site

import (
	"bytes"
	"context"
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"deal.digital/internal/content"
)

const testTemplates = `{{define "home"}}home{{end}}
{{define "page"}}{{.Page.HTML}}{{end}}
{{define "blog_page"}}page{{end}}
{{define "list"}}list{{end}}
{{define "blog_list"}}list{{end}}
{{define "taxonomy"}}taxonomy{{end}}
{{define "404"}}not found{{end}}`

func TestRebuildReplacesOutput(t *testing.T) {
	app := newTestGenerator(t, `
base_url = "https://example.test"
title = "Example"
default_language = "en"
`, map[string]string{
		"_index.md": "+++\ntitle = \"Example\"\n+++\n",
		"draft.md":  "+++\ntitle = \"Draft\"\ndraft = true\n+++\n",
	})
	app.opts.IncludeDrafts = true
	if err := app.Rebuild(context.Background()); err != nil {
		t.Fatalf("build with drafts: %v", err)
	}
	draftPath := filepath.Join(app.opts.OutputDir, "draft", "index.html")
	if _, err := os.Stat(draftPath); err != nil {
		t.Fatalf("draft output missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(app.opts.OutputDir, "search_index.en.json")); !os.IsNotExist(err) {
		t.Fatalf("disabled search index exists: %v", err)
	}
	if _, err := os.Stat(filepath.Join(app.opts.OutputDir, "atom.xml")); !os.IsNotExist(err) {
		t.Fatalf("disabled feed exists: %v", err)
	}

	app.opts.IncludeDrafts = false
	if err := app.Rebuild(context.Background()); err != nil {
		t.Fatalf("build without drafts: %v", err)
	}
	if _, err := os.Stat(draftPath); !os.IsNotExist(err) {
		t.Fatalf("draft output still exists: %v", err)
	}
}

func TestRebuildWritesConfiguredSearchAndFeeds(t *testing.T) {
	app := newTestGenerator(t, `
base_url = "https://example.test"
title = "Example"
description = "Example site"
default_language = "en"
author = "Example Author"
build_search_index = true
generate_feeds = true
taxonomies = [{name = "tags", feed = true}]
`, map[string]string{
		"_index.md": "+++\ntitle = \"Example\"\n+++\n",
		"post.md":   "+++\ntitle = \"Post\"\ndate = 2024-01-02\n[taxonomies]\ntags = [\"go\"]\n+++\nHello **world**.",
	})
	if err := app.Rebuild(context.Background()); err != nil {
		t.Fatalf("build: %v", err)
	}

	for _, path := range []string{
		filepath.Join(app.opts.OutputDir, "search_index.en.json"),
		filepath.Join(app.opts.OutputDir, "atom.xml"),
		filepath.Join(app.opts.OutputDir, "tags", "go", "atom.xml"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected output %s: %v", path, err)
		}
	}

	feed, err := os.ReadFile(filepath.Join(app.opts.OutputDir, "atom.xml"))
	if err != nil {
		t.Fatalf("read feed: %v", err)
	}
	if !strings.Contains(string(feed), `<content type="html"`) {
		t.Fatalf("feed does not contain rendered HTML: %s", feed)
	}
	var document struct {
		XMLName xml.Name `xml:"feed"`
	}
	if err := xml.Unmarshal(feed, &document); err != nil {
		t.Fatalf("parse feed XML: %v", err)
	}
}

func TestRebuildWritesPageVariants(t *testing.T) {
	app := newTestGenerator(t, `
base_url = "https://example.test"
title = "Example"
default_language = "en"
`, map[string]string{
		"_index.md": "+++\ntitle = \"Example\"\n+++\n",
		"resume/index.md": `+++
title = "Resume"
[variants]
default = "general"
items = [
  { id = "general", label = "General IT" },
  { id = "neteng", label = "Network Engineering" },
]
+++
Shared {% variant(include="neteng") %}network{% endvariant %} experience.

{% variant(include="neteng") %}
* Network-only list item
{% endvariant %}

{% variant(exclude="neteng") %}General-only text.{% endvariant %}

{{ variant_menu() }}`,
	})
	if err := app.Rebuild(context.Background()); err != nil {
		t.Fatalf("build: %v", err)
	}

	generalPath := filepath.Join(app.opts.OutputDir, "resume", "index.html")
	netengPath := filepath.Join(app.opts.OutputDir, "resume", "neteng", "index.html")
	general, err := os.ReadFile(generalPath)
	if err != nil {
		t.Fatalf("read general variant: %v", err)
	}
	neteng, err := os.ReadFile(netengPath)
	if err != nil {
		t.Fatalf("read network variant: %v", err)
	}
	if strings.Contains(string(general), "Network-only") || !strings.Contains(string(general), "General-only text.") {
		t.Fatalf("general variant content is incorrect: %s", general)
	}
	if strings.Contains(string(neteng), "General-only") || !strings.Contains(string(neteng), "Network-only list item") {
		t.Fatalf("network variant content is incorrect: %s", neteng)
	}
	if !strings.Contains(string(general), `<a href="/resume/neteng/">Network Engineering</a>`) || !strings.Contains(string(neteng), `<a href="/resume/">General IT</a>`) {
		t.Fatalf("variant navigation is incorrect")
	}
	if _, err := os.Stat(filepath.Join(app.opts.OutputDir, "resume", "general", "index.html")); !os.IsNotExist(err) {
		t.Fatalf("default variant duplicate exists: %v", err)
	}
}

func TestRebuildWarnsOnOverwrite(t *testing.T) {
	app := newTestGenerator(t, `
base_url = "https://example.test"
title = "Example"
default_language = "en"
`, map[string]string{
		"_index.md": "+++\ntitle = \"Example\"\n+++\n",
	})
	mustWriteFile(t, filepath.Join(app.opts.StaticDir, "index.html"), "static")

	previous := log.Writer()
	var logs bytes.Buffer
	log.SetOutput(&logs)
	defer log.SetOutput(previous)

	if err := app.Rebuild(context.Background()); err != nil {
		t.Fatalf("build: %v", err)
	}
	if !strings.Contains(logs.String(), "warning: section _index.md overwrites") {
		t.Fatalf("missing overwrite warning: %s", logs.String())
	}
}

func newTestGenerator(t *testing.T, config string, files map[string]string) *Generator {
	t.Helper()
	root := t.TempDir()
	configPath := filepath.Join(root, "config.toml")
	contentDir := filepath.Join(root, "content")
	staticDir := filepath.Join(root, "static")
	templateDir := filepath.Join(root, "templates")
	if err := os.MkdirAll(staticDir, 0o755); err != nil {
		t.Fatalf("create static directory: %v", err)
	}
	mustWriteFile(t, configPath, config)
	mustWriteFile(t, filepath.Join(templateDir, "site.tmpl"), testTemplates)
	for relative, body := range files {
		mustWriteFile(t, filepath.Join(contentDir, relative), body)
	}

	app, err := NewGenerator(content.Options{
		ConfigPath:  configPath,
		ContentDir:  contentDir,
		StaticDir:   staticDir,
		TemplateDir: templateDir,
		OutputDir:   filepath.Join(root, "build", "site"),
	})
	if err != nil {
		t.Fatalf("new generator: %v", err)
	}
	return app
}

func mustWriteFile(t *testing.T, path string, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
