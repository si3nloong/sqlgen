package cli

// assertAs will assert the value to the data type you defined
func assertAs[T any](v any) *T {
	vi, ok := v.(*T)
	if !ok {
		return nil
	}
	return vi
}
