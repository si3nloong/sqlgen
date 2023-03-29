package sql

import (
	"context"
	"database/sql"
	"strings"
)

// InsertInto is a helper function to insert your records.
func InsertInto[T Valuer[T]](ctx context.Context, db DB, values []T) (sql.Result, error) {
	if len(values) == 0 {
		return new(emptyResult), nil
	}

	stmt := AcquireStmt()
	defer ReleaseStmt(stmt)
	var ent T
	columns := ent.Columns()
	stmt.WriteQuery("INSERT INTO " + dialect.Wrap(ent.Table()) + " (" + dialect.Wrap(strings.Join(columns, dialect.Wrap(","))) + ") VALUES ")
	noOfCols := len(columns)
	valueStr := "(" + strings.Repeat(",?", noOfCols)[1:] + ")"
	for n := len(values); n > 0; {
		ent = values[0]
		stmt.WriteQuery(valueStr, ent.Values()...)
		if n > 1 {
			stmt.WriteQuery(",")
		}
		values = values[1:]
		n = len(values)
	}
	stmt.WriteQuery(";")

	return db.ExecContext(ctx, stmt.Query(), stmt.Args()...)
}
