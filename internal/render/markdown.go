package render

import (
	"bytes"
	"fmt"
	stdhtml "html"
	"net/url"
	"regexp"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	gmhtml "github.com/yuin/goldmark/renderer/html"
)

var tooltipPattern = regexp.MustCompile(`\{\{([^|{}]+)\|([^{}|]+)\}\}`)

type MarkdownRenderer struct {
	engine       goldmark.Markdown
	highlightCSS string
}

func NewMarkdownRenderer() *MarkdownRenderer {
	style := styles.Get("github")
	formatter := chromahtml.New(chromahtml.WithClasses(true))
	var css bytes.Buffer
	if style != nil {
		_ = formatter.WriteCSS(&css, style)
	}

	engine := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
				highlighting.WithFormatOptions(chromahtml.WithClasses(true)),
			),
		),
		goldmark.WithRendererOptions(gmhtml.WithUnsafe()),
	)

	return &MarkdownRenderer{
		engine:       engine,
		highlightCSS: css.String(),
	}
}

func (r *MarkdownRenderer) Render(markdown string) (string, error) {
	prepared := injectTooltipMarkup(markdown)
	var out bytes.Buffer
	if err := r.engine.Convert([]byte(prepared), &out); err != nil {
		return "", err
	}
	return out.String(), nil
}

func (r *MarkdownRenderer) HighlightCSS() string {
	return r.highlightCSS
}

func PlainText(markdown string) string {
	text := tooltipPattern.ReplaceAllString(markdown, "$1")
	var lines []string
	inFence := false
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		lines = append(lines, trimmed)
	}
	joined := strings.Join(lines, " ")
	replacer := strings.NewReplacer("#", " ", "*", " ", "_", " ", "`", " ", ">", " ", "[", " ", "]", " ", "(", " ", ")", " ", "-", " ")
	joined = replacer.Replace(joined)
	return strings.Join(strings.Fields(joined), " ")
}

func injectTooltipMarkup(markdown string) string {
	lines := strings.Split(markdown, "\n")
	inFence := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		lines[i] = tooltipPattern.ReplaceAllStringFunc(line, func(match string) string {
			parts := tooltipPattern.FindStringSubmatch(match)
			if len(parts) != 3 {
				return match
			}
			label := stdhtml.EscapeString(strings.TrimSpace(parts[1]))
			slug := strings.TrimSpace(parts[2])
			escapedSlug := stdhtml.EscapeString(url.PathEscape(slug))
			return fmt.Sprintf(`<button type="button" class="tooltip-term" data-tooltip-slug="%s"><span class="tooltip-label">%s</span><span class="tooltip-bubble" hidden></span></button>`, escapedSlug, label)
		})
	}
	return strings.Join(lines, "\n")
}
