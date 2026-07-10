package render

import (
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"

	"deal.digital/internal/content"
)

var (
	inlineFitImagePattern   = regexp.MustCompile(`\{\{\s*fitimg\((.*?)\)\s*\}\}`)
	inlineWavePattern       = regexp.MustCompile(`\{\{\s*wave\((.*?)\)\s*\}\}`)
	blockquotePattern       = regexp.MustCompile(`(?s)\{%\s*blockquote\((.*?)\)\s*%\}(.*?)\{%\s*end\s*%\}`)
	rawBlockPattern         = regexp.MustCompile(`(?s)\{%\s*raw(?:\s*\(\s*\))?\s*%\}(.*?)\{%\s*(?:end|endraw)\s*%\}`)
	variantBlockPattern     = regexp.MustCompile(`(?s)\{%\s*variant\((.*?)\)\s*%\}(.*?)\{%\s*endvariant\s*%\}`)
	variantMenuPattern      = regexp.MustCompile(`\{\{\s*variant_menu\(\s*\)\s*\}\}`)
	shortcodeArgument       = regexp.MustCompile(`([A-Za-z0-9_-]+)\s*=\s*("(?:[^"\\]|\\.)*"|'(?:[^'\\]|\\.)*'|[^,]+)`)
	standalonePlaceholderRE = regexp.MustCompile(`(?m)<p>\s*(%%SHORTCODE_[0-9]+%%)\s*</p>`)
)

type shortcodeState struct {
	values map[string]string
	nextID int
}

type variantContext struct {
	current  string
	variants []content.PageVariant
}

func newShortcodeState() *shortcodeState {
	return &shortcodeState{values: map[string]string{}}
}

func (s *shortcodeState) token(fragment string) string {
	token := fmt.Sprintf("%%SHORTCODE_%d%%", s.nextID)
	s.nextID++
	s.values[token] = fragment
	return token
}

func (s *shortcodeState) restore(input string) string {
	input = standalonePlaceholderRE.ReplaceAllStringFunc(input, func(match string) string {
		token := standalonePlaceholderRE.FindStringSubmatch(match)[1]
		return s.values[token]
	})
	for token, value := range s.values {
		input = strings.ReplaceAll(input, token, value)
	}
	return input
}

func (r *Renderer) preprocessShortcodes(input string, depth int, variant *variantContext) (string, error) {
	if depth > 8 {
		return input, nil
	}
	state := newShortcodeState()

	input = rawBlockPattern.ReplaceAllStringFunc(input, func(match string) string {
		body := rawBlockPattern.FindStringSubmatch(match)[1]
		return state.token(trimRawBlockBody(body))
	})
	var err error
	input, err = filterVariantBlocks(input, variant)
	if err != nil {
		return "", err
	}
	if variantMenuPattern.MatchString(input) {
		if variant == nil {
			return "", fmt.Errorf("variant_menu shortcode requires page variants")
		}
		input = variantMenuPattern.ReplaceAllStringFunc(input, func(string) string {
			return state.token(variantMenuHTML(variant))
		})
	}

	var blockquoteErr error
	input = blockquotePattern.ReplaceAllStringFunc(input, func(match string) string {
		submatch := blockquotePattern.FindStringSubmatch(match)
		args := parseShortcodeArgs(submatch[1])
		bodyHTML, _, plain, err := r.renderMarkdownWithVariant(strings.TrimSpace(submatch[2]), depth+1, variant)
		if err != nil {
			blockquoteErr = err
			return match
		}

		var b strings.Builder
		b.WriteString(`<blockquote class="shortcode-blockquote">`)
		b.WriteString(bodyHTML)
		if author := args["author"]; author != "" || args["date"] != "" {
			b.WriteString(`<footer>`)
			if author != "" {
				label := "—" + html.EscapeString(author)
				if url := args["url"]; url != "" {
					b.WriteString(`<p class="author"><a href="`)
					b.WriteString(html.EscapeString(url))
					b.WriteString(`">`)
					b.WriteString(label)
					b.WriteString(`</a></p>`)
				} else {
					b.WriteString(`<p class="author">`)
					b.WriteString(label)
					b.WriteString(`</p>`)
				}
			}
			if date := args["date"]; date != "" {
				b.WriteString(`<p class="date"><time datetime="`)
				b.WriteString(html.EscapeString(date))
				b.WriteString(`">`)
				b.WriteString(html.EscapeString(date))
				b.WriteString(`</time></p>`)
			}
			b.WriteString(`</footer>`)
		}
		b.WriteString(`</blockquote>`)
		if plain != "" {
			b.WriteString("\n")
		}
		return state.token(b.String())
	})
	if blockquoteErr != nil {
		return "", blockquoteErr
	}

	input = inlineFitImagePattern.ReplaceAllStringFunc(input, func(match string) string {
		args := parseShortcodeArgs(inlineFitImagePattern.FindStringSubmatch(match)[1])
		pathValue := args["path"]
		if pathValue == "" {
			return match
		}
		var styleParts []string
		if width := cssDimension(args["width"]); width != "" {
			styleParts = append(styleParts, "width:"+width)
		}
		if height := cssDimension(args["height"]); height != "" {
			styleParts = append(styleParts, "height:"+height)
		}
		style := ""
		if len(styleParts) > 0 {
			style = ` style="` + strings.Join(styleParts, ";") + `"`
		}
		fragment := `<p><a href="` + html.EscapeString(pathValue) + `"><img src="` + html.EscapeString(pathValue) + `"` + style + ` loading="lazy"></a></p>`
		return state.token(fragment)
	})

	input = inlineWavePattern.ReplaceAllStringFunc(input, func(match string) string {
		args := parseShortcodeArgs(inlineWavePattern.FindStringSubmatch(match)[1])
		width := args["width"]
		if width == "" {
			width = "100%"
		}
		color := args["color"]
		if color == "" {
			color = "var(--color-fg)"
		}
		fragment := `<svg width="1200" height="4" xmlns="http://www.w3.org/2000/svg" style="width:` + html.EscapeString(width) + `"><defs><pattern id="wave-1" x="0" y="0" width="15" height="4" patternUnits="userSpaceOnUse"><path d="M0 1C3.80745 1 3.80745 3 7.6149 3C11.4223 3 11.4223 1 15.2298 1C19.0372 1 19.0372 3 22.8447" stroke-width="1" fill="none" style="stroke:` + html.EscapeString(color) + `;"></path></pattern></defs><rect x="0" y="0" width="1200" height="4" fill="url(#wave-1)"></rect></svg>`
		return state.token(fragment)
	})

	return state.restore(input), nil
}

