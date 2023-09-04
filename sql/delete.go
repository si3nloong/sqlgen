package sql

import (
	"context"
	"database/sql"
)

// DeleteOne is to update single record using primary key.
func DeleteOne[T KeyValuer[T]](ctx context.Context, db DB, v T) (sql.Result, error) {
	pkName, _, pk := v.PK()

	stmt := acquireString()
	defer releaseString(stmt)
	stmt.WriteString("DELETE FROM " + dialect.Wrap(v.Table()) + " WHERE " + dialect.Wrap(pkName) + " = " + dialect.Var(1) + ";")

	return db.ExecContext(ctx, stmt.String(), pk)
}
