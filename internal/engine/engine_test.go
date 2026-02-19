package engine

import (
	"strings"
	"testing"
)

func TestRender_SimpleTemplate(t *testing.T) {
	content := `---
name: greet
description: A greeting prompt
required_vars:
  - name
---
Hello, {{ .name }}!`

	result, err := Render(content, map[string]any{"name": "World"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Output != "Hello, World!" {
		t.Errorf("expected 'Hello, World!', got %q", result.Output)
	}
	if result.Meta.Name != "greet" {
		t.Errorf("expected meta name 'greet', got %q", result.Meta.Name)
	}
}

func TestRender_WithHelpers(t *testing.T) {
	content := `---
name: helpers
description: Test helper functions
required_vars:
  - text
---
{{ .text | upper }} has {{ .text | word_count }} words`

	result, err := Render(content, map[string]any{"text": "hello world"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "HELLO WORLD has 2 words"
	if result.Output != expected {
		t.Errorf("expected %q, got %q", expected, result.Output)
	}
}

func TestRender_WithIncludes(t *testing.T) {
	content := `---
name: with-include
description: Test includes
required_vars: []
---
{{ template "header" }}
Body content here.`

	includes := map[string]string{
		"header": "=== HEADER ===",
	}

	result, err := Render(content, nil, includes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result.Output, "=== HEADER ===") {
		t.Errorf("expected header include in output, got %q", result.Output)
	}
	if !strings.Contains(result.Output, "Body content here.") {
		t.Errorf("expected body in output, got %q", result.Output)
	}
}

func TestRender_NoFrontmatter(t *testing.T) {
	content := "Hello, {{ .name }}!"

	result, err := Render(content, map[string]any{"name": "Test"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Output != "Hello, Test!" {
		t.Errorf("expected 'Hello, Test!', got %q", result.Output)
	}
	if result.Meta.Name != "" {
		t.Errorf("expected empty meta name, got %q", result.Meta.Name)
	}
}

func TestRender_Conditional(t *testing.T) {
	content := `{{ if .verbose }}Detailed output{{ else }}Brief output{{ end }}`

	result, err := Render(content, map[string]any{"verbose": true}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Output != "Detailed output" {
		t.Errorf("expected 'Detailed output', got %q", result.Output)
	}
}

func TestRender_InvalidTemplate(t *testing.T) {
	content := "{{ .foo | nonexistent }}"
	_, err := Render(content, map[string]any{"foo": "bar"}, nil)
	if err == nil {
		t.Fatal("expected error for invalid template function")
	}
}
