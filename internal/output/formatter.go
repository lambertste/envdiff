package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// ColorCode represents ANSI color escape codes.
type ColorCode string

const (
	ColorReset  ColorCode = "\033[0m"
	ColorRed    ColorCode = "\033[31m"
	ColorGreen  ColorCode = "\033[32m"
	ColorYellow ColorCode = "\033[33m"
	ColorCyan   ColorCode = "\033[36m"
)

// Format controls output rendering style.
type Format string

const (
	FormatText  Format = "text"
	FormatColor Format = "color"
	FormatDotenv Format = "dotenv"
)

// Write renders a slice of diff.Entry to the given writer using the specified format.
func Write(w io.Writer, entries []diff.Entry, format Format) error {
	for _, e := range entries {
		var line string
		switch format {
		case FormatColor:
			line = colorLine(e)
		case FormatDotenv:
			line = dotenvLine(e)
		default:
			line = textLine(e)
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func textLine(e diff.Entry) string {
	switch e.Kind {
	case diff.Added:
		return fmt.Sprintf("+ %s=%s", e.Key, e.NewValue)
	case diff.Removed:
		return fmt.Sprintf("- %s=%s", e.Key, e.OldValue)
	case diff.Modified:
		return fmt.Sprintf("~ %s: %s -> %s", e.Key, e.OldValue, e.NewValue)
	default:
		return fmt.Sprintf("  %s=%s", e.Key, e.OldValue)
	}
}

func colorLine(e diff.Entry) string {
	raw := textLine(e)
	switch e.Kind {
	case diff.Added:
		return strings.Join([]string{string(ColorGreen), raw, string(ColorReset)}, "")
	case diff.Removed:
		return strings.Join([]string{string(ColorRed), raw, string(ColorReset)}, "")
	case diff.Modified:
		return strings.Join([]string{string(ColorYellow), raw, string(ColorReset)}, "")
	default:
		return raw
	}
}

func dotenvLine(e diff.Entry) string {
	switch e.Kind {
	case diff.Added, diff.Modified:
		return fmt.Sprintf("%s=%s", e.Key, e.NewValue)
	case diff.Removed:
		return fmt.Sprintf("# REMOVED: %s", e.Key)
	default:
		return fmt.Sprintf("%s=%s", e.Key, e.OldValue)
	}
}
