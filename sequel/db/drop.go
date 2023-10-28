package db

import (
	"context"

	"github.com/si3nloong/sqlgen/sequel"
)

func DropTable[T sequel.Tabler](ctx context.Context, db sequel.DB) error {
	var v T
	stmt := acquireString()
	defer releaseString(stmt)
	stmt.WriteString("DROP TABLE " + v.TableName() + ";")
	if _, err := db.ExecContext(ctx, stmt.String()); err != nil {
		return err
	}
	return nil
}
