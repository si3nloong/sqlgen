package sqlite

import (
	"context"

	"github.com/si3nloong/sqlgen/codegen/dialect"
)

func (s *sqliteDriver) Migrate(ctx context.Context, dsn string, w dialect.Writer, schema dialect.Schema) error {
	return nil
}
