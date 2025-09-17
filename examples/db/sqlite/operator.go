package sqlite

import (
	"github.com/si3nloong/sqlgen/sequel"
)

func And(clauses ...sequel.WhereClause) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		if n := len(clauses); n > 0 {
			stmt.WriteByte('(')
			clauses[0](stmt)
			for i := 1; i < n; i++ {
				stmt.WriteString(" AND ")
				clauses[i](stmt)
			}
			stmt.WriteByte(')')
		}
	}
}

func Or(clauses ...sequel.WhereClause) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		if n := len(clauses); n > 0 {
			stmt.WriteByte('(')
			clauses[0](stmt)
			for i := 1; i < n; i++ {
				stmt.WriteString(" OR ")
				clauses[i](stmt)
			}
			stmt.WriteByte(')')
		}
	}
}

func Equal[T comparable](column sequel.ColumnClause[T], value T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		switch vi := column.(type) {
		case sequel.SQLColumnClause[T]:
			stmt.WriteString(vi.ColumnName() + " = " + stmt.Var(vi.Convert(value)))
		case sequel.ColumnConvertClause[T]:
			stmt.WriteString(vi.ColumnName() + " = " + stmt.Var(vi.Convert(value)))
		default:
			stmt.WriteString(column.ColumnName() + " = " + stmt.Var(value))
		}
	}
}

func NotEqual[T comparable](column sequel.ColumnClause[T], value T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		switch vi := column.(type) {
		case sequel.SQLColumnClause[T]:
			stmt.WriteString(vi.ColumnName() + " <> " + stmt.Var(vi.Convert(value)))
		case sequel.ColumnConvertClause[T]:
			stmt.WriteString(vi.ColumnName() + " <> " + stmt.Var(vi.Convert(value)))
		default:
			stmt.WriteString(column.ColumnName() + " <> " + stmt.Var(value))
		}
	}
}

func In[T any](column sequel.ColumnClause[T], values []T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		args := make([]any, len(values))
		switch vi := column.(type) {
		case sequel.SQLColumnClause[T]:
			for idx := range values {
				args[idx] = vi.Convert(values[idx])
			}
		case sequel.ColumnConvertClause[T]:
			for idx := range values {
				args[idx] = vi.Convert(values[idx])
			}
		default:
			for idx := range values {
				args[idx] = values[idx]
			}
		}
		stmt.WriteString(column.ColumnName() + " IN " + stmt.Vars(args))
	}
}

func NotIn[T any](column sequel.ColumnClause[T], values []T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		args := make([]any, len(values))
		switch vi := column.(type) {
		case sequel.SQLColumnClause[T]:
			for idx := range values {
				args[idx] = vi.Convert(values[idx])
			}
		case sequel.ColumnConvertClause[T]:
			for idx := range values {
				args[idx] = vi.Convert(values[idx])
			}
		default:
			for idx := range values {
				args[idx] = values[idx]
			}
		}
		stmt.WriteString(column.ColumnName() + " NOT IN " + stmt.Vars(args))
	}
}

func GreaterThan[T comparable](column sequel.ColumnClause[T], value T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		switch vi := column.(type) {
		case sequel.SQLColumnClause[T]:
			stmt.WriteString(vi.ColumnName() + " > " + stmt.Var(vi.Convert(value)))
		case sequel.ColumnConvertClause[T]:
			stmt.WriteString(vi.ColumnName() + " > " + stmt.Var(vi.Convert(value)))
		default:
			stmt.WriteString(column.ColumnName() + " > " + stmt.Var(value))
		}
	}
}

func GreaterThanOrEqual[T comparable](column sequel.ColumnClause[T], value T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		switch vi := column.(type) {
		case sequel.SQLColumnClause[T]:
			stmt.WriteString(vi.ColumnName() + " >= " + stmt.Var(vi.Convert(value)))
		case sequel.ColumnConvertClause[T]:
			stmt.WriteString(vi.ColumnName() + " >= " + stmt.Var(vi.Convert(value)))
		default:
			stmt.WriteString(column.ColumnName() + " >= " + stmt.Var(value))
		}
	}
}

func LessThan[T comparable](column sequel.ColumnClause[T], value T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		switch vi := column.(type) {
		case sequel.SQLColumnClause[T]:
			stmt.WriteString(vi.ColumnName() + " < " + stmt.Var(vi.Convert(value)))
		case sequel.ColumnConvertClause[T]:
			stmt.WriteString(vi.ColumnName() + " < " + stmt.Var(vi.Convert(value)))
		default:
			stmt.WriteString(column.ColumnName() + " < " + stmt.Var(value))
		}
	}
}

