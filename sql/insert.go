package sql

import (
	"context"
	"database/sql"
)

type InsertOnly[T any] interface {
	Valuer[T]
	InsertQuery() string
}

func InsertOne[T Valuer[T], Ptr interface {
	Valuer[T]
	Scanner[T]
}](ctx context.Context, db DB, v Ptr) (sql.Result, error) {
	columns, args := v.Columns(), v.Values()
	switch vi := any(v).(type) {
	case Keyer:
		if vi.IsAutoIncr() {
			idx, _ := vi.PK()
			columns = append(columns[:idx], columns[idx+1:]...)
			args = append(args[:idx], args[idx+1:]...)
		}
	}

	noOfCols := len(columns)
	stmt := acquireString()
	defer releaseString(stmt)
	stmt.WriteString("INSERT INTO " + dialect.Wrap(v.Table()) + " (")
	for i := 0; i < noOfCols; i++ {
		if i > 0 {
			stmt.WriteString("," + dialect.Wrap(columns[i]))
		} else {
			stmt.WriteString(dialect.Wrap(columns[i]))
		}
	}
	stmt.WriteString(") VALUES (")
	for i := range args {
		if i > 0 {
			stmt.WriteByte(',')
		}
		stmt.WriteString(dialect.Var(i + 1))
	}
	stmt.WriteByte(')')
	// stmt.WriteString(`) RETURNING `)
	// for i, col := range v.Columns() {
	// 	if i > 0 {
	// 		stmt.WriteString("," + dialect.Wrap(col))
	// 	} else {
	// 		stmt.WriteString(dialect.Wrap(col))
	// 	}
	// }
	// return db.QueryRowContext(ctx, stmt.String(), args...).Scan(v.Addrs()...)
	stmt.WriteByte(';')
	return db.ExecContext(ctx, stmt.String(), args...)
}

// InsertInto is a helper function to insert your records.
func InsertInto[T Valuer[T]](ctx context.Context, db DB, values []T) (sql.Result, error) {
	n := len(values)
	if n == 0 {
		return new(emptyResult), nil
	}

	model := values[0]
	columns := model.Columns()
	idx := -1
	switch vi := any(model).(type) {
	case Keyer:
		if vi.IsAutoIncr() {
			idx, _ = vi.PK()
			columns = append(columns[:idx], columns[idx+1:]...)
		}
	}
	noOfCols := len(columns)
	args := make([]any, 0)

	stmt := acquireString()
	defer releaseString(stmt)
	stmt.WriteString("INSERT INTO " + dialect.Wrap(model.Table()) + " (")
	for i := 0; i < noOfCols; i++ {
		if i > 0 {
			stmt.WriteString("," + dialect.Wrap(columns[i]))
		} else {
			stmt.WriteString(dialect.Wrap(columns[i]))
		}
	}
	stmt.WriteString(") VALUES ")
	var (
		i int
	)
	for k, v := range values {
		if k > 0 {
			stmt.WriteString(",(")
		} else {
			stmt.WriteByte('(')
		}
		for j := 0; j < noOfCols; j++ {
			if j > 0 {
				stmt.WriteByte(',')
			}
			stmt.WriteString(dialect.Var(i + 1))
			i++
		}
		if idx > -1 {
			values := v.Values()
			values = append(values[:idx], values[idx+1:]...)
			args = append(args, values...)
		} else {
			args = append(args, v.Values()...)
		}
		stmt.WriteByte(')')
	}
	stmt.WriteByte(';')

	return db.ExecContext(ctx, stmt.String(), args...)
}
