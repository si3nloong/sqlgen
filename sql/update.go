package sql

import (
	"context"
	"database/sql"
)

// UpdateOne is to update single record using primary key.
func UpdateOne[T KeyValuer[T]](ctx context.Context, db DB, v T) (sql.Result, error) {
	stmt := AcquireStmt()
	defer ReleaseStmt(stmt)

	columns := v.Columns()
	values := v.Values()
	stmt.WriteQuery("UPDATE " + Wrap(v.Table()) + " SET ")
	noOfCols := len(columns)
	for i := 0; i < noOfCols; i++ {
		if i > 0 {
			stmt.WriteQuery(",")
		}
		stmt.WriteQuery(Wrap(columns[i])+" = "+Var(i+1), values[i])
	}
	pkName, pk := v.PK()
	stmt.WriteQuery(" WHERE "+Wrap(pkName)+" = "+Var(noOfCols+2)+";", pk)

	return db.ExecContext(ctx, stmt.Query(), stmt.Args()...)
}
