// Package chain executes multi-step prompt pipelines.
package chain

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/devaloi/promptkit/internal/engine"
	"github.com/devaloi/promptkit/internal/registry"
	"github.com/devaloi/promptkit/internal/validator"
)

// Step defines a single step in a prompt chain.
type Step struct {
	Template  string            `yaml:"template"`
	Vars      map[string]string `yaml:"vars"`
	OutputVar string            `yaml:"output_var"`
}

// Definition is a parsed chain YAML file.
type Definition struct {
	Name  string `yaml:"name"`
	Steps []Step `yaml:"steps"`
}

// Result holds the outputs from executing a chain.
type Result struct {
	Final         string
	Intermediates map[string]string
}

// ParseFile reads and parses a chain definition from a YAML file.
func ParseFile(path string) (Definition, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Definition{}, fmt.Errorf("reading chain file: %w", err)
	}
	return Parse(data)
}

// Parse parses a chain definition from YAML bytes.
func Parse(data []byte) (Definition, error) {
	var def Definition
	if err := yaml.Unmarshal(data, &def); err != nil {
		return Definition{}, fmt.Errorf("parsing chain YAML: %w", err)
	}
	if len(def.Steps) == 0 {
		return Definition{}, fmt.Errorf("chain %q has no steps", def.Name)
	}
	return def, nil
}

// Execute runs a chain definition against a registry, passing initial vars.
// Each step renders a template and captures output into the variable namespace.
func Execute(def Definition, reg *registry.Registry, initialVars map[string]any) (Result, error) {
	vars := make(map[string]any, len(initialVars))
	for k, v := range initialVars {
		vars[k] = v
	}

	intermediates := make(map[string]string, len(def.Steps))
	var lastOutput string

	for i, step := range def.Steps {
		tmpl, err := reg.Get(step.Template)
		if err != nil {
			return Result{}, fmt.Errorf("step %d: %w", i+1, err)
		}

		// Build step vars: resolve any template references from current var namespace.
		stepVars := make(map[string]any, len(step.Vars))
		for k, v := range step.Vars {
			stepVars[k] = resolveVar(v, vars)
		}

		// Validate required vars.
		if len(tmpl.Meta.RequiredVars) > 0 {
			if err := validator.Validate(tmpl.Meta.RequiredVars, stepVars); err != nil {
				return Result{}, fmt.Errorf("step %d (%s): %w", i+1, step.Template, err)
			}
		}

		// Render the template.
		result, err := engine.Render(tmpl.Content, stepVars, reg.Includes())
		if err != nil {
			return Result{}, fmt.Errorf("step %d (%s): rendering: %w", i+1, step.Template, err)
		}

		lastOutput = result.Output

		// Capture output into the variable namespace.
		if step.OutputVar != "" {
			vars[step.OutputVar] = result.Output
			intermediates[step.OutputVar] = result.Output
		}
	}

	return Result{Final: lastOutput, Intermediates: intermediates}, nil
}

// resolveVar resolves simple {{ .varname }} references in a string value.
func resolveVar(val string, vars map[string]any) string {
	result, err := engine.Render(val, vars, nil)
	if err != nil {
		return val
	}
	return result.Output
}
