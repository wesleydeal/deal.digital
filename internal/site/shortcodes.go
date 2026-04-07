package site

import (
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"
)

var (
	inlineFitImagePattern   = regexp.MustCompile(`\{\{\s*fitimg\((.*?)\)\s*\}\}`)
	inlineWavePattern       = regexp.MustCompile(`\{\{\s*wave\((.*?)\)\s*\}\}`)
	blockquotePattern       = regexp.MustCompile(`(?s)\{%\s*blockquote\((.*?)\)\s*%\}(.*?)\{%\s*end\s*%\}`)
	rawBlockPattern         = regexp.MustCompile(`(?s)\{%\s*raw\s*\(\s*\)\s*%\}(.*?)\{%\s*end\s*%\}`)
	shortcodeArgument       = regexp.MustCompile(`([A-Za-z0-9_-]+)\s*=\s*("(?:[^"\\]|\\.)*"|'(?:[^'\\]|\\.)*'|[^,]+)`)
	standalonePlaceholderRE = regexp.MustCompile(`(?m)<p>\s*(%%SHORTCODE_[0-9]+%%)\s*</p>`)
)

type shortcodeState struct {
	values map[string]string
	nextID int
}

func newShortcodeState() *shortcodeState {
	return &shortcodeState{values: map[string]string{}}
}

func (s *shortcodeState) token(html string) string {
	token := fmt.Sprintf("%%SHORTCODE_%d%%", s.nextID)
	s.nextID++
	s.values[token] = html
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

func (r *Renderer) preprocessShortcodes(input string, depth int) string {
	if depth > 8 {
		return input
	}
	state := newShortcodeState()

	input = rawBlockPattern.ReplaceAllStringFunc(input, func(match string) string {
		body := rawBlockPattern.FindStringSubmatch(match)[1]
		return state.token(strings.TrimSpace(body))
	})

	input = blockquotePattern.ReplaceAllStringFunc(input, func(match string) string {
		submatch := blockquotePattern.FindStringSubmatch(match)
		args := parseShortcodeArgs(submatch[1])
		bodyHTML, _, plain, err := r.renderMarkdown(strings.TrimSpace(submatch[2]), depth+1)
		if err != nil {
			bodyHTML = templateEscape(strings.TrimSpace(submatch[2]))
			plain = strings.TrimSpace(submatch[2])
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
		html := `<p><a href="` + html.EscapeString(pathValue) + `"><img src="` + html.EscapeString(pathValue) + `"` + style + ` loading="lazy"></a></p>`
		return state.token(html)
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
		html := `<svg width="1200" height="4" xmlns="http://www.w3.org/2000/svg" style="width:` + html.EscapeString(width) + `"><defs><pattern id="wave-1" x="0" y="0" width="15" height="4" patternUnits="userSpaceOnUse"><path d="M0 1C3.80745 1 3.80745 3 7.6149 3C11.4223 3 11.4223 1 15.2298 1C19.0372 1 19.0372 3 22.8447" stroke-width="1" fill="none" style="stroke:` + html.EscapeString(color) + `;"></path></pattern></defs><rect x="0" y="0" width="1200" height="4" fill="url(#wave-1)"></rect></svg>`
		return state.token(html)
	})

	return state.restore(input)
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

func templateEscape(v string) string {
	return "<p>" + html.EscapeString(v) + "</p>"
}
