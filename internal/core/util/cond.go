package util

// If is because Golang DO NEED the fucking ternary operator.
func If[T any](cond bool, a, b T) T {
	if cond {
		return a
	} else {
		return b
	}
}

// Or returns first of not-zero value of given ones.
func Or[T comparable](v1, v2 T) T {
	if IsNotZero(v1) {
		return v1
	} else {
		return v2
	}
}

// OrTo is the same as [Or], but compares `v` against the value stored
// by the ref of `to` (if it's not nil).
func OrTo[T comparable](to *T, v T) {
	if to != nil {
		*to = Or(*to, v)
	}
}
