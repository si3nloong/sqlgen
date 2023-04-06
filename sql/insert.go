package sql

import (
	"context"
	"database/sql"
)

type InsertOnly[T any] interface {
	Valuer[T]
	InsertQuery() string
}

func InsertOne[T interface {
	InsertOnly[T]
}, Ptr Scanner[T]](ctx context.Context, db DB, v Ptr) error {
	err := db.QueryRowContext(ctx, (T)(*v).InsertQuery(), (T)(*v).Values()...).Scan(v.Addrs()...)
	return err
}

// InsertInto is a helper function to insert your records.
func InsertInto[T Valuer[T]](ctx context.Context, db DB, values []T) (sql.Result, error) {
	n := len(values)
	if n == 0 {
		return new(emptyResult), nil
	}

	stmt := acquireString()
	defer releaseString(stmt)
	columns := values[0].Columns()
	noOfCols := len(columns)
	args := make([]any, 0, n*noOfCols)

	stmt.WriteString("INSERT INTO " + dialect.Wrap(values[0].Table()) + " (")
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
			if columns[j] == "id" {
				stmt.WriteString("DEFAULT")
				continue
			}
			stmt.WriteString(dialect.Var(i + 1))
			i++
		}
		args = append(args, v.Values()[1:]...)
		stmt.WriteByte(')')
	}
	stmt.WriteByte(';')

	return db.ExecContext(ctx, stmt.String(), args...)
}
