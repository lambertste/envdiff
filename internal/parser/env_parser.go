package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap represents a parsed environment file as a key-value map.
type EnvMap map[string]string

// ParseFile reads and parses a .env file from the given path.
// It skips blank lines and comments (lines starting with #).
func ParseFile(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening env file %q: %w", path, err)
	}
	defer f.Close()

	return ParseReader(bufio.NewScanner(f))
}

// ParseReader parses environment variables from a bufio.Scanner source.
func ParseReader(scanner *bufio.Scanner) (EnvMap, error) {
	env := make(EnvMap)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseLine(line, lineNum)
		if err != nil {
			return nil, err
		}
		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning env file: %w", err)
	}

	return env, nil
}

// parseLine splits a single KEY=VALUE line into its components.
func parseLine(line string, lineNum int) (string, string, error) {
	idx := strings.IndexByte(line, '=')
	if idx < 0 {
		return "", "", fmt.Errorf("line %d: missing '=' in %q", lineNum, line)
	}

	key := strings.TrimSpace(line[:idx])
	value := strings.TrimSpace(line[idx+1:])

	if key == "" {
		return "", "", fmt.Errorf("line %d: empty key", lineNum)
	}

	// Strip optional surrounding quotes from value.
	value = stripQuotes(value)

	return key, value, nil
}

// stripQuotes removes matching surrounding single or double quotes.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
