# Build promptkit — LLM Prompt Template Engine in Go

You are building a **portfolio project** for a Senior AI Engineer's public GitHub. It must be impressive, clean, and production-grade. Read these docs before writing any code:

1. **`A06-go-prompt-engine.md`** — Complete project spec: architecture, template format with YAML frontmatter, helper functions, chaining pipeline, CLI design, phased build plan, commit plan. This is your primary blueprint. Follow it phase by phase.
2. **`github-portfolio.md`** — Portfolio goals and Definition of Done (Level 1 + Level 2). Understand the quality bar.
3. **`github-portfolio-checklist.md`** — Pre-publish checklist. Every item must pass before you're done.

---

## Instructions

### Read first, build second
Read all three docs completely before writing a single line of code. Understand the template format (YAML frontmatter + Go template body), the helper function map, the registry's directory-loading approach, the chaining pipeline, and the CLI commands.

### Follow the phases in order
The project spec has 4 phases. Do them in order:
1. **Foundation** — project setup, YAML frontmatter parser, LLM helper functions (truncate, json_encode, token_estimate, word_count, etc.)
2. **Engine + Validation** — core render engine with text/template and custom FuncMap, variable injection, template includes, required-variable validation
3. **Registry + Chaining** — directory-based template registry with name lookup, prompt chaining with sequential pipeline execution and variable passing between steps
4. **CLI + Polish** — CLI commands (render, validate, list, chain), example templates, README

### Commit frequently
Follow the commit plan in the spec. Use **conventional commits**. Each commit should be a logical unit.

### Quality non-negotiables
- **YAML frontmatter.** Every template has `---` delimited YAML frontmatter with `name`, `description`, `required_vars`, and `model_hint`. The parser must cleanly separate frontmatter from template body. Use `gopkg.in/yaml.v3`.
- **text/template based.** The render engine must use Go's `text/template` package — not `html/template`, not a third-party template library. LLM prompts are plain text.
- **Custom FuncMap.** All helper functions (`truncate`, `json_encode`, `word_count`, `token_estimate`, `upper`, `lower`, `join`, `default`) must be registered via `template.FuncMap` and individually tested.
- **Validation before render.** The validator checks that all `required_vars` from frontmatter are present in the provided variables map before rendering. Missing variables produce a clear, structured error listing all missing names.
- **Template includes.** Templates can reference other templates via `{{ template "name" }}`. The registry makes included templates available. The `includes/` subdirectory holds reusable blocks.
- **Prompt chaining works end-to-end.** A chain YAML defines steps; each step's rendered output is captured into a variable available to subsequent steps. The chain executor returns both the final output and all intermediates.
- **Lint clean.** `golangci-lint run` and `go vet` must pass.
- **No external LLM calls.** This is a prompt *construction* library. It renders template strings. It does not call OpenAI or any LLM API.

### What NOT to do
- Don't use `html/template`. LLM prompts are plain text — HTML escaping will corrupt them.
- Don't skip frontmatter validation. Templates without frontmatter should still work but without validation benefits.
- Don't hardcode template paths. The registry accepts a configurable directory root.
- Don't couple to any specific LLM provider. This library renders prompt strings — it is provider-agnostic.
- Don't leave `// TODO` or `// FIXME` comments anywhere.
- Don't use a third-party template engine (Mustache, Handlebars, etc.). `text/template` is the engine.

---

## GitHub Username

The GitHub username is **devaloi**. For Go module paths, use `github.com/devaloi/promptkit`. All internal imports must use this module path.

## Start

Read the three docs. Then begin Phase 1 from `A06-go-prompt-engine.md`.
