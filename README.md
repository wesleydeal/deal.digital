# Hello, I'm Wesley

## [🔗deal.digital](https://deal.digital/)

This repo is for my **🚧 under construction 🚧** web playground. It started as a Zola site and now includes a Go-based CMS/build pipeline that keeps markdown as the source of truth while caching rendered HTML for Caddy to serve like a static site.

## Go CMS

The new CMS lives under [`cmd/deald`](/home/wesley/deal.digital/cmd/deald/main.go) and [`internal/site`](/home/wesley/deal.digital/internal/site/app.go).

It currently does the following:

- parses existing markdown files with Zola-style TOML frontmatter
- normalizes pages, sections, and taxonomies into a Go-native site graph
- supports the current shortcode set: `blockquote`, `fitimg`, `wave`, and `raw`
- renders HTML into [`build/site`](/home/wesley/deal.digital/build/site) and writes `.gz`, `.br`, and `.zst` sidecar files for each page
- copies static assets and page-bundle assets into the Caddy-served tree
- stores build metadata in SQLite at [`build/data/site.db`](/home/wesley/deal.digital/build/data/site.db) with WAL mode enabled
- exposes `POST /__cms/rebuild` and `GET /__cms/healthz` for rebuild/admin hooks

## Local Use

Build the site:

```bash
go run ./cmd/deald build
```

Run the CMS and local file server:

```bash
go run ./cmd/deald serve
```

Run Caddy against the generated site tree:

```bash
caddy run --config ./Caddyfile
```

The provided [`Caddyfile`](/home/wesley/deal.digital/Caddyfile) serves `build/site` directly and proxies `/__cms/*` to the Go process. You can override paths with `SITE_ROOT` and `CMS_UPSTREAM`.
