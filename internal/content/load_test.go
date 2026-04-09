package content

import (
	"os"
	"path/filepath"
	"testing"
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
