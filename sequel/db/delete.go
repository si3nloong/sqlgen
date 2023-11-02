package db

import (
	"context"
	"database/sql"

	"github.com/si3nloong/sqlgen/sequel"
)

// DeleteOne is to update single record using primary key.
func DeleteOne[T sequel.KeyValuer[T]](ctx context.Context, db sequel.DB, v T) (sql.Result, error) {
	var (
		pkName, _, pk = v.PK()
		dialect       = sequel.DefaultDialect()
		stmt          = acquireString()
	)
	defer releaseString(stmt)
	stmt.WriteString("DELETE FROM " + dialect.Wrap(v.TableName()) + " WHERE " + dialect.Wrap(pkName) + " = " + dialect.Var(1) + ";")

	return db.ExecContext(ctx, stmt.String(), pk)
}
