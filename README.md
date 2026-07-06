# Hello, I'm Wesley

## [🔗deal.digital](https://deal.digital/)

This repo is for my **🚧 under construction 🚧** web playground. It started as a Zola site and now uses a small Go-based static site generator with markdown as the source of truth.

## Static Site Generator

The generator entrypoint lives under [`cmd/deald`](/home/wesley/deal.digital/cmd/deald/main.go) and is split into:

- [`internal/content`](/home/wesley/deal.digital/internal/content/load.go) for flatfile loading, frontmatter parsing, and the site graph
- [`internal/render`](/home/wesley/deal.digital/internal/render/renderer.go) for markdown, templates, shortcodes, and page rendering
- [`internal/runtime`](/home/wesley/deal.digital/internal/runtime/watch.go) for watched rebuilds during local development
- [`internal/site`](/home/wesley/deal.digital/internal/site/app.go) for orchestration, build output, compression, and the local development file server
- [`templates`](/home/wesley/deal.digital/templates/base.tmpl) for the live on-disk Go templates

It currently does the following:

- parses existing markdown files with Zola-style TOML frontmatter
- normalizes pages, sections, and taxonomies into a Go-native site graph
- supports the current shortcode set: `blockquote`, `fitimg`, `wave`, and `raw`
- renders HTML into [`build/site`](/home/wesley/deal.digital/build/site)
- writes `.gz`, `.br`, and `.zst` sidecar files for generated output assets during `build`, using parallel workers, skipping any compressed variant that would be larger than the source, blacklisting obviously unhelpful formats like SWF/media/archive assets, and only brotli-compressing web-text formats
- copies static assets and page-bundle assets into the Caddy-served tree

## Local Use

Build the site:

```bash
go run ./cmd/deald build
```

Run the generator and local file server:

```bash
go run ./cmd/deald serve
```

`serve` watches the config, content, static, and template directories and rebuilds on change. It intentionally skips precompression so local iteration stays fast.

Run Caddy against the generated site tree:

```bash
caddy run --config ./Caddyfile
```

The provided [`Caddyfile`](/home/wesley/deal.digital/Caddyfile) serves `build/site` directly with precompressed asset support. You can override the site root with `SITE_ROOT`.
