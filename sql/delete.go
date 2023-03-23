package sql

import (
	"context"
	"database/sql"
)

// DeleteOne is to update single record using primary key.
func DeleteOne[T KeyValuer[T]](ctx context.Context, db DB, v T) (sql.Result, error) {
	stmt := AcquireStmt()
	defer ReleaseStmt(stmt)

	pkName, pk := v.PK()
	stmt.WriteQuery("DELETE FROM "+Wrap(v.Table())+" WHERE "+Wrap(pkName)+" = "+Var(1)+";", pk)

	return db.ExecContext(ctx, stmt.Query(), stmt.Args()...)
}
