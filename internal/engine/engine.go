package engine

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/devaloi/promptkit/internal/frontmatter"
)

// RenderResult holds the output of rendering a template.
type RenderResult struct {
	Output string
	Meta   frontmatter.Metadata
}

// Render parses frontmatter from content, then renders the template body with
// the provided variables and optional include templates.
func Render(content string, vars map[string]any, includes map[string]string) (RenderResult, error) {
	parsed, fmErr := frontmatter.Parse(content)

	meta := parsed.Meta
	body := parsed.Body

	// If no frontmatter was found, render the entire content as a template.
	if fmErr != nil {
		body = content
		meta = frontmatter.Metadata{}
	}

	tmpl := template.New("main").Funcs(FuncMap())

	// Register include templates.
	for name, incBody := range includes {
		if _, err := tmpl.New(name).Parse(incBody); err != nil {
			return RenderResult{}, fmt.Errorf("parsing include %q: %w", name, err)
		}
	}

	if _, err := tmpl.Parse(body); err != nil {
		return RenderResult{}, fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return RenderResult{}, fmt.Errorf("executing template: %w", err)
	}

	return RenderResult{Output: buf.String(), Meta: meta}, nil
}
