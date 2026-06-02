package markdown

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderer_Render(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		wantHTML string
	}{
		{
			name:     "empty",
			markdown: "",
			wantHTML: "",
		},
		{
			name:     "header",
			markdown: "# Hello World",
			wantHTML: "<h1>Hello World</h1>\n",
		},
		{
			name:     "list",
			markdown: "- item 1\n- item 2",
			wantHTML: "<ul>\n<li>item 1</li>\n<li>item 2</li>\n</ul>\n",
		},
		{
			name:     "code_block",
			markdown: "