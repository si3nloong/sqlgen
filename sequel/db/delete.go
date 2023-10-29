package db

import (
	"context"
	"database/sql"

	"github.com/si3nloong/sqlgen/sequel"
)

// DeleteOne is to update single record using primary key.
func DeleteOne[T sequel.KeyValuer[T]](ctx context.Context, db sequel.DB, v T) (sql.Result, error) {
	pkName, _, pk := v.PK()

	stmt := acquireString()
	defer releaseString(stmt)
	stmt.WriteString("DELETE FROM `" + v.TableName() + "` WHERE " + pkName + " = ?;")

	return db.ExecContext(ctx, stmt.String(), pk)
}
