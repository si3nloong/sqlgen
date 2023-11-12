package db

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/si3nloong/sqlgen/sequel"
)

func NewStmt() sequel.Stmt {
	return &sqlStmt{}
}

type sqlStmt struct {
	strings.Builder
	pos  uint
	args []any
}

func (s *sqlStmt) Var(query string, values ...any) {
	s.WriteString(query)
	noOfLen := len(values)
	if noOfLen == 1 {
		s.WriteByte('?')
	} else if noOfLen > 1 {
		s.WriteString("(" + strings.Repeat("?,", noOfLen)[:(noOfLen*2)-1] + ")")
	}
	s.args = append(s.args, values...)
	s.pos++
}

func (s sqlStmt) Args() []any {
	return s.args
}

func (s *sqlStmt) Reset() {
	s.args = nil
	s.pos = 0
	s.Builder.Reset()
}

type SelectStmt struct {
	Select    []string
	FromTable string
	Where     sequel.WhereClause
	OrderBy   []string
	Limit     uint16
}

func QueryStmt[T any, Ptr interface {
	*T
	sequel.Scanner[T]
}, Stmt interface{ SelectStmt }](ctx context.Context, dbConn sequel.DB, stmt Stmt) ([]T, error) {
	blr := NewStmt()
	switch vi := any(stmt).(type) {
	case SelectStmt:
		blr.WriteString("SELECT ")
		for i := range vi.Select {
			if i > 0 {
				blr.WriteByte(',')
			}
			blr.WriteString(vi.Select[i])
		}
		blr.WriteString(" FROM " + vi.FromTable)
		if vi.Where != nil {
			blr.WriteString(" WHERE ")
			vi.Where(blr)
		}
		if len(vi.OrderBy) > 0 {
			blr.WriteString(" ORDER BY ")
			for i := range vi.OrderBy {
				if i > 0 {
					blr.WriteByte(',')
				}
				blr.WriteString(vi.OrderBy[i])
			}
		}
		if vi.Limit > 0 {
			blr.WriteString(" LIMIT " + strconv.FormatUint(uint64(vi.Limit), 10))
		}
		blr.WriteByte(';')
	}

	log.Println(blr.String())
	rows, err := dbConn.QueryContext(ctx, blr.String(), blr.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []T
	for rows.Next() {
		var v T
		if err := rows.Scan(Ptr(&v).Addrs()...); err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}

type UpdateStmt struct {
	FromTable string
	Set       []string
	Where     sequel.WhereClause
	OrderBy   []sequel.OrderByClause
	Limit     uint16
}

type DeleteStmt struct {
	FromTable string
	Where     sequel.WhereClause
	OrderBy   []sequel.OrderByClause
	Limit     uint16
}

func ExecStmt[T any, Stmt interface {
	UpdateStmt | DeleteStmt
}](ctx context.Context, dbConn sequel.DB, stmt Stmt) error {
	blr := NewStmt()

	switch vi := any(stmt).(type) {
	case UpdateStmt:
		blr.WriteString("UPDATE " + vi.FromTable)
		if vi.Where != nil {
			blr.WriteString(" WHERE ")
			vi.Where(blr)
		}
		if len(vi.Set) > 0 {
			blr.WriteString(" SET ")
		}
		if len(vi.OrderBy) > 0 {
			blr.WriteString(" ORDER BY ")
			for i := range vi.OrderBy {
				if i > 0 {
					blr.WriteByte(',')
				}
				vi.OrderBy[i](blr)
			}
		}
		if vi.Limit > 0 {
			blr.WriteString(" LIMIT " + strconv.FormatUint(uint64(vi.Limit), 10))
		}
		blr.WriteByte(';')

	case DeleteStmt:
		blr.WriteString("DELETE FROM " + vi.FromTable)
		if vi.Where != nil {
			blr.WriteString(" WHERE ")
			vi.Where(blr)
		}
		if len(vi.OrderBy) > 0 {
			blr.WriteString(" ORDER BY ")
			for i := range vi.OrderBy {
				if i > 0 {
					blr.WriteByte(',')
				}
				vi.OrderBy[i](blr)
			}
		}
		if vi.Limit > 0 {
			blr.WriteString(" LIMIT " + strconv.FormatUint(uint64(vi.Limit), 10))
		}
		blr.WriteByte(';')
	}
	return nil
}

func And(stmts ...sequel.WhereClause) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		stmt.WriteByte('(')
		for i := range stmts {
			if i > 0 {
				stmt.WriteString(" AND ")
			}
			stmts[i](stmt)
		}
		stmt.WriteByte(')')
	}
}

func Or(stmts ...sequel.WhereClause) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		stmt.WriteByte('(')
		for i := range stmts {
			if i > 0 {
				stmt.WriteString(" OR ")
			}
			stmts[i](stmt)
		}
		stmt.WriteByte(')')
	}
}

func Equal[T comparable](f sequel.ColumnValuer[T], value T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		stmt.Var(f.ColumnName()+" = ", f.Convert(value))
	}
}

func NotEqual[T comparable](f sequel.ColumnValuer[T], value T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		stmt.Var(f.ColumnName()+" <> ", f.Convert(value))
	}
}

func In[T any](f sequel.ColumnValuer[T], values ...T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		args := make([]any, len(values))
		for idx := range values {
			args[idx] = f.Convert(values[idx])
		}
		stmt.Var(f.ColumnName()+" IN ", args...)
	}
}

func NotIn[T any](f sequel.ColumnValuer[T], values ...T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		args := make([]any, len(values))
		for idx := range values {
			args[idx] = f.Convert(values[idx])
		}
		stmt.Var(f.ColumnName()+" NOT IN ", args...)
	}
}

func GreaterThan[T comparable](f sequel.ColumnValuer[T], value T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		stmt.Var(f.ColumnName()+" > ", f.Convert(value))
	}
}

func GreaterThanOrEqual[T comparable](f sequel.ColumnValuer[T], value T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		stmt.Var(f.ColumnName()+" >= ", f.Convert(value))
	}
}

func LessThan[T comparable](f sequel.ColumnValuer[T], value T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		stmt.Var(f.ColumnName()+" < ", f.Convert(value))
	}
}

func LessThanOrEqual[T comparable](f sequel.ColumnValuer[T], value T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		stmt.Var(f.ColumnName()+" >= ", f.Convert(value))
	}
}

func Like[T comparable](f sequel.ColumnValuer[T], value T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		stmt.Var(f.ColumnName()+" LIKE ", f.Convert(value))
	}
}

func NotLike[T comparable](f sequel.ColumnValuer[T], value T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		stmt.Var(f.ColumnName()+" NOT LIKE ", f.Convert(value))
	}
}

func IsNull[T any](f sequel.ColumnValuer[T]) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		stmt.WriteString(f.ColumnName() + " IS NULL")
	}
}

func IsNotNull[T any](f sequel.ColumnValuer[T]) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		stmt.WriteString(f.ColumnName() + " IS NOT NULL")
	}
}

func Asc[T any](f sequel.ColumnValuer[T]) sequel.OrderByClause {
	return func(sw sequel.StmtWriter) {
		sw.WriteString(f.ColumnName() + " ASC")
	}
}

func Desc[T any](f sequel.ColumnValuer[T]) sequel.OrderByClause {
	return func(sw sequel.StmtWriter) {
		sw.WriteString(f.ColumnName() + " DESC")
	}
}
