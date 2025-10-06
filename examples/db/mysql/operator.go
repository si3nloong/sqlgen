package mysqldb

import (
	"github.com/si3nloong/sqlgen/sequel"
)

func And(clauses ...sequel.WhereClause) sequel.WhereClause {
	return func(w sequel.StmtWriter) {
		if n := len(clauses); n > 0 {
			w.WriteByte('(')
			clauses[0](w)
			for i := 1; i < n; i++ {
				w.WriteString(" AND ")
				clauses[i](w)
			}
			w.WriteByte(')')
		}
	}
}

func Or(clauses ...sequel.WhereClause) sequel.WhereClause {
	return func(w sequel.StmtWriter) {
		if n := len(clauses); n > 0 {
			w.WriteByte('(')
			clauses[0](w)
			for i := 1; i < n; i++ {
				w.WriteString(" OR ")
				clauses[i](w)
			}
			w.WriteByte(')')
		}
	}
}

func Equal[T comparable](column sequel.ColumnClause[T], value T) sequel.WhereClause {
	return func(w sequel.StmtWriter) {
		switch vi := column.(type) {
		case sequel.SQLColumnClause[T]:
			w.WriteString(vi.ColumnName() + " = " + w.Var(vi.Convert(value)))
		case sequel.ColumnConvertClause[T]:
			w.WriteString(vi.ColumnName() + " = " + w.Var(vi.Convert(value)))
		default:
			w.WriteString(column.ColumnName() + " = " + w.Var(value))
		}
	}
}

func NotEqual[T comparable](column sequel.ColumnClause[T], value T) sequel.WhereClause {
	return func(w sequel.StmtWriter) {
		switch vi := column.(type) {
		case sequel.SQLColumnClause[T]:
			w.WriteString(vi.ColumnName() + " <> " + w.Var(vi.Convert(value)))
		case sequel.ColumnConvertClause[T]:
			w.WriteString(vi.ColumnName() + " <> " + w.Var(vi.Convert(value)))
		default:
			w.WriteString(column.ColumnName() + " <> " + w.Var(value))
		}
	}
}

func In[T any](column sequel.ColumnClause[T], values []T) sequel.WhereClause {
	return func(w sequel.StmtWriter) {
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
		w.WriteString(column.ColumnName() + " IN " + w.Vars(args))
	}
}

func NotIn[T any](column sequel.ColumnClause[T], values []T) sequel.WhereClause {
	return func(w sequel.StmtWriter) {
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
		w.WriteString(column.ColumnName() + " NOT IN " + w.Vars(args))
	}
}

func GreaterThan[T comparable](column sequel.ColumnClause[T], value T) sequel.WhereClause {
	return func(w sequel.StmtWriter) {
		switch vi := column.(type) {
		case sequel.SQLColumnClause[T]:
			w.WriteString(vi.ColumnName() + " > " + w.Var(vi.Convert(value)))
		case sequel.ColumnConvertClause[T]:
			w.WriteString(vi.ColumnName() + " > " + w.Var(vi.Convert(value)))
		default:
			w.WriteString(column.ColumnName() + " > " + w.Var(value))
		}
	}
}

func GreaterThanOrEqual[T comparable](column sequel.ColumnClause[T], value T) sequel.WhereClause {
	return func(w sequel.StmtWriter) {
		switch vi := column.(type) {
		case sequel.SQLColumnClause[T]:
			w.WriteString(vi.ColumnName() + " >= " + w.Var(vi.Convert(value)))
		case sequel.ColumnConvertClause[T]:
			w.WriteString(vi.ColumnName() + " >= " + w.Var(vi.Convert(value)))
		default:
			w.WriteString(column.ColumnName() + " >= " + w.Var(value))
		}
	}
}

func LessThan[T comparable](column sequel.ColumnClause[T], value T) sequel.WhereClause {
	return func(w sequel.StmtWriter) {
		switch vi := column.(type) {
		case sequel.SQLColumnClause[T]:
			w.WriteString(vi.ColumnName() + " < " + w.Var(vi.Convert(value)))
		case sequel.ColumnConvertClause[T]:
			w.WriteString(vi.ColumnName() + " < " + w.Var(vi.Convert(value)))
		default:
			w.WriteString(column.ColumnName() + " < " + w.Var(value))
		}
	}
}

