package qparams

// ptr is a helper that returns a pointer to v.
func ptr[T any](v T) *T {
	return &v
}
