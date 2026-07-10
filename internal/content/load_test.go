package content

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseMarkdownDocumentStripsUTF8BOMBeforeFrontMatter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bom.md")
	input := "\uFEFF+++\n" +
		"title = \"BOM Title\"\n" +
		"[taxonomies]\n" +
		"tags = [\"one\", \"two\"]\n" +
		"+++\n\n" +
		"Body text.\n"

	if err := os.WriteFile(path, []byte(input), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	doc, err := parseMarkdownDocument(path, "bom.md")
	if err != nil {
		t.Fatalf("parse markdown: %v", err)
	}

	if doc.FrontMatter.Title != "BOM Title" {
		t.Fatalf("title = %q, want %q", doc.FrontMatter.Title, "BOM Title")
	}

	if got := doc.FrontMatter.Taxonomies["tags"]; len(got) != 2 || got[0] != "one" || got[1] != "two" {
		t.Fatalf("tags = %#v, want %#v", got, []string{"one", "two"})
	}

	if doc.Body != "\nBody text.\n" {
		t.Fatalf("body = %q, want %q", doc.Body, "\nBody text.\n")
	}
}

func TestParseMarkdownDocumentRejectsMalformedKnownFrontMatter(t *testing.T) {
	tests := []struct {
		name  string
		front string
	}{
		{name: "string instead of boolean", front: `draft = "true"`},
		{name: "invalid date", front: `date = "not-a-date"`},
		{name: "number instead of title", front: `title = 1`},
		{name: "string instead of authors", front: `authors = "Wesley"`},
		{name: "invalid taxonomy", front: "[taxonomies]\ntags = \"go\""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "invalid.md")
			input := "+++\n" + tt.front + "\n+++\n"
			if err := os.WriteFile(path, []byte(input), 0o644); err != nil {
				t.Fatalf("write fixture: %v", err)
			}
			if _, err := parseMarkdownDocument(path, "invalid.md"); err == nil {
				t.Fatal("malformed frontmatter was accepted")
			}
		})
	}
}

func TestParseMarkdownDocumentTreatsBareDatesAsUTC(t *testing.T) {
	path := filepath.Join(t.TempDir(), "date.md")
	input := "+++\n" +
		"date = 2025-07-17\n" +
		"updated = 2025-07-18\n" +
		"+++\n"
	if err := os.WriteFile(path, []byte(input), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	doc, err := parseMarkdownDocument(path, "date.md")
	if err != nil {
		t.Fatalf("parse markdown: %v", err)
	}
	for name, got := range map[string]*time.Time{"date": doc.FrontMatter.Date, "updated": doc.FrontMatter.Updated} {
		if got == nil || got.Format(time.RFC3339) != map[string]string{"date": "2025-07-17T00:00:00Z", "updated": "2025-07-18T00:00:00Z"}[name] {
			t.Fatalf("%s = %v", name, got)
		}
	}
}
