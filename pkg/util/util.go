package util

// Mimics a ternary operator found in other languages.
func If[T any](cond bool, tval, fval T) T {
	if cond {
		return tval
	} else {
		return fval
	}
}
