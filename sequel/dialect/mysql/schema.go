package mysql

import (
	"github.com/samber/lo"
	"github.com/si3nloong/sqlgen/sequel"
)

func (s *mysqlDriver) TableSchemas(table sequel.GoTableSchema) sequel.TableSchema {
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
	return schema
}
