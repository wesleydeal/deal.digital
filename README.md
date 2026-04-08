# Hello, I'm Wesley

## [🔗deal.digital](https://deal.digital/)

This repo is for my **🚧 under construction 🚧** web playground. It started as a Zola site and now includes a Go-based flatfile CMS that keeps markdown as the source of truth while caching rendered HTML for Caddy to serve like a static site.

## Go CMS

The CMS entrypoint lives under [`cmd/deald`](/home/wesley/deal.digital/cmd/deald/main.go) and is split into:

- [`internal/content`](/home/wesley/deal.digital/internal/content/load.go) for flatfile loading, frontmatter parsing, and the site graph
- [`internal/render`](/home/wesley/deal.digital/internal/render/renderer.go) for markdown, templates, shortcodes, and page rendering
- [`internal/store`](/home/wesley/deal.digital/internal/store/sqlite.go) for SQLite-backed runtime state
- [`internal/runtime`](/home/wesley/deal.digital/internal/runtime/watch.go) for watched rebuilds during local development
- [`internal/cms`](/home/wesley/deal.digital/internal/cms/app.go) for orchestration, HTTP routes, and build output
- [`templates`](/home/wesley/deal.digital/templates/base.tmpl) for the live on-disk Go templates

It currently does the following:

- parses existing markdown files with Zola-style TOML frontmatter
- normalizes pages, sections, and taxonomies into a Go-native site graph
- supports the current shortcode set: `blockquote`, `fitimg`, `wave`, and `raw`
- renders HTML into [`build/site`](/home/wesley/deal.digital/build/site)
- writes `.gz`, `.br`, and `.zst` sidecar files for generated output assets during `build`, using parallel workers, skipping any compressed variant that would be larger than the source, blacklisting obviously unhelpful formats like SWF/media/archive assets, and only brotli-compressing web-text formats
- copies static assets and page-bundle assets into the Caddy-served tree
- stores build metadata in SQLite at [`build/data/site.db`](/home/wesley/deal.digital/build/data/site.db) with WAL mode enabled
- exposes a public dynamic route namespace under `/_/` for server-side interactivity
- exposes a sample SQLite-backed interactive counter at `GET/POST /_/fragments/sample-counter`, wired from `/sample/` via µJS
- exposes `POST /_/admin/rebuild` and `GET /_/healthz` for rebuild/admin hooks

## Local Use

Build the site:

```bash
go run ./cmd/deald build
```

Run the CMS and local file server:

```bash
go run ./cmd/deald serve
```

`serve` watches the config, content, static, and template directories and rebuilds on change. It intentionally skips precompression so local iteration stays fast.

Run Caddy against the generated site tree:

```bash
caddy run --config ./Caddyfile
```

The provided [`Caddyfile`](/home/wesley/deal.digital/Caddyfile) serves `build/site` directly and proxies `/_/*` to the Go process. That is the answer to the “how does interactivity work through Caddy?” question: static routes are served from disk, while dynamic fragment/admin routes are reverse proxied to Go under one explicit namespace. You can override paths with `SITE_ROOT` and `CMS_UPSTREAM`.
