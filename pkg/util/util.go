package util

import "strings"

// Escapes a string for usage in Markdown.
func EscapeMD(s string) string {
	replacer := strings.NewReplacer(
		"*", "\\*",
		"_", "\\_",
		"`", "\\`",
		"#", "\\#",
		"-", "\\-",
		"+", "\\+",
		".", "\\.",
		"!", "\\!",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
	)
	return replacer.Replace(s)
}

// Mimics a ternary operator found in other languages.
func If[T any](cond bool, tval, fval T) T {
	if cond {
		return tval
	} else {
		return fval
	}
}
