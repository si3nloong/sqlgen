package sql

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
			stmt := new(strings.Builder)
			return stmt
		},
	}
)

func acquireString() *strings.Builder {
	return pool.Get().(*strings.Builder)
}

func releaseString(blr *strings.Builder) {
	blr.Reset()
	pool.Put(blr)
}
