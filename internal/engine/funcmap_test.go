package engine

import (
	"testing"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		max      int
		input    string
		expected string
	}{
		{"no truncation needed", 10, "hello", "hello"},
		{"exact length", 5, "hello", "hello"},
		{"truncated with ellipsis", 8, "hello world", "hello..."},
		{"very short max", 2, "hello", "he"},
		{"empty string", 10, "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncate(tt.max, tt.input)
			if got != tt.expected {
				t.Errorf("truncate(%d, %q) = %q, want %q", tt.max, tt.input, got, tt.expected)
			}
		})
	}
}

func TestJSONEncode(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"string", "hello", `"hello"`},
		{"number", 42, "42"},
		{"map", map[string]string{"key": "val"}, `{"key":"val"}`},
		{"nil", nil, "null"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsonEncode(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("jsonEncode(%v) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestWordCount(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"normal sentence", "hello world foo", 3},
		{"single word", "hello", 1},
		{"empty string", "", 0},
		{"extra whitespace", "  hello   world  ", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wordCount(tt.input)
			if got != tt.expected {
				t.Errorf("wordCount(%q) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}

func TestTokenEstimate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"empty", "", 0},
		{"short", "hi", 1},
		{"four chars", "test", 1},
		{"five chars", "hello", 2},
		{"sixteen chars", "1234567890123456", 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tokenEstimate(tt.input)
			if got != tt.expected {
				t.Errorf("tokenEstimate(%q) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}

func TestJoinSlice(t *testing.T) {
	tests := []struct {
		name     string
		sep      string
		elems    []string
		expected string
	}{
		{"comma join", ", ", []string{"a", "b", "c"}, "a, b, c"},
		{"empty slice", ", ", []string{}, ""},
		{"single element", ", ", []string{"only"}, "only"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := joinSlice(tt.sep, tt.elems)
			if got != tt.expected {
				t.Errorf("joinSlice(%q, %v) = %q, want %q", tt.sep, tt.elems, got, tt.expected)
			}
		})
	}
}

func TestDefaultVal(t *testing.T) {
	tests := []struct {
		name     string
		fallback string
		value    string
		expected string
	}{
		{"empty value uses fallback", "default_val", "", "default_val"},
		{"non-empty value used", "default_val", "actual", "actual"},
		{"both empty", "", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := defaultVal(tt.fallback, tt.value)
			if got != tt.expected {
				t.Errorf("defaultVal(%q, %q) = %q, want %q", tt.fallback, tt.value, got, tt.expected)
			}
		})
	}
}

func TestFuncMapRegistered(t *testing.T) {
	fm := FuncMap()
	expected := []string{"truncate", "json_encode", "word_count", "token_estimate", "upper", "lower", "join", "default"}
	for _, name := range expected {
		if _, ok := fm[name]; !ok {
			t.Errorf("FuncMap missing function %q", name)
		}
	}
}
