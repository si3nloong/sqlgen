package strpool

import (
	"strings"
	"sync"
)

var (
	pool = sync.Pool{
		New: func() any {
			// The Pool's New function should generally only return pointer
			// types, since a pointer can be put into the return interface
			// value without an allocation:
			return new(strings.Builder)
		},
	}
)

func AcquireString() *strings.Builder {
	return pool.Get().(*strings.Builder)
}

func ReleaseString(blr *strings.Builder) {
	blr.Reset()
	pool.Put(blr)
}
