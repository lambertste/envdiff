package env

import "fmt"

// PivotResult holds the result of pivoting an env set by a grouping key.
type PivotResult struct {
	// Columns are the distinct values of the pivot key (e.g. "staging", "prod").
	Columns []string
	// Rows maps each non-pivot key to a map of column -> value.
	Rows map[string]map[string]string
}

// Pivot restructures the env set by treating pivotKey's value as a column
// header and grouping all other keys under it. Useful when multiple env files
// are merged with a SCOPE or ENV key distinguishing them.
//
// Example input (merged set):
//   ENV=staging  DB_HOST=db1  API_KEY=abc
//   ENV=prod     DB_HOST=db2  API_KEY=xyz
//
// Pivot(set, "ENV") produces columns=["prod","staging"] and rows for DB_HOST, API_KEY.
func Pivot(sets []*Set, pivotKey string) (*PivotResult, error) {
	if len(sets) == 0 {
		return &PivotResult{Rows: make(map[string]map[string]string)}, nil
	}

	colSet := map[string]struct{}{}
	var columns []string

	for _, s := range sets {
		col, ok := s.Get(pivotKey)
		if !ok || col == "" {
			return nil, fmt.Errorf("pivot: set missing pivot key %q", pivotKey)
		}
		if _, seen := colSet[col]; !seen {
			colSet[col] = struct{}{}
			columns = append(columns, col)
		}
	}

	rows := make(map[string]map[string]string)

	for _, s := range sets {
		col, _ := s.Get(pivotKey)
		for _, k := range s.Keys() {
			if k == pivotKey {
				continue
			}
			if rows[k] == nil {
				rows[k] = make(map[string]string)
			}
			v, _ := s.Get(k)
			rows[k][col] = v
		}
	}

	return &PivotResult{
		Columns: columns,
		Rows:    rows,
	}, nil
}

// FormatPivot returns a human-readable table string for the pivot result.
func FormatPivot(pr *PivotResult) string {
	if len(pr.Columns) == 0 || len(pr.Rows) == 0 {
		return "(empty pivot)\n"
	}

	out := fmt.Sprintf("%-30s", "KEY")
	for _, col := range pr.Columns {
		out += fmt.Sprintf("  %-20s", col)
	}
	out += "\n"

	keys := make([]string, 0, len(pr.Rows))
	for k := range pr.Rows {
		keys = append(keys, k)
	}
	sortStrings(keys)

	for _, k := range keys {
		row := fmt.Sprintf("%-30s", k)
		for _, col := range pr.Columns {
			v := pr.Rows[k][col]
			if v == "" {
				v = "-"
			}
			row += fmt.Sprintf("  %-20s", v)
		}
		out += row + "\n"
	}
	return out
}
