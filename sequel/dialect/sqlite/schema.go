package sqlite

import (
	"github.com/si3nloong/sqlgen/sequel"
)

func (s *sqliteDriver) TableSchemas(table sequel.GoTableSchema) sequel.TableSchema {
	return nil
}
