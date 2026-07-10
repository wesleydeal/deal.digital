package site

import (
	"context"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"deal.digital/internal/content"
	"deal.digital/internal/render"
	siteruntime "deal.digital/internal/runtime"
)

type Generator struct {
	opts     content.Options
	renderer *render.Renderer
}

type buildTiming struct {
	render   time.Duration
	compress time.Duration
}

func NewGenerator(opts content.Options) (*Generator, error) {
	renderer, err := render.NewRenderer(opts.TemplateDir)
	if err != nil {
		return nil, err
	}

	return &Generator{
		opts:     opts,
		renderer: renderer,
	}, nil
}

func (a *Generator) Options() content.Options {
	return a.opts
}

func (a *Generator) Rebuild(ctx context.Context) error {
	start := time.Now().UTC()
	cfg, err := content.LoadConfig(a.opts.ConfigPath)
	if err != nil {
		return err
	}
	siteModel, err := content.LoadSite(cfg, a.opts)
	if err != nil {
		return err
	}
	stagingDir, err := createStagingDir(a.opts.OutputDir)
	if err != nil {
		return err
	}
	keepStaging := true
	defer func() {
		if keepStaging {
			_ = os.RemoveAll(stagingDir)
		}
	}()

	timings, err := a.writeOutput(siteModel, stagingDir)
	if err != nil {
		return err
	}
	if err := replaceOutput(stagingDir, a.opts.OutputDir); err != nil {
		return err
	}
	keepStaging = false
	total := time.Since(start).Round(time.Millisecond)
	renderDur := timings.render.Round(time.Millisecond)
	if a.opts.CompressOutput {
		log.Printf("rebuilt site in %s (render=%s, compress=%s, %s)", total, renderDur, timings.compress.Round(time.Millisecond), siteModel)
		return nil
	}
	log.Printf("rebuilt site in %s (render=%s, %s)", total, renderDur, siteModel)
	return nil
}

func (a *Generator) Watch(ctx context.Context) {
	siteruntime.Watch(ctx, siteruntime.NormalizeRoots([]string{
		a.opts.ConfigPath,
		a.opts.ContentDir,
		a.opts.StaticDir,
		a.opts.TemplateDir,
	}), a.Rebuild)
}

func (a *Generator) writeOutput(site *content.Site, outputDir string) (buildTiming, error) {
	start := time.Now()
	if err := a.renderer.RenderSite(site); err != nil {
		return buildTiming{}, err
	}
	if err := copyTree(a.opts.StaticDir, outputDir, func(string) bool { return true }); err != nil {
		return buildTiming{}, err
	}
	if err := copyTree(a.opts.ContentDir, outputDir, func(rel string) bool {
		return filepath.Ext(rel) != ".md"
	}); err != nil {
		return buildTiming{}, err
	}

	for _, section := range site.Sections {
		doc, err := a.renderer.RenderSectionDocument(site, section)
		if err != nil {
			return buildTiming{}, err
		}
		if err := writeRouteDocument(outputDir, section.Route, doc, "section "+section.RelativePath); err != nil {
			return buildTiming{}, err
		}
	}
	for _, page := range site.Pages {
		doc, err := a.renderer.RenderPageDocument(site, page)
		if err != nil {
			return buildTiming{}, err
		}
		if err := writeRouteDocument(outputDir, page.Route, doc, "page "+page.RelativePath); err != nil {
			return buildTiming{}, err
		}
		for _, variant := range page.Variants {
			if variant.ID == page.DefaultVariant {
				continue
			}
			doc, err := a.renderer.RenderPageVariantDocument(site, page, variant)
			if err != nil {
				return buildTiming{}, err
			}
			if err := writeRouteDocument(outputDir, variant.Route, doc, "page variant "+page.RelativePath+"/"+variant.ID); err != nil {
				return buildTiming{}, err
			}
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
			if err := writeRouteDocument(outputDir, term.Route, doc, "taxonomy "+name+"/"+term.Name); err != nil {
				return buildTiming{}, err
			}
		}
	}

	if site.Config.BuildSearchIndex {
		searchIndex := filepath.Join(outputDir, site.SearchIndexName)
		if err := warnOverwrite(searchIndex, "search index"); err != nil {
			return buildTiming{}, err
		}
		if err := render.WriteJSON(searchIndex, content.SearchDocuments(site)); err != nil {
			return buildTiming{}, err
		}
	}
	notFound, err := a.renderer.RenderNotFoundDocument(site)
	if err != nil {
		return buildTiming{}, err
	}
	if err := writeOutputFile(filepath.Join(outputDir, "404.html"), notFound, "404 page"); err != nil {
		return buildTiming{}, err
	}
	if site.Config.GenerateFeeds {
		feed, err := render.RenderFeed(site, site.Config.Title, "/", site.Pages)
		if err != nil {
			return buildTiming{}, err
		}
		if err := writeRouteFile(outputDir, "/", "atom.xml", feed, "site feed"); err != nil {
			return buildTiming{}, err
		}
		for name, terms := range site.Taxonomies {
			for _, term := range terms {
				if !term.Config.ShouldRender() || !term.Config.Feed {
					continue
				}
				feed, err := render.RenderFeed(site, site.Config.Title+" - "+term.Name, term.Route, term.Pages)
				if err != nil {
					return buildTiming{}, err
				}
				if err := writeRouteFile(outputDir, term.Route, "atom.xml", feed, "taxonomy feed "+name+"/"+term.Name); err != nil {
					return buildTiming{}, err
				}
			}
		}
	}

	timings := buildTiming{render: time.Since(start)}
	if a.opts.CompressOutput {
		compressStart := time.Now()
		if err := compressTree(outputDir, func(done, total int) {
			log.Printf("compressing output: %d/%d files", done, total)
		}); err != nil {
			return buildTiming{}, err
		}
		timings.compress = time.Since(compressStart)
	}
	return timings, nil
}
func createStagingDir(outputDir string) (string, error) {
	outputDir = filepath.Clean(outputDir)
	parent := filepath.Dir(outputDir)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return "", err
	}
	return os.MkdirTemp(parent, "."+filepath.Base(outputDir)+"-")
}

func replaceOutput(stagingDir, outputDir string) error {
	if err := os.RemoveAll(outputDir); err != nil {
		return err
	}
	return os.Rename(stagingDir, outputDir)
}

func writeRouteDocument(root, route string, doc []byte, source string) error {
	return writeRouteFile(root, route, "index.html", doc, source)
}

func writeRouteFile(root, route, name string, data []byte, source string) error {
	output := filepath.Join(root, name)
	if route != "/" {
		output = filepath.Join(root, filepath.FromSlash(strings.Trim(route, "/")), name)
	}
	return writeOutputFile(output, data, source)
}

func writeOutputFile(path string, data []byte, source string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := warnOverwrite(path, source); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func warnOverwrite(path string, source string) error {
	if _, err := os.Lstat(path); err == nil {
		log.Printf("warning: %s overwrites %s", source, path)
		return nil
	} else if os.IsNotExist(err) {
		return nil
	} else {
		return err
	}
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
		if err := warnOverwrite(dst, "asset "+src); err != nil {
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
