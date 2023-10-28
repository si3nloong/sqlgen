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
	stmt.WriteString("DELETE FROM `" + v.TableName() + "` WHERE " + pkName + " = " + ";")

	return db.ExecContext(ctx, stmt.String(), pk)
}
