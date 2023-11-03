package types

func Ptr[T any, DoublePtr interface{ **T }, Ptr interface{ *T }](v DoublePtr) Ptr {
	if v == nil {
		return nil
	}
	return *v
}
