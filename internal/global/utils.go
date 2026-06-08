package global

import "strings"

// FirstNonEmpty returns the first non-empty string from the given values.
func FirstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

// FirstNonEmptyTrimmed returns the first non-empty string after trimming whitespace from the given values.
func FirstNonEmptyTrimmed(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

// TrimAndCollapse whitespace and collapse multiple spaces into single spaces.
func TrimAndCollapse(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// TrimRunes truncates a string to a maximum number of runes, appending an ellipsis if truncated.
func TrimRunes(s string, max int) string {
	if max <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	if max <= 1 {
		return "…"
	}
	return string(runes[:max-1]) + "…"
}

// TrimLines removes trailing empty line from a slice of strings.
func TrimLines(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}
