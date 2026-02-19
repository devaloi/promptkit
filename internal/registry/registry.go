// Package registry loads and indexes template files from a directory.
package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/devaloi/promptkit/internal/config"
	"github.com/devaloi/promptkit/internal/frontmatter"
)

// Template holds a parsed template file with its metadata and raw content.
type Template struct {
	Name    string
	Meta    frontmatter.Metadata
	Body    string
	Content string
}

// Registry holds loaded templates indexed by name.
type Registry struct {
	templates map[string]*Template
	includes  map[string]string
}

// New creates an empty Registry.
func New() *Registry {
	return &Registry{
		templates: make(map[string]*Template),
		includes:  make(map[string]string),
	}
}

// LoadDir loads all .tmpl files from dir and its includes/ subdirectory.
func (r *Registry) LoadDir(dir string) error {
	includesDir := filepath.Join(dir, config.IncludesDir)

	// Load includes first.
	if info, err := os.Stat(includesDir); err == nil && info.IsDir() {
		if err := r.loadIncludes(includesDir); err != nil {
			return fmt.Errorf("loading includes: %w", err)
		}
	}

	// Load top-level templates.
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading directory %q: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".tmpl") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		if err := r.loadTemplate(path); err != nil {
			return fmt.Errorf("loading template %q: %w", path, err)
		}
	}

	return nil
}

func (r *Registry) loadTemplate(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	content := string(data)
	parsed, fmErr := frontmatter.Parse(content)

	name := parsed.Meta.Name
	if name == "" {
		// Use filename without extension as fallback name.
		name = strings.TrimSuffix(filepath.Base(path), ".tmpl")
	}

	tmpl := &Template{
		Name:    name,
		Content: content,
		Body:    parsed.Body,
	}

	if fmErr == nil {
		tmpl.Meta = parsed.Meta
	}

	r.templates[name] = tmpl
	return nil
}

func (r *Registry) loadIncludes(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".tmpl") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		name := strings.TrimSuffix(entry.Name(), ".tmpl")
		r.includes[name] = string(data)
	}

	return nil
}

// Get retrieves a template by name.
func (r *Registry) Get(name string) (*Template, error) {
	tmpl, ok := r.templates[name]
	if !ok {
		return nil, fmt.Errorf("template %q not found", name)
	}
	return tmpl, nil
}

// List returns all loaded templates.
func (r *Registry) List() []*Template {
	result := make([]*Template, 0, len(r.templates))
	for _, tmpl := range r.templates {
		result = append(result, tmpl)
	}
	return result
}

// Includes returns the loaded include templates.
func (r *Registry) Includes() map[string]string {
	return r.includes
}
