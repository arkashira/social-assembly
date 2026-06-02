package markdown

import (
	"bytes"
	"html/template"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
)

// Renderer implements a secure markdown rendering service with sanitization and accessibility features
type Renderer struct {
	opts *Options
	policy *bluemonday.Policy
}

// Options configures the markdown renderer
type Options struct {
	SafeLinks bool // Sanitize external links
	Strict    bool // Enable strict parsing
	AllowHTML bool // Allow limited safe HTML (default false)
}

// NewRenderer creates a new markdown renderer with secure defaults
func NewRenderer(opts *Options) *Renderer {
	if opts == nil {
		opts = &Options{
			SafeLinks: true,
			Strict:    true,
			AllowHTML: false,
		}
	}

	// Create a strict policy that allows only safe HTML
	policy := bluemonday.UGCPolicy()
	policy.AllowAttrs("href").OnElements("a")
	policy.AllowAttrs("title").OnElements("a")
	policy.AllowAttrs("class", "id").OnElements("div", "span", "p", "h1", "h2", "h3", "h4", "h5", "h6")
	policy.AllowAttrs("alt").OnElements("img")

	return &Renderer{
		opts:   opts,
		policy: policy,
	}
}

// Render converts markdown text to safe HTML
func (r *Renderer) Render(mdText string) template.HTML {
	if mdText == "" {
		return ""
	}

	// Pre-sanitize markdown input
	if r.opts.Strict {
		mdText = sanitizeMarkdown(mdText)
	}

	// Parse markdown with safe extensions
	p := parser.NewWithExtensions(
		parser.CommonExtensions |
			parser.AutoHeadingIDs |
			parser.NoEmptyLineBeforeBlock,
	)

	doc := p.Parse([]byte(mdText))

	// Configure HTML renderer with security flags
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	if r.opts.SafeLinks {
		htmlFlags |= html.Safelink
	}

	renderer := html.NewRenderer(html.RendererOptions{
		Flags: htmlFlags,
		RenderNodeHook: func(w html.Buffer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
			// Enforce heading level constraints
			if heading, ok := node.(*ast.Heading); ok && entering {
				if heading.Level < 1 {
					heading.Level = 1
				}
				if heading.Level > 6 {
					heading.Level = 6
				}
			}
			return ast.GoToNext, false
		},
	})

	// Render to HTML
	htmlBytes := markdown.Render(doc, renderer)
	htmlStr := string(htmlBytes)

	// Post-process for accessibility
	htmlStr = fixLists(htmlStr)

	// Sanitize final output
	if !r.opts.AllowHTML {
		htmlStr = r.policy.Sanitize(htmlStr)
	}

	return template.HTML(htmlStr)
}

// sanitizeMarkdown removes potentially dangerous markdown constructs
func sanitizeMarkdown(md string) string {
	// Remove raw HTML blocks
	reHTML := regexp.MustCompile(`(?s)<[^>]*>`)
	md = reHTML.ReplaceAllString(md, "")

	// Remove javascript: links
	reJS := regexp.MustCompile(`javascript\s*:`)
	md = reJS.ReplaceAllString(md, "")

	// Limit heading levels to prevent abuse
	reHeadings := regexp.MustCompile(`(?m)^#{1,10}\s`)
	md = reHeadings.ReplaceAllStringFunc(md, func(s string) string {
		count := strings.Count(s, "#")
		if count > 6 {
			return strings.Repeat("#", 6) + " "
		}
		return s
	})

	return md
}

// fixLists ensures proper list structure for accessibility
func fixLists(html string) string {
	// Convert loose lists to proper structure
	reLooseList := regexp.MustCompile(`(?m)<li>\s*<p>`)
	html = reLooseList.ReplaceAllString(html, "<li>")

	reLooseListClose := regexp.MustCompile(`(?m)</p>\s*</li>`)
	html = reLooseListClose.ReplaceAllString(html, "</li>")

	return html
}

// RenderToString renders markdown to a plain string (stripped of HTML tags)
func (r *Renderer) RenderToString(mdText string) string {
	html := r.Render(mdText)
	// Strip HTML tags for plain text rendering
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(string(html), "")
}

// RenderWithExcerpt renders markdown with a maximum length excerpt
func (r *Renderer) RenderWithExcerpt(mdText string, maxLength int) (template.HTML, string) {
	html := r.Render(mdText)
	plain := r.RenderToString(mdText)

	if len(plain) <= maxLength {
		return html, plain
	}

	excerpt := plain[:maxLength]
	return html, excerpt
}