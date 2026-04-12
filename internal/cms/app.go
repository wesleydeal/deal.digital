package cms

import (
	"context"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"deal.digital/internal/content"
	"deal.digital/internal/render"
	cmsruntime "deal.digital/internal/runtime"
	"deal.digital/internal/store"
)

type App struct {
	opts     content.Options
	cfg      content.Config
	renderer *render.Renderer
	store    *store.SQLite
	mu       sync.RWMutex
	site     *content.Site
}

type buildTiming struct {
	render   time.Duration
	compress time.Duration
}

func NewApp(opts content.Options) (*App, error) {
	cfg, err := content.LoadConfig(opts.ConfigPath)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(opts.OutputDir, 0o755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(opts.DataDir, 0o755); err != nil {
		return nil, err
	}

	sqliteStore, err := store.NewSQLite(opts.DataDir)
	if err != nil {
		return nil, err
	}
	renderer, err := render.NewRenderer(opts.TemplateDir)
	if err != nil {
		_ = sqliteStore.Close()
		return nil, err
	}

	return &App{
		opts:     opts,
		cfg:      cfg,
		renderer: renderer,
		store:    sqliteStore,
	}, nil
}

func (a *App) Options() content.Options {
	return a.opts
}

func (a *App) Close() error {
	if a.store != nil {
		return a.store.Close()
	}
	return nil
}

func (a *App) Rebuild(ctx context.Context) error {
	start := time.Now().UTC()
	cfg, err := content.LoadConfig(a.opts.ConfigPath)
	if err != nil {
		return err
	}
	siteModel, err := content.LoadSite(cfg, a.opts)
	if err != nil {
		return err
	}
	siteModel.RootOutputDir = a.opts.OutputDir
	timings, err := a.writeOutput(siteModel)
	if err != nil {
		return err
	}
	if err := a.store.PersistSite(ctx, siteModel); err != nil {
		return err
	}

	a.mu.Lock()
	a.cfg = cfg
	a.site = siteModel
	a.mu.Unlock()
	total := time.Since(start).Round(time.Millisecond)
	renderDur := timings.render.Round(time.Millisecond)
	if a.opts.CompressOutput {
		log.Printf("rebuilt site in %s (render=%s, compress=%s, %s)", total, renderDur, timings.compress.Round(time.Millisecond), siteModel)
		return nil
	}
	log.Printf("rebuilt site in %s (render=%s, %s)", total, renderDur, siteModel)
	return nil
}

func (a *App) Watch(ctx context.Context) {
	cmsruntime.Watch(ctx, cmsruntime.NormalizeRoots([]string{
		a.opts.ConfigPath,
		a.opts.ContentDir,
		a.opts.StaticDir,
		a.opts.TemplateDir,
	}), a.Rebuild)
}

func (a *App) writeOutput(site *content.Site) (buildTiming, error) {
	start := time.Now()
	if err := a.renderer.RenderSite(site); err != nil {
		return buildTiming{}, err
	}
	if err := os.MkdirAll(a.opts.OutputDir, 0o755); err != nil {
		return buildTiming{}, err
	}
	if err := copyTree(a.opts.StaticDir, a.opts.OutputDir, func(string) bool { return true }); err != nil {
		return buildTiming{}, err
	}
	if err := copyTree(a.opts.ContentDir, a.opts.OutputDir, func(rel string) bool {
		return filepath.Ext(rel) != ".md"
	}); err != nil {
		return buildTiming{}, err
	}

	for _, section := range site.Sections {
		doc, err := a.renderer.RenderSectionDocument(site, section)
		if err != nil {
			return buildTiming{}, err
		}
		if err := writeRouteDocument(a.opts.OutputDir, section.Route, doc); err != nil {
			return buildTiming{}, err
		}
	}
	for _, page := range site.Pages {
		doc, err := a.renderer.RenderPageDocument(site, page)
		if err != nil {
			return buildTiming{}, err
		}
		if err := writeRouteDocument(a.opts.OutputDir, page.Route, doc); err != nil {
			return buildTiming{}, err
		}
	}
	for name, terms := range site.Taxonomies {
		for _, term := range terms {
			if !term.Config.ShouldRender() {
				continue
			}
			doc, err := a.renderer.RenderTaxonomyDocument(site, name, term)
			if err != nil {
				return buildTiming{}, err
			}
			if err := writeRouteDocument(a.opts.OutputDir, term.Route, doc); err != nil {
				return buildTiming{}, err
			}
		}
	}

	if err := render.WriteJSON(filepath.Join(a.opts.OutputDir, site.SearchIndexName), content.SearchDocuments(site)); err != nil {
		return buildTiming{}, err
	}
	notFound, err := a.renderer.RenderNotFoundDocument(site)
	if err != nil {
		return buildTiming{}, err
	}
	if err := os.WriteFile(filepath.Join(a.opts.OutputDir, "404.html"), notFound, 0o644); err != nil {
		return buildTiming{}, err
	}

	timings := buildTiming{render: time.Since(start)}
	if a.opts.CompressOutput {
		compressStart := time.Now()
		if err := compressTree(a.opts.OutputDir, func(done, total int) {
			log.Printf("compressing output: %d/%d files", done, total)
		}); err != nil {
			return buildTiming{}, err
		}
		timings.compress = time.Since(compressStart)
	}
	return timings, nil
}
func writeRouteDocument(root, route string, doc []byte) error {
	if route == "/" {
		return os.WriteFile(filepath.Join(root, "index.html"), doc, 0o644)
	}

	output := filepath.Join(root, filepath.FromSlash(strings.Trim(route, "/")), "index.html")
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return err
	}
	return os.WriteFile(output, doc, 0o644)
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
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	dstInfo, err := os.Stat(dst)
	if err == nil && dstInfo.Size() == srcInfo.Size() && dstInfo.ModTime().Equal(srcInfo.ModTime()) {
		return nil
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	return os.Chtimes(dst, srcInfo.ModTime(), srcInfo.ModTime())
}
