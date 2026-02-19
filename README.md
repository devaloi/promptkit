# promptkit

[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Lint](https://img.shields.io/badge/lint-golangci--lint-blue)](https://golangci-lint.run)

A template engine for LLM prompts with variable injection, validation, includes, a registry of reusable templates, prompt chaining, and a CLI.

Built on Go's `text/template` — purpose-built for prompt engineering workflows.

## Features

- **YAML frontmatter** — templates carry metadata (name, description, required variables, model hints)
- **Variable injection** — render templates with dynamic values via `map[string]any`
- **Template includes** — compose prompts from reusable blocks (`{{ template "name" }}`)
- **Validation** — check required variables before rendering, catch errors at build time
- **Template registry** — load templates from a directory, look up by name
- **Prompt chaining** — define multi-step pipelines where output feeds into the next step
- **LLM helper functions** — `truncate`, `json_encode`, `word_count`, `token_estimate`, and more
- **CLI tool** — render, validate, list, and chain prompts from the terminal

## Installation

```bash
go install github.com/devaloi/promptkit/cmd/promptkit@latest
```

Or clone and build:

```bash
git clone https://github.com/devaloi/promptkit.git
cd promptkit
make build
```

## Prerequisites

- Go 1.22+

## Template Format

Templates use YAML frontmatter between `---` delimiters:

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

{{ template "json_format" }}
```

### Frontmatter Fields

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Template identifier for registry lookup |
| `description` | string | Human-readable description |
| `required_vars` | list | Variables that must be provided |
| `model_hint` | string | Suggested LLM model |

## Helper Functions

All helpers are registered as `template.FuncMap` and available in every template:

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

## CLI Usage

### Render a template

```bash
promptkit render summarize \
  --dir ./templates \
  --var document="The quick brown fox..." \
  --var max_words=100
```

### Validate required variables

```bash
promptkit validate summarize --dir ./templates
```

Output:
```
Required variables for "summarize":
  - document
  - max_words
```

### List available templates

```bash
promptkit list --dir ./templates
```

Output:
```
classify             Classify text into provided categories
extract              Extract entities from text
summarize            Summarize a document with configurable length
```

### Execute a prompt chain

```bash
promptkit chain ./templates/chain_example.yaml \
  --dir ./templates \
  --var input_document="AI is transforming technology." \
  --var categories="tech, science, politics"
```

## Prompt Chaining

Define multi-step pipelines in YAML. Each step's output is captured into a variable available to subsequent steps:

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

## Project Structure

```
promptkit/
├── cmd/promptkit/          # CLI entry point
├── internal/
│   ├── chain/              # Prompt chaining pipeline
│   ├── config/             # Default configuration
│   ├── engine/             # Render engine + helper functions
│   ├── frontmatter/        # YAML frontmatter parser
│   ├── registry/           # Template directory loading
│   └── validator/          # Required variable validation
├── templates/              # Example templates
│   ├── includes/           # Reusable template blocks
│   ├── summarize.tmpl
│   ├── classify.tmpl
│   ├── extract.tmpl
│   └── chain_example.yaml
├── Makefile
└── LICENSE
```

## Development

```bash
make test       # Run all tests
make lint       # Run golangci-lint
make build      # Build the binary
make cover      # Generate coverage report
```

## License

[MIT](LICENSE)
