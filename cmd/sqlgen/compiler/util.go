package compiler

import "go/types"

// This is to count the underlying pointer of a type
func ptrCount(t types.Type) int {
	var total int
loop:
	for t != nil {
		if v, ok := t.(*types.Pointer); ok {
			total++
			t = v.Elem()
		} else {
			break loop
		}
	}
	return total
}

func assertAsPtr[T any](v any) *T {
	t, ok := v.(*T)
	if ok {
		return t
	}
	return nil
}
