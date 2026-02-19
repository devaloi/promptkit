// Package engine provides the core prompt template rendering engine.
package engine

import (
	"encoding/json"
	"strings"
	"text/template"
)

// FuncMap returns the template.FuncMap with all LLM helper functions.
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"truncate":       truncate,
		"json_encode":    jsonEncode,
		"word_count":     wordCount,
		"token_estimate": tokenEstimate,
		"upper":          strings.ToUpper,
		"lower":          strings.ToLower,
		"join":           joinSlice,
		"default":        defaultVal,
	}
}

// truncate limits text to maxChars characters, appending "..." if truncated.
func truncate(maxChars int, text string) string {
	if len(text) <= maxChars {
		return text
	}
	if maxChars <= 3 {
		return text[:maxChars]
	}
	return text[:maxChars-3] + "..."
}

// jsonEncode marshals a value to a JSON string.
func jsonEncode(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// wordCount returns the number of whitespace-delimited words in text.
func wordCount(text string) int {
	return len(strings.Fields(text))
}

// tokenEstimate estimates the number of tokens in text (~4 chars per token).
func tokenEstimate(text string) int {
	n := len(text)
	if n == 0 {
		return 0
	}
	return (n + 3) / 4
}

// joinSlice joins a slice of strings with the given separator.
func joinSlice(sep string, elems []string) string {
	return strings.Join(elems, sep)
}

// defaultVal returns fallback if value is an empty string, otherwise value.
func defaultVal(fallback, value string) string {
	if value == "" {
		return fallback
	}
	return value
}
