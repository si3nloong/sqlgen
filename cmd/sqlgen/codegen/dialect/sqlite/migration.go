package sqlite

import (
	"context"

	"github.com/si3nloong/sqlgen/cmd/sqlgen/codegen/dialect"
)

func (s *sqliteDriver) Migrate(ctx context.Context, dsn string, w dialect.Writer, schema dialect.TableMigrator) error {
	return nil
}
