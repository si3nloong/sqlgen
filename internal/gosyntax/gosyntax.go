package gosyntax

func PtrOf[T any](v T) *T {
	return &v
}
