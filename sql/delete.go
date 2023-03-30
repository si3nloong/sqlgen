package sql

import (
	"context"
	"database/sql"
)

// DeleteOne is to update single record using primary key.
func DeleteOne[T KeyValuer[T]](ctx context.Context, db DB, v T) (sql.Result, error) {
	pk, err := v.PK()
	if err != nil {
		return nil, err
	}

	stmt := AcquireStmt()
	defer ReleaseStmt(stmt)
	stmt.WriteQuery("DELETE FROM "+dialect.Wrap(v.Table())+" WHERE "+dialect.Wrap(v.PKName())+" = "+dialect.Var(1)+";", pk)

	return db.ExecContext(ctx, stmt.Query(), stmt.Args()...)
}
