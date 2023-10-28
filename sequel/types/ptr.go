package types

func Ptr[T any, Ptr interface{ *T }](v Ptr) Ptr {
	if v == nil {
		return nil
	}
	return v
}
