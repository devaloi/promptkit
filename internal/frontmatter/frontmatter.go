// Package frontmatter parses YAML frontmatter from template files.
package frontmatter

import (
	"errors"
	"strings"

	"gopkg.in/yaml.v3"
)

// Metadata holds parsed YAML frontmatter fields from a template file.
type Metadata struct {
	Name         string   `yaml:"name"`
	Description  string   `yaml:"description"`
	RequiredVars []string `yaml:"required_vars"`
	ModelHint    string   `yaml:"model_hint"`
}

// Result contains parsed frontmatter metadata and the remaining template body.
type Result struct {
	Meta Metadata
	Body string
}

// ErrNoFrontmatter indicates the template has no YAML frontmatter delimiters.
var ErrNoFrontmatter = errors.New("no frontmatter found")

const delimiter = "---"

// Parse splits a template string into YAML frontmatter metadata and a body.
// Frontmatter must be delimited by lines containing only "---".
// If no frontmatter is present, ErrNoFrontmatter is returned and the full
// content is placed in Result.Body.
func Parse(content string) (Result, error) {
	trimmed := strings.TrimSpace(content)
	if !strings.HasPrefix(trimmed, delimiter) {
		return Result{Body: content}, ErrNoFrontmatter
	}

	// Find the closing delimiter after the opening one.
	rest := trimmed[len(delimiter):]
	idx := strings.Index(rest, "\n"+delimiter)
	if idx < 0 {
		return Result{Body: content}, ErrNoFrontmatter
	}

	rawYAML := rest[:idx]
	body := rest[idx+len("\n"+delimiter):]
	body = strings.TrimLeft(body, "\r\n")

	var meta Metadata
	if err := yaml.Unmarshal([]byte(rawYAML), &meta); err != nil {
		return Result{}, err
	}

	return Result{Meta: meta, Body: body}, nil
}
