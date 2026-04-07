package site

import (
	"bytes"
	"compress/gzip"
	"context"
	"database/sql"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
	_ "modernc.org/sqlite"
)

type App struct {
	opts     Options
	cfg      Config
	renderer *Renderer
	db       *sql.DB
	mu       sync.RWMutex
	site     *Site
}

func NewApp(opts Options) (*App, error) {
	cfg, err := LoadConfig(opts.ConfigPath)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(opts.OutputDir, 0o755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(opts.DataDir, 0o755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(opts.DataDir, "site.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	if err := configureSQLite(db); err != nil {
		return nil, err
	}
	if err := initSchema(db); err != nil {
		return nil, err
	}

	renderer, err := NewRenderer()
	if err != nil {
		return nil, err
	}

	return &App{
		opts:     opts,
		cfg:      cfg,
		renderer: renderer,
		db:       db,
	}, nil
}

func (a *App) Options() Options {
	return a.opts
}

func (a *App) Close() error {
	if a.db != nil {
		return a.db.Close()
	}
	return nil
}

func configureSQLite(db *sql.DB) error {
	pragmas := []string{
		`PRAGMA journal_mode=WAL;`,
		`PRAGMA synchronous=NORMAL;`,
		`PRAGMA foreign_keys=ON;`,
		`PRAGMA temp_store=MEMORY;`,
	}
	for _, statement := range pragmas {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}
	return nil
}

func initSchema(db *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS documents (
  route TEXT PRIMARY KEY,
  kind TEXT NOT NULL,
  source_path TEXT NOT NULL,
  checksum TEXT NOT NULL,
  title TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS taxon_entries (
  taxonomy TEXT NOT NULL,
  term_slug TEXT NOT NULL,
  page_route TEXT NOT NULL,
  PRIMARY KEY (taxonomy, term_slug, page_route)
);`
	_, err := db.Exec(schema)
	return err
}

func (a *App) Rebuild(ctx context.Context) error {
	start := time.Now()
	siteModel, err := LoadSite(a.cfg, a.opts)
	if err != nil {
		return err
	}
	siteModel.RootOutputDir = a.opts.OutputDir
	if err := a.renderer.RenderSite(siteModel); err != nil {
		return err
	}
	if err := a.writeOutput(siteModel); err != nil {
		return err
	}
	if err := a.persistSite(ctx, siteModel); err != nil {
		return err
	}

	a.mu.Lock()
	a.site = siteModel
	a.mu.Unlock()
	log.Printf("rebuilt site in %s (%s)", time.Since(start).Round(time.Millisecond), siteModel)
	return nil
}

func (a *App) writeOutput(site *Site) error {
	if err := os.RemoveAll(a.opts.OutputDir); err != nil {
		return err
	}
	if err := os.MkdirAll(a.opts.OutputDir, 0o755); err != nil {
		return err
	}
	if err := copyTree(a.opts.StaticDir, a.opts.OutputDir, func(string) bool { return true }); err != nil {
		return err
	}
	if err := copyTree(a.opts.ContentDir, a.opts.OutputDir, func(rel string) bool {
		return filepath.Ext(rel) != ".md"
	}); err != nil {
		return err
	}

	for _, section := range site.Sections {
		doc, err := a.renderer.RenderSectionDocument(site, section)
		if err != nil {
			return err
		}
		if err := writeRouteDocument(a.opts.OutputDir, section.Route, doc); err != nil {
			return err
		}
	}
	for _, page := range site.Pages {
		doc, err := a.renderer.RenderPageDocument(site, page)
		if err != nil {
			return err
		}
		if err := writeRouteDocument(a.opts.OutputDir, page.Route, doc); err != nil {
			return err
		}
	}
	for name, terms := range site.Taxonomies {
		for _, term := range terms {
			if !term.Config.ShouldRender() {
				continue
			}
			doc, err := a.renderer.RenderTaxonomyDocument(site, name, term)
			if err != nil {
				return err
			}
			if err := writeRouteDocument(a.opts.OutputDir, term.Route, doc); err != nil {
				return err
			}
		}
	}

	if err := writeJSON(filepath.Join(a.opts.OutputDir, site.SearchIndexName), searchDocuments(site)); err != nil {
		return err
	}
	if err := writeNotFound(filepath.Join(a.opts.OutputDir, "404.html")); err != nil {
		return err
	}
	return nil
}

func writeRouteDocument(root, route string, doc []byte) error {
	if route == "/" {
		path := filepath.Join(root, "index.html")
		if err := os.WriteFile(path, doc, 0o644); err != nil {
			return err
		}
		return writeCompressedVariants(path, doc)
	}

	output := filepath.Join(root, filepath.FromSlash(strings.Trim(route, "/")), "index.html")
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(output, doc, 0o644); err != nil {
		return err
	}
	return writeCompressedVariants(output, doc)
}

func writeCompressedVariants(basePath string, plain []byte) error {
	if err := writeGzip(basePath+".gz", plain); err != nil {
		return err
	}
	if err := writeBrotli(basePath+".br", plain); err != nil {
		return err
	}
	if err := writeZstd(basePath+".zst", plain); err != nil {
		return err
	}
	return nil
}

func writeGzip(path string, plain []byte) error {
	var buf bytes.Buffer
	zw, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return err
	}
	if _, err := zw.Write(plain); err != nil {
		return err
	}
	if err := zw.Close(); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o644)
}

func writeBrotli(path string, plain []byte) error {
	var buf bytes.Buffer
	zw := brotli.NewWriterLevel(&buf, brotli.BestCompression)
	if _, err := zw.Write(plain); err != nil {
		return err
	}
	if err := zw.Close(); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o644)
}

func writeZstd(path string, plain []byte) error {
	var buf bytes.Buffer
	zw, err := zstd.NewWriter(&buf, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
	if err != nil {
		return err
	}
	if _, err := zw.Write(plain); err != nil {
		return err
	}
	if err := zw.Close(); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o644)
}

func copyTree(srcRoot, dstRoot string, include func(rel string) bool) error {
	return filepath.WalkDir(srcRoot, func(src string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(srcRoot, src)
		if err != nil {
			return err
		}
		if !include(rel) {
			return nil
		}
		dst := filepath.Join(dstRoot, rel)
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		return copyFile(src, dst)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

func writeNotFound(path string) error {
	html := `<!doctype html><html lang="en"><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><title>Not Found</title><style>body{font:16px/1.6 sans-serif;max-width:60rem;margin:4rem auto;padding:0 1rem}a{color:#056}</style><h1>404</h1><p>The page you asked for is not here.</p><p><a href="/">Return home</a></p></html>`
	if err := os.WriteFile(path, []byte(html), 0o644); err != nil {
		return err
	}
	return writeCompressedVariants(path, []byte(html))
}

func (a *App) persistSite(ctx context.Context, site *Site) error {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM documents`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM taxon_entries`); err != nil {
		return err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO documents(route, kind, source_path, checksum, title, updated_at) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, section := range site.Sections {
		if _, err := stmt.ExecContext(ctx, section.Route, "section", section.SourcePath, section.Checksum, section.Title, now); err != nil {
			return err
		}
	}
	for _, page := range site.Pages {
		if _, err := stmt.ExecContext(ctx, page.Route, "page", page.SourcePath, page.Checksum, page.Title, now); err != nil {
			return err
		}
	}
	for name, terms := range site.Taxonomies {
		for _, term := range terms {
			for _, page := range term.Pages {
				if _, err := tx.ExecContext(ctx, `INSERT INTO taxon_entries(taxonomy, term_slug, page_route) VALUES (?, ?, ?)`, name, term.Slug, page.Route); err != nil {
					return err
				}
			}
		}
	}

	return tx.Commit()
}

func (a *App) ListenAndServe(ctx context.Context, addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/__cms/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("/__cms/rebuild", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := a.Rebuild(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	})

	fsHandler := http.FileServer(http.Dir(a.opts.OutputDir))
	mux.Handle("/", fsHandler)

	server := &http.Server{
		Addr:    addr,
		Handler: loggingMiddleware(mux),
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/__cms/") {
			w.Header().Set("Cache-Control", "public, max-age=60")
		}
		next.ServeHTTP(w, r)
	})
}