func LessThanOrEqual[T comparable](column sequel.ColumnClause[T], value T) sequel.WhereClause {
	return func(w sequel.StmtWriter) {
		switch vi := column.(type) {
		case sequel.SQLColumnClause[T]:
			w.WriteString(vi.ColumnName() + " <= " + w.Var(vi.Convert(value)))
		case sequel.ColumnConvertClause[T]:
			w.WriteString(vi.ColumnName() + " <= " + w.Var(vi.Convert(value)))
		default:
			w.WriteString(column.ColumnName() + " <= " + w.Var(value))
		}
	}
}

func Like[T comparable](column sequel.ColumnClause[T], value T) sequel.WhereClause {
	return func(w sequel.StmtWriter) {
		switch vi := column.(type) {
		case sequel.SQLColumnClause[T]:
			w.WriteString(vi.ColumnName() + " LIKE " + w.Var(vi.Convert(value)))
		case sequel.ColumnConvertClause[T]:
			w.WriteString(vi.ColumnName() + " LIKE " + w.Var(vi.Convert(value)))
		default:
			w.WriteString(column.ColumnName() + " LIKE " + w.Var(value))
		}
	}
}

func NotLike[T comparable](column sequel.ColumnClause[T], value T) sequel.WhereClause {
	return func(w sequel.StmtWriter) {
		switch vi := column.(type) {
		case sequel.SQLColumnClause[T]:
			w.WriteString(vi.ColumnName() + " NOT LIKE " + w.Var(vi.Convert(value)))
		case sequel.ColumnConvertClause[T]:
			w.WriteString(vi.ColumnName() + " NOT LIKE " + w.Var(vi.Convert(value)))
		default:
			w.WriteString(column.ColumnName() + " NOT LIKE " + w.Var(value))
		}
	}
}

func IsNull[T any](column sequel.ColumnClause[T]) sequel.WhereClause {
	return func(w sequel.StmtWriter) {
		w.WriteString(column.ColumnName() + " IS NULL")
	}
}

func IsNotNull[T any](column sequel.ColumnClause[T]) sequel.WhereClause {
	return func(w sequel.StmtWriter) {
		w.WriteString(column.ColumnName() + " IS NOT NULL")
	}
}

func Between[T comparable](column sequel.ColumnClause[T], from, to T) sequel.WhereClause {
	return func(w sequel.StmtWriter) {
		switch vi := column.(type) {
		case sequel.SQLColumnClause[T]:
			w.WriteString(vi.ColumnName() + " BETWEEN " + w.Var(vi.Convert(from)) + " AND " + w.Var(vi.Convert(to)))
		case sequel.ColumnConvertClause[T]:
			w.WriteString(vi.ColumnName() + " BETWEEN " + w.Var(vi.Convert(from)) + " AND " + w.Var(vi.Convert(to)))
		default:
			w.WriteString(column.ColumnName() + " BETWEEN " + w.Var(from) + " AND " + w.Var(to))
		}
	}
}

func NotBetween[T comparable](column sequel.ColumnClause[T], from, to T) sequel.WhereClause {
	return func(w sequel.StmtWriter) {
		switch vi := column.(type) {
		case sequel.SQLColumnClause[T]:
			w.WriteString(vi.ColumnName() + " NOT BETWEEN " + w.Var(vi.Convert(from)) + " AND " + w.Var(vi.Convert(to)))
		case sequel.ColumnConvertClause[T]:
			w.WriteString(vi.ColumnName() + " NOT BETWEEN " + w.Var(vi.Convert(from)) + " AND " + w.Var(vi.Convert(to)))
		default:
			w.WriteString(column.ColumnName() + " NOT BETWEEN " + w.Var(from) + " AND " + w.Var(to))
		}
	}
}

func Set[T any](column sequel.ColumnClause[T], value T) sequel.SetClause {
	return func(w sequel.StmtWriter) {
		switch vi := column.(type) {
		case sequel.SQLColumnClause[T]:
			w.WriteString(column.ColumnName() + " = " + w.Var(vi.Convert(value)))
		case sequel.ColumnConvertClause[T]:
			w.WriteString(column.ColumnName() + " = " + w.Var(vi.Convert(value)))
		default:
			w.WriteString(column.ColumnName() + " = " + w.Var(value))
		}
	}
}

func Asc[T any](column sequel.ColumnClause[T]) sequel.OrderByClause {
	return sequel.OrderByColumn(column.ColumnName(), true)
}

func Desc[T any](column sequel.ColumnClause[T]) sequel.OrderByClause {
	return sequel.OrderByColumn(column.ColumnName(), false)
}
