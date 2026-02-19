package validator

import (
	"errors"
	"testing"
)

func TestValidate_AllPresent(t *testing.T) {
	required := []string{"name", "age"}
	vars := map[string]any{"name": "Alice", "age": 30}

	err := Validate(required, vars)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidate_OneMissing(t *testing.T) {
	required := []string{"name", "age"}
	vars := map[string]any{"name": "Alice"}

	err := Validate(required, vars)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var mve *MissingVarsError
	if !errors.As(err, &mve) {
		t.Fatalf("expected *MissingVarsError, got %T", err)
	}
	if len(mve.Missing) != 1 || mve.Missing[0] != "age" {
		t.Errorf("expected [age], got %v", mve.Missing)
	}
}

func TestValidate_MultipleMissing(t *testing.T) {
	required := []string{"a", "b", "c"}
	vars := map[string]any{"b": "present"}

	err := Validate(required, vars)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var mve *MissingVarsError
	if !errors.As(err, &mve) {
		t.Fatalf("expected *MissingVarsError, got %T", err)
	}
	if len(mve.Missing) != 2 {
		t.Errorf("expected 2 missing vars, got %d: %v", len(mve.Missing), mve.Missing)
	}
}

func TestValidate_EmptyRequired(t *testing.T) {
	err := Validate(nil, map[string]any{"foo": "bar"})
	if err != nil {
		t.Fatalf("expected no error for empty required_vars, got %v", err)
	}
}

func TestValidate_EmptyVars(t *testing.T) {
	err := Validate([]string{"x"}, nil)
	if err == nil {
		t.Fatal("expected error for nil vars with required var")
	}
}

func TestMissingVarsError_Message(t *testing.T) {
	e := &MissingVarsError{Missing: []string{"foo", "bar"}}
	expected := "missing required variables: foo, bar"
	if e.Error() != expected {
		t.Errorf("expected %q, got %q", expected, e.Error())
	}
}
