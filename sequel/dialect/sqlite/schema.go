package sqlite

import (
	"github.com/samber/lo"
	"github.com/si3nloong/sqlgen/sequel"
)

func (s *sqliteDriver) TableSchemas(table sequel.GoTableSchema) sequel.TableSchema {
	schema := new(tableDefinition)
	autoIncrKey, hasAutoIncr := table.AutoIncrKey()
	schema.keys = append(schema.keys, table.Keys()...)
	for i := range table.Columns() {
		col := table.Column(i)
		colDef := dataType(col)
		// key shouldn't have default value
		if lo.IndexOf(schema.keys, col.ColumnName()) > -1 {
			colDef.defaultValue = nil
		}
		if hasAutoIncr && autoIncrKey == col {
			colDef.extra = "AUTO_INCREMENT"
		}
		schema.cols = append(schema.cols, colDef)
	}
	for i := range table.Indexes() {
		idx := table.Index(i)
		idxDef := indexDefinition{}
		idxDef.cols = idx.Columns()
		switch idx.Type() {
		case "unique":
			idxDef.indexType = unique
		}
		schema.idxs = append(schema.idxs, &idxDef)
	}
	return schema
}