func filterVariantBlocks(input string, variant *variantContext) (string, error) {
	for variantBlockPattern.MatchString(input) {
		var blockErr error
		input = variantBlockPattern.ReplaceAllStringFunc(input, func(match string) string {
			submatch := variantBlockPattern.FindStringSubmatch(match)
			visible, err := variantVisible(submatch[1], variant)
			if err != nil {
				blockErr = err
				return match
			}
			if visible {
				return submatch[2]
			}
			return ""
		})
		if blockErr != nil {
			return "", blockErr
		}
	}
	if strings.Contains(input, "{% variant") || strings.Contains(input, "{% endvariant") {
		return "", fmt.Errorf("variant shortcode must use {%% variant(...) %%} and {%% endvariant %%}")
	}
	return input, nil
}

func variantVisible(raw string, variant *variantContext) (bool, error) {
	if variant == nil {
		return false, fmt.Errorf("variant shortcode requires page variants")
	}
	args := parseShortcodeArgs(raw)
	include, hasInclude := args["include"]
	exclude, hasExclude := args["exclude"]
	if len(args) != 1 || hasInclude == hasExclude {
		return false, fmt.Errorf("variant shortcode requires exactly one include or exclude argument")
	}
	targets, err := variantTargets(include)
	if hasExclude {
		targets, err = variantTargets(exclude)
	}
	if err != nil {
		return false, err
	}
	for _, target := range targets {
		if !variantExists(target, variant.variants) {
			return false, fmt.Errorf("variant shortcode references unknown variant %q", target)
		}
	}
	selected := false
	for _, target := range targets {
		if target == variant.current {
			selected = true
			break
		}
	}
	if hasInclude {
		return selected, nil
	}
	return !selected, nil
}

func variantTargets(raw string) ([]string, error) {
	var targets []string
	for _, part := range strings.Split(raw, ",") {
		target := strings.TrimSpace(part)
		if target == "" {
			return nil, fmt.Errorf("variant shortcode requires at least one variant id")
		}
		targets = append(targets, target)
	}
	return targets, nil
}

func variantExists(id string, variants []content.PageVariant) bool {
	for _, variant := range variants {
		if variant.ID == id {
			return true
		}
	}
	return false
}

func variantMenuHTML(variant *variantContext) string {
	var b strings.Builder
	b.WriteString(`<menu class="page-variants">`)
	for _, item := range variant.variants {
		b.WriteString(`<li>`)
		if item.ID == variant.current {
			b.WriteString(`<span aria-current="page">`)
			b.WriteString(html.EscapeString(item.Label))
			b.WriteString(`</span>`)
		} else {
			b.WriteString(`<a href="`)
			b.WriteString(html.EscapeString(item.Route))
			b.WriteString(`">`)
			b.WriteString(html.EscapeString(item.Label))
			b.WriteString(`</a>`)
		}
		b.WriteString(`</li>`)
	}
	b.WriteString(`</menu>`)
	return b.String()
}

func parseShortcodeArgs(raw string) map[string]string {
	args := map[string]string{}
	for _, part := range shortcodeArgument.FindAllStringSubmatch(raw, -1) {
		value := strings.TrimSpace(part[2])
		value = strings.TrimPrefix(value, `"`)
		value = strings.TrimSuffix(value, `"`)
		value = strings.TrimPrefix(value, `'`)
		value = strings.TrimSuffix(value, `'`)
		args[part[1]] = value
	}
	return args
}

func cssDimension(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if _, err := strconv.Atoi(raw); err == nil {
		return raw + "px"
	}
	return raw
}

func trimRawBlockBody(body string) string {
	body = strings.TrimPrefix(body, "\r\n")
	body = strings.TrimPrefix(body, "\n")
	body = strings.TrimSuffix(body, "\r\n")
	body = strings.TrimSuffix(body, "\n")
	return body
}
