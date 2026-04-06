package util

import (
	"cmp"
)

// IsZero reports whether given value is default (zero) one.
// Check whether integers are zero, strings are empty, pointers are nil, etc.
func IsZero[T comparable](v T) bool {
	return v == *new(T)
}

// IsNotZero is opposite to the [IsDefault].
// Useful when you are passing func as argument.
func IsNotZero[T comparable](v T) bool {
	return !IsZero(v)
}

// IsInRange reports whether given value belongs provided [a...b] range.
func IsInRange[T cmp.Ordered](v, a, b T) bool {
	return a <= v && v <= b
}

// Clamp is a clamp.
// https://en.wikipedia.org/wiki/Clamp_(function)
func Clamp[T cmp.Ordered](v, a, b T) T {
	return max(min(a, b), min(max(a, b), v))
}

// ClampTo is the same as [Clamp] but using given pointer
// as both source and destination of the value.
func ClampTo[T cmp.Ordered](v *T, a, b T) {
	if v != nil {
		*v = Clamp(*v, a, b)
	}
}
