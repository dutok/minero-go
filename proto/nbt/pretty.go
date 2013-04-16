package nbt

import (
	"strings"
)

func Pretty(s string) string {
	var indent int
	var l []rune
	var lines []string

	// Split by lines
	for _, r := range s {
		l = append(l, r)
		switch r {
		case '{', '[', ',':
			lines = append(lines, string(l))
			l = nil
		}
	}
	// Append final line (compound close brace)
	lines = append(lines, string(l))

	for index, line := range lines {
		switch {
		case strings.HasSuffix(line, "{"), strings.HasSuffix(line, "["):
			lines[index] = tabs(indent) + line
			indent++
		case strings.HasPrefix(line, "}"), strings.HasPrefix(line, "]"):
			indent--
			lines[index] = tabs(indent) + line
		default:
			lines[index] = tabs(indent) + line
		}
	}

	return strings.Join(lines, "\n")
}

func tabs(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat("\t", n)
}
