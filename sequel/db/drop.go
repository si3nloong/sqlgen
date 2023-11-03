package db

import (
	"context"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/strpool"
)

func DropTable[T sequel.Tabler](ctx context.Context, db sequel.DB) error {
	var (
		v       T
		dialect = sequel.DefaultDialect()
		stmt    = strpool.AcquireString()
	)
	defer strpool.ReleaseString(stmt)
	stmt.WriteString("DROP TABLE " + dialect.Wrap(v.TableName()) + ";")
	if _, err := db.ExecContext(ctx, stmt.String()); err != nil {
		return err
	}
	return nil
}
