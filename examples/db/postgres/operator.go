package postgresdb

import (
	"github.com/si3nloong/sqlgen/sequel"
)

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
		stmt.Vars(f.ColumnName()+" IN ", args)
	}
}

func NotIn[T any](f sequel.ColumnValuer[T], values ...T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		args := make([]any, len(values))
		for idx := range values {
			args[idx] = f.Convert(values[idx])
		}
		stmt.Vars(f.ColumnName()+" NOT IN ", args)
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
		stmt.Var(f.ColumnName()+" <= ", f.Convert(value))
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

func Between[T comparable](f sequel.ColumnValuer[T], from, to T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		stmt.Var(f.ColumnName()+" BETWEEN ", from)
		stmt.Var(" AND ", to)
	}
}

func Set[T any](f sequel.ColumnValuer[T], value ...T) sequel.SetClause {
	return func(stmt sequel.StmtBuilder) {
		defaultValue := f.Value()
		if len(value) > 0 {
			defaultValue = f.Convert(value[0])
		}
		stmt.Var(f.ColumnName()+" = ", defaultValue)
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
