package sqlgen

import (
	"strings"
	"sync"
)

type SqlStmt interface {
	WriteQuery(string, ...any)
	WriteQueryByte(byte)
	Query() string
	Args() []any
	Reset()
}

type sqlStmt struct {
	w    strings.Builder
	args []any
}

var (
	pool = sync.Pool{
		New: func() any {
			// The Pool's New function should generally only return pointer
			// types, since a pointer can be put into the return interface
			// value without an allocation:
			return new(sqlStmt)
		},
	}
)

func (s *sqlStmt) WriteQuery(query string, args ...any) {
	s.w.WriteString(query)
	s.args = append(s.args, args...)
}

func (s *sqlStmt) WriteQueryByte(b byte) {
	s.w.WriteByte(b)
}

func (s *sqlStmt) Query() string {
	return s.w.String()
}

func (s *sqlStmt) Args() []any {
	return s.args
}

func (s *sqlStmt) Reset() {
	s.w.Reset()
	s.args = nil
}

// AcquireStmt
func AcquireStmt() SqlStmt {
	return pool.Get().(SqlStmt)
}

// ReleaseStmt
func ReleaseStmt(stmt SqlStmt) {
	stmt.Reset()
	pool.Put(stmt.(*sqlStmt))
}
