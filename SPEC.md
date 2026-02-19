# A06: promptkit — LLM Prompt Template Engine in Go

**Catalog ID:** A06 | **Size:** S | **Language:** Go
**Repo name:** `promptkit`
**One-liner:** A template engine for LLM prompts with variable injection, conditionals, includes, validation, a registry of reusable templates, prompt chaining, and a CLI.

---

## Why This Stands Out

- **Go text/template engine extended for LLM workflows** — not a generic template library, purpose-built for prompt engineering with LLM-specific helpers
- **Template validation** — validates required variables are present before rendering, catches errors at build time not runtime
- **Template registry** — loads templates from a directory, lookups by name, supports template includes for composable system/user/assistant blocks
- **Built-in LLM helpers** — `truncate`, `json_encode`, `word_count`, `token_estimate`, `upper`, `lower`, `join` — registered as template functions
- **Prompt chaining** — define multi-step pipelines where the output of one template render feeds as input to the next
- **YAML frontmatter** — templates carry metadata (name, description, required variables, model hints) in YAML frontmatter, parsed and validated
- **CLI tool** — `promptkit render template.tmpl --var key=value` for quick iteration without writing Go code

---

## Architecture

```
promptkit/
├── cmd/
│   └── promptkit/
│       └── main.go              # CLI: render, validate, list, chain
├── internal/
│   ├── config/
│   │   └── config.go            # Template directory path, default settings
│   ├── engine/
│   │   ├── engine.go            # Core render engine: parse frontmatter, inject vars, execute template
│   │   ├── engine_test.go
│   │   ├── funcmap.go           # LLM helper functions registered to template.FuncMap
│   │   └── funcmap_test.go
│   ├── frontmatter/
│   │   ├── frontmatter.go       # Parse YAML frontmatter from template files
│   │   └── frontmatter_test.go
│   ├── registry/
│   │   ├── registry.go          # Load templates from directory, lookup by name
│   │   └── registry_test.go
│   ├── validator/
│   │   ├── validator.go         # Check required vars present, type hints
│   │   └── validator_test.go
│   └── chain/
│       ├── chain.go             # Prompt chaining: pipeline of template renders
│       └── chain_test.go
├── templates/
│   ├── summarize.tmpl           # Example: summarization prompt
│   ├── classify.tmpl            # Example: classification prompt
│   ├── extract.tmpl             # Example: entity extraction prompt
│   ├── chain_example.yaml       # Example: chain definition
│   └── includes/
│       ├── system_default.tmpl  # Reusable system message block
│       └── json_format.tmpl     # Reusable JSON output instruction
├── go.mod
├── go.sum
├── Makefile
├── .gitignore
├── .golangci.yml
├── LICENSE
└── README.md
```

---

## Template Format

```
---
name: summarize
description: Summarize a document with configurable length
required_vars:
  - document
  - max_words
model_hint: gpt-4
---
{{ template "system_default" }}

Summarize the following document in {{ .max_words }} words or fewer.

{{ .document | truncate 8000 }}

Respond in JSON:
{{ template "json_format" }}
```

---

## Helper Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `truncate` | `truncate <max_chars> <text>` | Truncate text to max characters with ellipsis |
| `json_encode` | `json_encode <value>` | Marshal value to JSON string |
| `word_count` | `word_count <text>` | Count words in text |
| `token_estimate` | `token_estimate <text>` | Estimate token count (~4 chars/token) |
| `upper` | `upper <text>` | Uppercase |
| `lower` | `lower <text>` | Lowercase |
| `join` | `join <sep> <slice>` | Join slice elements with separator |
| `default` | `default <fallback> <value>` | Return fallback if value is empty |

---

## Chain Definition (YAML)

```yaml
name: summarize-and-classify
steps:
  - template: summarize
    vars:
      document: "{{ .input_document }}"
      max_words: "100"
    output_var: summary

  - template: classify
    vars:
      text: "{{ .summary }}"
      categories: "{{ .categories }}"
    output_var: classification
```

---

## Tech Stack

| Component | Choice |
|-----------|--------|
| Language | Go 1.26 |
| Templates | `text/template` (stdlib) |
| Frontmatter | `gopkg.in/yaml.v3` |
| CLI | `cobra` + `pflag` |
| Testing | stdlib |
| Linting | golangci-lint |

---

## Phased Build Plan

### Phase 1: Foundation

**1.1 — Project setup**
- `go mod init github.com/devaloi/promptkit`
- Directory structure, Makefile, `.gitignore`, `.golangci.yml`

**1.2 — Frontmatter parser**
- Parse YAML frontmatter delimited by `---` from template file content
- Return metadata struct (Name, Description, RequiredVars, ModelHint) + template body
- Tests: valid frontmatter, missing frontmatter, empty required_vars

**1.3 — Helper functions**
- Implement all functions in funcmap table
- Register as `template.FuncMap`
- Tests: each function with normal input, edge cases (empty string, zero, nil)

### Phase 2: Engine + Validation

**2.1 — Render engine**
- Parse template body with `text/template` + custom FuncMap
- Inject variables as `map[string]any`
- Return rendered string
- Support `{{ template "name" }}` includes via `template.ParseFiles` or manual registration
- Tests: simple render, conditionals, includes, missing variable behavior

**2.2 — Validator**
- Given frontmatter `required_vars` and provided variables map, check all required vars are present
- Return list of missing variables as structured error
- Tests: all present → pass, one missing → error with name, multiple missing → all listed

### Phase 3: Registry + Chaining

**3.1 — Template registry**
- Load all `.tmpl` files from a directory recursively
- Parse frontmatter for each, index by `name` field
- Lookup by name, list all available templates
- Includes directory loaded separately and available to all templates
- Tests: load directory with 3 templates, lookup by name, missing name → error

**3.2 — Prompt chaining**
- Parse chain YAML definition (list of steps with template name, vars, output_var)
- Execute steps sequentially: render template, capture output, inject into next step's vars
- Return final output + all intermediate outputs
- Tests: two-step chain, variable passing between steps, missing template in chain → error

### Phase 4: CLI + Polish

**4.1 — CLI commands**
- `promptkit render <template> --var key=value [--var key2=value2] --dir ./templates` — render a single template
- `promptkit validate <template> --dir ./templates` — validate required vars against provided vars
- `promptkit list --dir ./templates` — list all templates with name and description
- `promptkit chain <chain.yaml> --var key=value --dir ./templates` — run a prompt chain

**4.2 — Example templates**
- `summarize.tmpl` — summarization with configurable word count
- `classify.tmpl` — text classification with category list
- `extract.tmpl` — entity extraction with output format
- `chain_example.yaml` — summarize → classify pipeline

**4.3 — README**
- Badges, install, quick start
- Template format with frontmatter example
- Helper function reference table
- Chain definition format
- CLI usage examples

---

## Commit Plan

1. `chore: scaffold project with directory structure`
2. `feat: add YAML frontmatter parser`
3. `feat: add LLM helper functions (truncate, json_encode, token_estimate, etc.)`
4. `feat: add template render engine with variable injection`
5. `feat: add template validation for required variables`
6. `feat: add template registry with directory loading`
7. `feat: add prompt chaining with sequential pipeline execution`
8. `feat: add CLI with render, validate, list, and chain commands`
9. `feat: add example templates and chain definition`
10. `docs: add README with usage examples and helper reference`
