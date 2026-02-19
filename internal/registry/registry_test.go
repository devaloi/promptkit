package registry

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	// Create includes directory.
	incDir := filepath.Join(dir, "includes")
	if err := os.Mkdir(incDir, 0o755); err != nil {
		t.Fatal(err)
	}

	writeFile(t, filepath.Join(incDir, "header.tmpl"), "=== HEADER ===")

	writeFile(t, filepath.Join(dir, "greet.tmpl"), `---
name: greet
description: A greeting
required_vars:
  - name
---
Hello, {{ .name }}!`)

	writeFile(t, filepath.Join(dir, "farewell.tmpl"), `---
name: farewell
description: A farewell
required_vars:
  - name
---
Goodbye, {{ .name }}!`)

	writeFile(t, filepath.Join(dir, "plain.tmpl"), `Just plain text.`)

	return dir
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestRegistry_LoadDir(t *testing.T) {
	dir := setupTestDir(t)
	reg := New()

	if err := reg.LoadDir(dir); err != nil {
		t.Fatalf("LoadDir error: %v", err)
	}

	if len(reg.List()) != 3 {
		t.Errorf("expected 3 templates, got %d", len(reg.List()))
	}
}

func TestRegistry_GetByName(t *testing.T) {
	dir := setupTestDir(t)
	reg := New()

	if err := reg.LoadDir(dir); err != nil {
		t.Fatalf("LoadDir error: %v", err)
	}

	tmpl, err := reg.Get("greet")
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if tmpl.Meta.Description != "A greeting" {
		t.Errorf("expected description 'A greeting', got %q", tmpl.Meta.Description)
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	reg := New()
	_, err := reg.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing template")
	}
}

func TestRegistry_Includes(t *testing.T) {
	dir := setupTestDir(t)
	reg := New()

	if err := reg.LoadDir(dir); err != nil {
		t.Fatalf("LoadDir error: %v", err)
	}

	includes := reg.Includes()
	if _, ok := includes["header"]; !ok {
		t.Error("expected 'header' include to be loaded")
	}
}

func TestRegistry_PlainTemplate(t *testing.T) {
	dir := setupTestDir(t)
	reg := New()

	if err := reg.LoadDir(dir); err != nil {
		t.Fatalf("LoadDir error: %v", err)
	}

	tmpl, err := reg.Get("plain")
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if tmpl.Meta.Name != "" {
		t.Errorf("expected empty meta name for plain template, got %q", tmpl.Meta.Name)
	}
}
