package gitutil

import "regexp"

func SplitWithDelimiters(text string, pattern *regexp.Regexp) []string {
	matches := pattern.FindAllStringIndex(text, -1)
	parts := make([]string, 0, len(matches)*2+1)
	start := 0
	for _, match := range matches {
		parts = append(parts, text[start:match[0]], text[match[0]:match[1]])
		start = match[1]
	}
	parts = append(parts, text[start:])
	return parts
}
