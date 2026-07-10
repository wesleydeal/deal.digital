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
- supports Markdown, page variants, and the `blockquote`, `fitimg`, `wave`, and `raw` shortcodes
- renders each build into a clean staging tree before replacing [`build/site`](/home/wesley/deal.digital/build/site)
- writes `.gz`, `.br`, and `.zst` sidecar files for generated output assets during `build`, using parallel workers, skipping any compressed variant that would be larger than the source, blacklisting obviously unhelpful formats like SWF/media/archive assets, and only brotli-compressing web-text formats
- copies static assets and page-bundle assets into the Caddy-served tree
- writes `/atom.xml` when `generate_feeds = true`; renderable taxonomies with `feed = true` receive `/taxonomy/term/atom.xml`
- writes `search_index.<language>.json` only when `build_search_index = true`

## Content Authoring

Pages use TOML frontmatter between `+++` markers. Untagged Markdown renders in every version of a page.

### Page Variants

Declare page variants in frontmatter. The declared default renders at the page's normal route; every other variant renders below it. For example, the following page renders at `/resume/`, `/resume/neteng/`, and `/resume/security/`:

```toml
[variants]
default = "general"
items = [
  { id = "general", label = "General IT" },
  { id = "neteng", label = "Network Engineering" },
  { id = "security", label = "Security Engineering" },
]
```

Variant IDs must be unique lowercase letters, numbers, or hyphens. Use `variant` blocks to include or exclude Markdown for named variants. The block can wrap inline text, a whole list item, or a larger Markdown block.

```md
Untagged content appears everywhere.

{% variant(include="neteng,security") %}
* This complete list item appears only in network and security variants.
{% endvariant %}

This {% variant(exclude="neteng") %}generalist {% endvariant %}description changes by variant.

{{ variant_menu() }}
```

`variant_menu()` creates a `<menu>` with one `<li>` per variant. The current variant is plain text with `aria-current="page"`; all other variants are links. Variant blocks and `variant_menu()` are valid only on pages that declare `[variants]`.

### Shortcodes

All shortcode arguments use `name="value"` syntax.

- Image: `{{ fitimg(path="/path/to/image.png", width=640, height=360) }}` writes a linked, lazy-loaded image. `height` is optional; numeric dimensions become pixels.
- Divider: `{{ wave(width="100%", color="var(--color-fg)") }}` writes the repeating wave divider. Both arguments are optional.
- Quote:

  ```md
  {% blockquote(author="Name", date="2025-01-01", url="https://example.com") %}
  Markdown quote content.
  {% end %}
  ```

  `author`, `date`, and `url` are optional.

- Raw HTML:

  ```md
  {% raw %}
  <section>HTML left untouched by Markdown.</section>
  {% endraw %}
  ```

  The legacy `{% raw () %}...{% end %}` form is also supported.

## Local Use

Build the site:

```bash
go run ./cmd/deald build
```

Include drafts with a conventional long flag before or after the command:

```bash
go run ./cmd/deald build --drafts
go run ./cmd/deald --drafts build
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
