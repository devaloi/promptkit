package chain

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/devaloi/promptkit/internal/registry"
)

func setupChainTest(t *testing.T) (*registry.Registry, string) {
	t.Helper()
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "step_one.tmpl"), `---
name: step_one
description: First step
required_vars:
  - input
---
Processed: {{ .input }}`)

	writeFile(t, filepath.Join(dir, "step_two.tmpl"), `---
name: step_two
description: Second step
required_vars:
  - data
---
Final: {{ .data }}`)

	reg := registry.New()
	if err := reg.LoadDir(dir); err != nil {
		t.Fatal(err)
	}

	chainYAML := `name: test-chain
steps:
  - template: step_one
    vars:
      input: "{{ .user_input }}"
    output_var: step_one_out

  - template: step_two
    vars:
      data: "{{ .step_one_out }}"
    output_var: final_out
`

	chainPath := filepath.Join(dir, "chain.yaml")
	writeFile(t, chainPath, chainYAML)

	return reg, chainPath
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestChain_TwoStepPipeline(t *testing.T) {
	reg, chainPath := setupChainTest(t)

	def, err := ParseFile(chainPath)
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}

	result, err := Execute(def, reg, map[string]any{"user_input": "hello"})
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.Final != "Final: Processed: hello" {
		t.Errorf("unexpected final output: %q", result.Final)
	}

	if result.Intermediates["step_one_out"] != "Processed: hello" {
		t.Errorf("unexpected intermediate: %q", result.Intermediates["step_one_out"])
	}
}

func TestChain_MissingTemplate(t *testing.T) {
	reg := registry.New()
	def := Definition{
		Name: "bad-chain",
		Steps: []Step{
			{Template: "nonexistent", Vars: map[string]string{"x": "y"}, OutputVar: "out"},
		},
	}

	_, err := Execute(def, reg, nil)
	if err == nil {
		t.Fatal("expected error for missing template in chain")
	}
}

func TestChain_ParseEmptySteps(t *testing.T) {
	data := []byte(`name: empty
steps: []
`)
	_, err := Parse(data)
	if err == nil {
		t.Fatal("expected error for empty steps")
	}
}

func TestChain_VariablePassingBetweenSteps(t *testing.T) {
	reg, chainPath := setupChainTest(t)

	def, err := ParseFile(chainPath)
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}

	result, err := Execute(def, reg, map[string]any{"user_input": "world"})
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if _, ok := result.Intermediates["step_one_out"]; !ok {
		t.Error("expected step_one_out in intermediates")
	}
	if _, ok := result.Intermediates["final_out"]; !ok {
		t.Error("expected final_out in intermediates")
	}
}