func LessThanOrEqual[T comparable](column sequel.ColumnClause[T], value T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		switch vi := column.(type) {
		case sequel.SQLColumnClause[T]:
			stmt.WriteString(vi.ColumnName() + " <= " + stmt.Var(vi.Convert(value)))
		case sequel.ColumnConvertClause[T]:
			stmt.WriteString(vi.ColumnName() + " <= " + stmt.Var(vi.Convert(value)))
		default:
			stmt.WriteString(column.ColumnName() + " <= " + stmt.Var(value))
		}
	}
}

func Like[T comparable](column sequel.ColumnClause[T], value T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		switch vi := column.(type) {
		case sequel.SQLColumnClause[T]:
			stmt.WriteString(vi.ColumnName() + " LIKE " + stmt.Var(vi.Convert(value)))
		case sequel.ColumnConvertClause[T]:
			stmt.WriteString(vi.ColumnName() + " LIKE " + stmt.Var(vi.Convert(value)))
		default:
			stmt.WriteString(column.ColumnName() + " LIKE " + stmt.Var(value))
		}
	}
}

func NotLike[T comparable](column sequel.ColumnClause[T], value T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		switch vi := column.(type) {
		case sequel.SQLColumnClause[T]:
			stmt.WriteString(vi.ColumnName() + " NOT LIKE " + stmt.Var(vi.Convert(value)))
		case sequel.ColumnConvertClause[T]:
			stmt.WriteString(vi.ColumnName() + " NOT LIKE " + stmt.Var(vi.Convert(value)))
		default:
			stmt.WriteString(column.ColumnName() + " NOT LIKE " + stmt.Var(value))
		}
	}
}

func IsNull[T any](column sequel.ColumnClause[T]) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		stmt.WriteString(column.ColumnName() + " IS NULL")
	}
}

func IsNotNull[T any](column sequel.ColumnClause[T]) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		stmt.WriteString(column.ColumnName() + " IS NOT NULL")
	}
}

func Between[T comparable](column sequel.ColumnClause[T], from, to T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		switch vi := column.(type) {
		case sequel.SQLColumnClause[T]:
			stmt.WriteString(vi.ColumnName() + " BETWEEN " + stmt.Var(vi.Convert(from)) + " AND " + stmt.Var(vi.Convert(to)))
		case sequel.ColumnConvertClause[T]:
			stmt.WriteString(vi.ColumnName() + " BETWEEN " + stmt.Var(vi.Convert(from)) + " AND " + stmt.Var(vi.Convert(to)))
		default:
			stmt.WriteString(column.ColumnName() + " BETWEEN " + stmt.Var(from) + " AND " + stmt.Var(to))
		}
	}
}

func NotBetween[T comparable](column sequel.ColumnClause[T], from, to T) sequel.WhereClause {
	return func(stmt sequel.StmtBuilder) {
		switch vi := column.(type) {
		case sequel.SQLColumnClause[T]:
			stmt.WriteString(vi.ColumnName() + " NOT BETWEEN " + stmt.Var(vi.Convert(from)) + " AND " + stmt.Var(vi.Convert(to)))
		case sequel.ColumnConvertClause[T]:
			stmt.WriteString(vi.ColumnName() + " NOT BETWEEN " + stmt.Var(vi.Convert(from)) + " AND " + stmt.Var(vi.Convert(to)))
		default:
			stmt.WriteString(column.ColumnName() + " NOT BETWEEN " + stmt.Var(from) + " AND " + stmt.Var(to))
		}
	}
}

func Set[T any](column sequel.ColumnClause[T], value ...T) sequel.SetClause {
	return func(stmt sequel.StmtBuilder) {
		if len(value) > 0 {
			switch vi := column.(type) {
			case sequel.SQLColumnClause[T]:
				stmt.WriteString(column.ColumnName() + " = " + stmt.Var(vi.Convert(value[0])))
			case sequel.ColumnConvertClause[T]:
				stmt.WriteString(column.ColumnName() + " = " + stmt.Var(vi.Convert(value[0])))
			default:
				stmt.WriteString(column.ColumnName() + " = " + stmt.Var(value[0]))
			}
		} else {
			stmt.WriteString(column.ColumnName() + " = " + stmt.Var(column.Value()))
		}
	}
}

func Asc[T any](column sequel.ColumnClause[T]) sequel.OrderByClause {
	return sequel.OrderByColumn(column.ColumnName(), true)
}

func Desc[T any](column sequel.ColumnClause[T]) sequel.OrderByClause {
	return sequel.OrderByColumn(column.ColumnName(), false)
}
