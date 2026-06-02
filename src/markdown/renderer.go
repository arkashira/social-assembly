package markdown

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// Renderer provides markdown rendering capabilities for posts and comments.
type Renderer struct {
	opts *Options
}

// Options configures the markdown renderer.
type Options struct {
	SafeLinks bool
}

// NewRenderer creates a new markdown renderer with the given options.
func NewRenderer(opts *Options) *Renderer {
	if opts == nil {
		opts = &Options{SafeLinks: true}
	}
	return &Renderer{opts: opts}
}

// Render converts markdown source into safe HTML suitable for web display.
func (r *Renderer) Render(src []byte) ([]byte, error) {
	// Create parser with extensions
	p := parser.NewWithExtensions(
		parser.CommonExtensions |
			parser.AutoHeadingIDs |
			parser.NoEmptyLineBeforeBlock,
	)

	// Parse markdown into AST
	doc := p.Parse(src)

	// Configure HTML renderer
	renderer := html.NewRenderer(html.RendererOptions{
		Flags: html.CommonFlags | html.HrefTargetBlank,
		CSS:   "markdown-body",
	})

	// Render AST to HTML
	htmlBytes := markdown.Render(doc, renderer)

	// Sanitize output for safety
	return sanitizeHTML(htmlBytes), nil
}

// sanitizeHTML strips dangerous tags and attributes from rendered HTML.
func sanitizeHTML(b []byte) []byte {
	s := string(b)
	// Remove script tags and other dangerous elements
	s = strings.ReplaceAll(s, "<script", "&lt;script")
	s = strings.ReplaceAll(s, "</script>", "&lt;/script&gt;")
	// Remove on* event handlers
	s = strings.ReplaceAll(s, " on", " data-on")
	return []byte(s)
}

// RenderToTemplate renders markdown and wraps it in a template-safe structure.
func (r *Renderer) RenderToTemplate(src []byte) (template.HTML, error) {
	htmlBytes, err := r.Render(src)
	if err != nil {
		return "", err
	}
	return template.HTML(htmlBytes), nil
}