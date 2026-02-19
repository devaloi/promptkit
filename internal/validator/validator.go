// Package validator checks required template variables are present before rendering.
package validator

import (
	"fmt"
	"strings"
)

// MissingVarsError is returned when required variables are not provided.
type MissingVarsError struct {
	Missing []string
}

func (e *MissingVarsError) Error() string {
	return fmt.Sprintf("missing required variables: %s", strings.Join(e.Missing, ", "))
}

// Validate checks that all requiredVars are present as keys in vars.
// Returns a *MissingVarsError listing any missing variables, or nil if all are present.
func Validate(requiredVars []string, vars map[string]any) error {
	var missing []string
	for _, name := range requiredVars {
		if _, ok := vars[name]; !ok {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return &MissingVarsError{Missing: missing}
	}
	return nil
}
