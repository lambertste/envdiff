package env

import "fmt"

// PatchOp represents a single patch operation kind.
type PatchOp string

const (
	PatchSet    PatchOp = "set"
	PatchDelete PatchOp = "delete"
	PatchRename PatchOp = "rename"
)

// PatchInstruction describes one mutation to apply to an EnvSet.
type PatchInstruction struct {
	Op      PatchOp
	Key     string
	Value   string // used by PatchSet
	NewKey  string // used by PatchRename
}

// Patch applies a slice of PatchInstructions to a copy of the given Set
// and returns the modified copy. The original Set is not mutated.
func Patch(s *Set, instructions []PatchInstruction) (*Set, error) {
	out := Clone(s)
	for _, ins := range instructions {
		switch ins.Op {
		case PatchSet:
			if ins.Key == "" {
				return nil, fmt.Errorf("patch set: key must not be empty")
			}
			out.Set(ins.Key, ins.Value)
		case PatchDelete:
			if ins.Key == "" {
				return nil, fmt.Errorf("patch delete: key must not be empty")
			}
			out.Delete(ins.Key)
		case PatchRename:
			if ins.Key == "" || ins.NewKey == "" {
				return nil, fmt.Errorf("patch rename: both key and new_key must not be empty")
			}
			v, ok := out.Get(ins.Key)
			if !ok {
				return nil, fmt.Errorf("patch rename: key %q not found", ins.Key)
			}
			out.Delete(ins.Key)
			out.Set(ins.NewKey, v)
		default:
			return nil, fmt.Errorf("unknown patch op: %q", ins.Op)
		}
	}
	return out, nil
}
