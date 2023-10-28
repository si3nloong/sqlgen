package db

import (
	"context"
	"database/sql"

	"github.com/si3nloong/sqlgen/sequel"
)

func Migrate[T sequel.Migrator](ctx context.Context, db sequel.DB) (sql.Result, error) {
	var v T
	stmt := acquireString()
	defer releaseString(stmt)
	// TODO: support alter table as well
	// SELECT *
	// FROM information_schema.tables
	// WHERE table_schema = 'yourdb'
	// 	AND table_name = 'testtable'
	// LIMIT 1;
	// log.Println(stmt.String())
	result, err := db.ExecContext(ctx, v.CreateTableStmt())
	if err != nil {
		return nil, err
	}
	if affected, _ := result.RowsAffected(); affected > 0 {
		return result, nil
	}
	return db.ExecContext(ctx, v.AlterTableStmt())
}
