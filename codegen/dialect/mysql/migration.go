package mysql

import (
	"context"

	"github.com/si3nloong/sqlgen/codegen/dialect"
)

func (s *mysqlDriver) Migrate(ctx context.Context, dsn string, w dialect.Writer, schema dialect.TableMigrator) error {
	return nil
}
