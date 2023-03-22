package sqlgen

import (
	"context"
	"database/sql"
	"log"
	"strings"
)

func Var(n int) string {
	return "?"
}

func Wrap(v string) string {
	return "`" + v + "`"
}

func InsertInto[T Valuer[T]](ctx context.Context, db DB, table string, values []T) (sql.Result, error) {
	if len(values) == 0 {
		return &emptyResult{}, nil
	}

	stmt := AcquireStmt()
	defer ReleaseStmt(stmt)
	columns := values[0].Columns()
	stmt.WriteQuery("INSERT INTO " + Wrap(table) + " (" + Wrap(strings.Join(columns, Wrap(","))) + ") VALUES ")
	noOfCols := len(columns)
	valueStr := "(" + strings.Repeat(",?", noOfCols)[1:] + ")"
	for n := len(values); n > 0; {
		value := values[0]
		stmt.WriteQuery(valueStr, value.Values()...)
		if n > 1 {
			stmt.WriteQuery(",")
		}
		values = values[1:]
		n = len(values)
	}
	stmt.WriteQuery(";")
	log.Println(stmt.Query())

	return db.ExecContext(ctx, stmt.Query(), stmt.Args()...)
}
