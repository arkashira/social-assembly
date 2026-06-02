package markdown

import (
	"strings"
	"testing"
)

func TestRenderer_Render(t *testing.T) {
	r := NewRenderer(RenderOptions{SafeLinks: true, Strict: true})

	tests := []struct {
		name        string
		input       string
		contains    []string
		notContains []string
	}{
		{
			name:   "empty string",
			input:  "",
		},
		{
			name:  "header levels",
			input: "# H1\n## H2\n### H3",
			contains: []string{
				"<h1", "<h2", "<h3",
			},
		},
		{
			name:  "unordered list",
			input: "- item 1\n- item 2",
			contains: []string{
				"<ul>", "<li>item 1", "<li>item 2",
			},
		},
		{
			name:  "ordered list",
			input: "1. first\n2. second",
			contains: []string{
				"<ol>", "<li>first", "<li>second",
			},
		},
		{
			name:  "code block",
			input: "