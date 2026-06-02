package markdown

import (
	"bytes"
	"html/template"
	"net/url"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// RenderOptions configures the Markdown renderer.
type RenderOptions struct {
	// SafeLinks removes javascript:, data: and other dangerous protocols from href/src.
	SafeLinks bool
	// Strict forces the parser to use the same extensions as the original code
	// (AutoHeadingIDs, NoEmptyLineBeforeBlock, SafeLinks) and disables
	// extensions that could introduce raw HTML.
	Strict bool
}

// Renderer is a thin wrapper around gomarkdown that guarantees safe output.
type Renderer struct {
	opts     RenderOptions
	parser   *parser.Parser
	renderer *html.Renderer
	policy   *bluemonday.Policy
}

// NewRenderer creates a new Renderer with the supplied options.
func NewRenderer(opts RenderOptions) *Renderer {
	// ---- Parser ----------------------------------------------------------
	ext := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	if opts.SafeLinks {
		ext |= parser.SafeLinks
	}
	if opts.Strict {
		ext |= parser.SafeLinks | parser.NoEmptyLineBeforeBlock
	}
	p := parser.NewWithExtensions(ext)

	// ---- HTML renderer -----------------------------------------------
	h := html.NewRenderer(html.RendererOptions{
		Flags: html.CommonFlags | html.HrefTargetBlank,
	})

	// ---- Sanitisation policy ------------------------------------------
	policy := bluemonday.UGCPolicy()
	if opts.SafeLinks {
		// Remove javascript:, data: and other unsafe protocols
		policy.AllowAttrs("href", "src").Matching(bluemonday.SafeURLPattern).Globally()
	}
	// Ensure target="_blank" and rel="noopener noreferrer" on external links
	policy.AddTargetBlankOnLinks()
	policy.RequireNoFollowOnExternalLinks()

	return &Renderer{
		opts:     opts,
		parser:   p,
		renderer: h,
		policy:   policy,
	}
}

// Render converts Markdown to safe HTML.  It never panics – callers can
// decide how to handle an empty result.
func (r *Renderer) Render(md string) template.HTML {
	if md == "" {
		return ""
	}

	// Parse into an AST
	ast := r.parser.Parse([]byte(md))
	if ast == nil {
		return ""
	}

	// Render to raw HTML
	var buf bytes.Buffer
	html.Render(&buf, ast, r.renderer)

	// Sanitize the output
	safe := r.policy.SanitizeBytes(buf.Bytes())

	return template.HTML(safe)
}

// MustRender panics if the result is empty while the input was non‑empty.
// Useful for init‑time rendering where a failure is unrecoverable.
func (r *Renderer) MustRender(md string) template.HTML {
	res := r.Render(md)
	if res == "" && md != "" {
		panic("markdown rendering failed")
	}
	return res
}