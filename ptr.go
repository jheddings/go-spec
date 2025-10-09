package gospec

// helper function to return a pointer to a value.
// useful for values or params that may be nil.
func Ptr[T any](v T) *T {
	return &v
}
