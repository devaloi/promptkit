package frontmatter

import (
	"errors"
	"testing"
)

func TestParse_ValidFrontmatter(t *testing.T) {
	input := `---
name: summarize
description: Summarize a document
required_vars:
  - document
  - max_words
model_hint: gpt-4
---
Summarize the following document.`

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Meta.Name != "summarize" {
		t.Errorf("expected name 'summarize', got %q", result.Meta.Name)
	}
	if result.Meta.Description != "Summarize a document" {
		t.Errorf("expected description 'Summarize a document', got %q", result.Meta.Description)
	}
	if len(result.Meta.RequiredVars) != 2 {
		t.Fatalf("expected 2 required vars, got %d", len(result.Meta.RequiredVars))
	}
	if result.Meta.RequiredVars[0] != "document" || result.Meta.RequiredVars[1] != "max_words" {
		t.Errorf("unexpected required_vars: %v", result.Meta.RequiredVars)
	}
	if result.Meta.ModelHint != "gpt-4" {
		t.Errorf("expected model_hint 'gpt-4', got %q", result.Meta.ModelHint)
	}
	if result.Body != "Summarize the following document." {
		t.Errorf("unexpected body: %q", result.Body)
	}
}

func TestParse_MissingFrontmatter(t *testing.T) {
	input := "Just a plain template with no frontmatter."
	result, err := Parse(input)
	if !errors.Is(err, ErrNoFrontmatter) {
		t.Fatalf("expected ErrNoFrontmatter, got %v", err)
	}
	if result.Body != input {
		t.Errorf("expected full content in body, got %q", result.Body)
	}
}

func TestParse_EmptyRequiredVars(t *testing.T) {
	input := `---
name: simple
description: A simple prompt
required_vars: []
model_hint: gpt-3.5-turbo
---
Hello, world.`

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Meta.Name != "simple" {
		t.Errorf("expected name 'simple', got %q", result.Meta.Name)
	}
	if len(result.Meta.RequiredVars) != 0 {
		t.Errorf("expected 0 required vars, got %d", len(result.Meta.RequiredVars))
	}
	if result.Body != "Hello, world." {
		t.Errorf("unexpected body: %q", result.Body)
	}
}

func TestParse_NoClosingDelimiter(t *testing.T) {
	input := `---
name: broken
description: Missing closing delimiter`

	_, err := Parse(input)
	if !errors.Is(err, ErrNoFrontmatter) {
		t.Fatalf("expected ErrNoFrontmatter, got %v", err)
	}
}

func TestParse_EmptyBody(t *testing.T) {
	input := `---
name: empty
description: No body content
---`

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Body != "" {
		t.Errorf("expected empty body, got %q", result.Body)
	}
}
