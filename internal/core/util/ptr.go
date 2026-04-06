package util

// Deref returns the T, de-referring the given `ptr` or the default (zero)
// T object, if given `ptr` is nil.
func Deref[T any](ptr *T) T {
	return DerefDefault(ptr, *new(T))
}

// DerefDefault is the same as [Deref], but returns provided `defaultValue`
// instead of Go's default one.
func DerefDefault[T any](ptr *T, defaultValue T) T {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// Ref returns *T to the given T.
func Ref[T any](v T) *T {
	return &v
}
