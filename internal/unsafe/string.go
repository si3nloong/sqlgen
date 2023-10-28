package unsafe

import "unsafe"

func String[T ~string | ~[]byte](v T) string {
	return *(*string)(unsafe.Pointer(&v))
}
