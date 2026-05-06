package env

import (
	"fmt"
	"regexp"
	"strings"
)

// SchemaField describes a single expected environment variable.
type SchemaField struct {
	Key         string
	Required    bool
	Pattern     *regexp.Regexp // optional value pattern
	Description string
}

// SchemaViolation records a single schema validation failure.
type SchemaViolation struct {
	Key     string
	Message string
}

func (v SchemaViolation) String() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// Schema holds an ordered list of field definitions.
type Schema struct {
	Fields []SchemaField
}

// Validate checks the given Set against the schema and returns any violations.
func (s *Schema) Validate(set *Set) []SchemaViolation {
	var violations []SchemaViolation

	for _, field := range s.Fields {
		val, ok := set.Get(field.Key)
		if !ok || strings.TrimSpace(val) == "" {
			if field.Required {
				violations = append(violations, SchemaViolation{
					Key:     field.Key,
					Message: "required key is missing or empty",
				})
			}
			continue
		}
		if field.Pattern != nil && !field.Pattern.MatchString(val) {
			violations = append(violations, SchemaViolation{
				Key:     field.Key,
				Message: fmt.Sprintf("value %q does not match pattern %s", val, field.Pattern),
			})
		}
	}

	return violations
}

// FormatViolations returns a human-readable summary of schema violations.
func FormatViolations(violations []SchemaViolation) string {
	if len(violations) == 0 {
		return "schema: all checks passed"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("schema: %d violation(s)\n", len(violations)))
	for _, v := range violations {
		sb.WriteString("  - ")
		sb.WriteString(v.String())
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}
