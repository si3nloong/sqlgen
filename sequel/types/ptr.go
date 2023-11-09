package types

func PtrOf[T any](v T) *T {
	return &v
}

// func DerefPtr[T any, DoublePtr interface{ **T }, Ptr interface{ *T }](v DoublePtr) Ptr {
// 	if v == nil {
// 		*(v) = new(T)
// 		return *v
// 	}
// 	return *v
// }
