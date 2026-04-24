package reconcile

import (
	"fmt"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Action represents the type of reconciliation action needed.
type Action int

const (
	ActionAdd Action = iota
	ActionRemove
	ActionUpdate
)

// Step describes a single reconciliation step to bring target in sync with source.
type Step struct {
	Action Action
	Key    string
	Value  string
}

// String returns a human-readable representation of the step.
func (s Step) String() string {
	switch s.Action {
	case ActionAdd:
		return fmt.Sprintf("+ %s=%s", s.Key, s.Value)
	case ActionRemove:
		return fmt.Sprintf("- %s", s.Key)
	case ActionUpdate:
		return fmt.Sprintf("~ %s=%s", s.Key, s.Value)
	default:
		return ""
	}
}

// Plan generates an ordered list of steps to reconcile target with source.
// Steps are sorted by key for deterministic output.
func Plan(entries []diff.Entry) []Step {
	steps := make([]Step, 0, len(entries))

	for _, e := range entries {
		switch e.Status {
		case diff.Added:
			steps = append(steps, Step{Action: ActionAdd, Key: e.Key, Value: e.NewValue})
		case diff.Removed:
			steps = append(steps, Step{Action: ActionRemove, Key: e.Key})
		case diff.Modified:
			steps = append(steps, Step{Action: ActionUpdate, Key: e.Key, Value: e.NewValue})
		}
	}

	sort.Slice(steps, func(i, j int) bool {
		return steps[i].Key < steps[j].Key
	})

	return steps
}

// Apply applies reconciliation steps to an existing env map, returning a new map.
func Apply(base map[string]string, steps []Step) map[string]string {
	result := make(map[string]string, len(base))
	for k, v := range base {
		result[k] = v
	}

	for _, s := range steps {
		switch s.Action {
		case ActionAdd, ActionUpdate:
			result[s.Key] = s.Value
		case ActionRemove:
			delete(result, s.Key)
		}
	}

	return result
}

// Format renders steps as a patch-style string.
func Format(steps []Step) string {
	lines := make([]string, 0, len(steps))
	for _, s := range steps {
		lines = append(lines, s.String())
	}
	return strings.Join(lines, "\n")
}
