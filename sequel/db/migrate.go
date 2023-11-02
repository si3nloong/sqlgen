package db

import (
	"context"

	"github.com/si3nloong/sqlgen/sequel"
)

func Migrate[T sequel.Migrator](ctx context.Context, db sequel.DB) error {
	var (
		dialect     = sequel.DefaultDialect()
		v           T
		table       string
		tableExists bool
		stmt        = acquireString()
	)
	defer releaseString(stmt)
	if err := db.QueryRowContext(ctx, "SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_NAME = "+dialect.Var(1)+" LIMIT 1;", v.TableName()).Scan(&table); err != nil {
		tableExists = false
	} else {
		tableExists = true
	}
	if tableExists {
		if _, err := db.ExecContext(ctx, v.AlterTableStmt()); err != nil {
			return err
		}
		return nil
	}
	if _, err := db.ExecContext(ctx, v.CreateTableStmt()); err != nil {
		return err
	}
	return nil
}
